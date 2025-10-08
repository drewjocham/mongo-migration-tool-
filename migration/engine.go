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

// Engine is the main migration engine for executing and tracking MongoDB migrations.
//
// The engine manages migration state, executes migrations in order, and provides
// rollback capabilities. It stores migration records in a MongoDB collection to
// track which migrations have been applied.
//
// Example usage:
//
//	engine := migration.NewEngine(db, "schema_migrations")
//	engine.Register(&CreateUsersCollection{})
//	engine.Register(&AddUserIndexes{})
//
//	// Run all pending migrations
//	if err := engine.Up(ctx, ""); err != nil {
//		log.Fatal(err)
//	}
type Engine struct {
	db                   *mongo.Database
	migrationsCollection string
	migrations           map[string]Migration
}

// NewEngine creates a new migration engine.
//
// Parameters:
//   - db: MongoDB database instance
//   - migrationsCollection: Name of the collection to store migration records
//
// The migrations collection will be created automatically if it doesn't exist.
func NewEngine(db *mongo.Database, migrationsCollection string) *Engine {
	return &Engine{
		db:                   db,
		migrationsCollection: migrationsCollection,
		migrations:           make(map[string]Migration),
	}
}

// Register registers a migration with the engine.
//
// All migrations must be registered before calling Up, Down, or GetStatus.
// Migrations with duplicate versions will overwrite previously registered migrations.
//
// Example:
//
//	engine.Register(&CreateUsersCollection{})
//	engine.Register(&AddUserIndexes{})
func (e *Engine) Register(migration Migration) {
	e.migrations[migration.Version()] = migration
}

// RegisterMany registers multiple migrations at once.
//
// This is a convenience method for registering multiple migrations:
//
//	engine.RegisterMany(
//		&CreateUsersCollection{},
//		&AddUserIndexes{},
//		&CreateOrdersCollection{},
//	)
func (e *Engine) RegisterMany(migrations ...Migration) {
	for _, migration := range migrations {
		e.Register(migration)
	}
}

// GetStatus returns the status of all migrations (both registered and applied).
//
// This includes migrations that have been applied but no longer have corresponding
// registered Migration implementations, which can happen during development or
// when cleaning up old migrations.
//
// Returns a slice of MigrationStatus sorted by version.
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

// Up runs migrations forward to the specified target version.
//
// If target is empty, all pending migrations will be executed.
// If target is specified, migrations will be executed up to and including that version.
//
// Migrations are executed in version order (lexicographic sorting).
//
// Example:
//
//	// Run all pending migrations
//	err := engine.Up(ctx, "")
//
//	// Run migrations up to a specific version
//	err := engine.Up(ctx, "20240101_003")
func (e *Engine) Up(ctx context.Context, target string) error {
	return e.migrate(ctx, DirectionUp, target)
}

// Down rolls back migrations to the specified target version.
//
// Migrations are rolled back in reverse version order until the target is reached.
// The target migration itself will NOT be rolled back.
//
// Example:
//
//	// Rollback to version 20240101_001 (this version will remain applied)
//	err := engine.Down(ctx, "20240101_001")
func (e *Engine) Down(ctx context.Context, target string) error {
	return e.migrate(ctx, DirectionDown, target)
}

// Force marks a migration as applied without actually running it.
//
// This is useful for handling migrations that were applied manually or
// for fixing migration state inconsistencies. Use with caution.
//
// Example:
//
//	// Mark a migration as applied without running it
//	err := engine.Force(ctx, "20240101_001")
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
		return nil
	}

	// Execute migrations
	for _, version := range migrationsToExecute {
		migration := e.migrations[version]
		if migration == nil {
			return fmt.Errorf("migration %s not found", version)
		}

		err := e.executeMigration(ctx, migration, direction)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s %s: %w", version, direction.String(), err)
		}
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
				if target != "" && version == target {
					break
				}
				toExecute = append(toExecute, version)
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
