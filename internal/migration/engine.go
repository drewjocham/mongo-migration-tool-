package migration

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Engine is the main migration engine
type Engine struct {
	db                   *mongo.Database
	migrationsCollection string
	migrations           map[string]Migration
}

// NewEngine creates a new migration engine
func NewEngine(db *mongo.Database, migrationsCollection string) *Engine {
	return &Engine{
		db:                   db,
		migrationsCollection: migrationsCollection,
		migrations:           make(map[string]Migration),
	}
}

// Register registers a migration with the engine
func (e *Engine) Register(migration Migration) {
	e.migrations[migration.Version()] = migration
}

// GetStatus returns the status of all migrations
func (e *Engine) GetStatus(ctx context.Context) ([]MigrationStatus, error) {
	// Get applied migrations from database
	appliedMigrations, err := e.getAppliedMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]MigrationRecord)
	for _, record := range appliedMigrations {
		appliedMap[record.Version] = record
	}

	// Get all available versions
	var allVersions []string
	for version := range e.migrations {
		allVersions = append(allVersions, version)
	}

	// Add applied migrations that might not have corresponding migration files
	for version := range appliedMap {
		if _, exists := e.migrations[version]; !exists {
			allVersions = append(allVersions, version)
		}
	}

	sort.Strings(allVersions)

	var status []MigrationStatus
	for _, version := range allVersions {
		migration := e.migrations[version]
		applied, exists := appliedMap[version]

		description := ""
		if migration != nil {
			description = migration.Description()
		} else if exists {
			description = applied.Description
		}

		migrationStatus := MigrationStatus{
			Version:     version,
			Description: description,
			Applied:     exists,
		}

		if exists {
			migrationStatus.AppliedAt = &applied.AppliedAt
		}

		status = append(status, migrationStatus)
	}

	return status, nil
}

// Up runs all pending migrations
func (e *Engine) Up(ctx context.Context, target string) error {
	return e.migrate(ctx, DirectionUp, target)
}

// Down rolls back migrations to the target version
func (e *Engine) Down(ctx context.Context, target string) error {
	return e.migrate(ctx, DirectionDown, target)
}

// Force marks a migration as applied without actually running it
func (e *Engine) Force(ctx context.Context, version string) error {
	migration, exists := e.migrations[version]
	if !exists {
		return fmt.Errorf("migration %s not found", version)
	}

	record := MigrationRecord{
		Version:     version,
		Description: migration.Description(),
		AppliedAt:   time.Now().UTC(),
		Checksum:    e.calculateChecksum(migration),
	}

	collection := e.db.Collection(e.migrationsCollection)
	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"version": version},
		record,
		options.Replace().SetUpsert(true),
	)

	if err != nil {
		return fmt.Errorf("failed to force migration %s: %w", version, err)
	}

	fmt.Printf("Forced migration %s\n", version)
	return nil
}

// migrate executes migrations in the specified direction
func (e *Engine) migrate(ctx context.Context, direction Direction, target string) error {
	appliedMigrations, err := e.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, record := range appliedMigrations {
		appliedMap[record.Version] = true
	}

	// Get migrations to execute
	migrationsToExecute, err := e.getMigrationsToExecute(direction, target, appliedMap)
	if err != nil {
		return fmt.Errorf("failed to determine migrations to execute: %w", err)
	}

	if len(migrationsToExecute) == 0 {
		fmt.Printf("No migrations to %s\n", direction.String())
		return nil
	}

	// Execute migrations
	for _, version := range migrationsToExecute {
		migration := e.migrations[version]
		if migration == nil {
			return fmt.Errorf("migration %s not found", version)
		}

		fmt.Printf("Running migration %s %s: %s\n", version, direction.String(), migration.Description())

		err := e.executeMigration(ctx, migration, direction)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s %s: %w", version, direction.String(), err)
		}

		fmt.Printf("Successfully executed migration %s %s\n", version, direction.String())
	}

	return nil
}

// executeMigration executes a single migration
func (e *Engine) executeMigration(ctx context.Context, migration Migration, direction Direction) error {
	version := migration.Version()

	switch direction {
	case DirectionUp:
		// Execute the migration
		if err := migration.Up(ctx, e.db); err != nil {
			return err
		}

		// Record the migration as applied
		record := MigrationRecord{
			Version:     version,
			Description: migration.Description(),
			AppliedAt:   time.Now().UTC(),
			Checksum:    e.calculateChecksum(migration),
		}

		collection := e.db.Collection(e.migrationsCollection)
		_, err := collection.InsertOne(ctx, record)
		return err

	case DirectionDown:
		// Execute the rollback
		if err := migration.Down(ctx, e.db); err != nil {
			return err
		}

		// Remove the migration record
		collection := e.db.Collection(e.migrationsCollection)
		_, err := collection.DeleteOne(ctx, bson.M{"version": version})
		return err

	default:
		return fmt.Errorf("unknown direction: %v", direction)
	}
}

// getAppliedMigrations retrieves all applied migrations from the database
func (e *Engine) getAppliedMigrations(ctx context.Context) ([]MigrationRecord, error) {
	collection := e.db.Collection(e.migrationsCollection)

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"version": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []MigrationRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// getMigrationsToExecute determines which migrations need to be executed
func (e *Engine) getMigrationsToExecute(direction Direction, target string, appliedMap map[string]bool) ([]string, error) {
	var versions []string
	for version := range e.migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	var toExecute []string

	switch direction {
	case DirectionUp:
		for _, version := range versions {
			if !appliedMap[version] {
				toExecute = append(toExecute, version)
				if target != "" && version == target {
					break
				}
			}
		}

	case DirectionDown:
		// Reverse order for rollbacks
		for i := len(versions) - 1; i >= 0; i-- {
			version := versions[i]
			if appliedMap[version] {
				toExecute = append(toExecute, version)
				if target != "" && version == target {
					break
				}
			}
		}
	}

	return toExecute, nil
}

// calculateChecksum calculates a checksum for the migration
func (e *Engine) calculateChecksum(migration Migration) string {
	data := fmt.Sprintf("%s:%s", migration.Version(), migration.Description())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
