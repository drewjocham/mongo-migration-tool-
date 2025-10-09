// Package migration provides MongoDB migration functionality for database schema versioning.
//
// This package allows you to version control your MongoDB schema changes using up/down
// migrations, similar to Liquibase or Flyway for relational databases.
//
// Basic Usage:
//
//	// Create a new migration engine
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	db := client.Database("myapp")
//	engine := migration.NewEngine(db, "schema_migrations")
//
//	// Register your migrations
//	engine.Register(&MyMigration{})
//
//	// Run pending migrations
//	if err := engine.Up(ctx, ""); err != nil {
//		log.Fatal(err)
//	}
package migration

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Migration represents a single database migration.
//
// Each migration must have a unique version identifier and provide both
// up and down implementations for bidirectional migrations.
type Migration interface {
	// Version returns the unique version identifier for this migration.
	// Recommended format: YYYYMMDD_NNN (e.g., "20240101_001")
	Version() string

	// Description returns a human-readable description of what this migration does.
	Description() string

	// Up executes the migration, applying changes to the database.
	Up(ctx context.Context, db *mongo.Database) error

	// Down rolls back the migration, undoing changes made by Up.
	Down(ctx context.Context, db *mongo.Database) error
}

// MigrationRecord represents a migration record stored in the database.
//
// This is used internally to track which migrations have been applied
// and when they were executed.
type MigrationRecord struct {
	Version     string    `bson:"version"`
	Description string    `bson:"description"`
	AppliedAt   time.Time `bson:"applied_at"`
	Checksum    string    `bson:"checksum,omitempty"`
}

// Direction represents the migration direction (up or down).
type Direction int

const (
	// DirectionUp indicates applying migrations forward
	DirectionUp Direction = iota
	// DirectionDown indicates rolling back migrations
	DirectionDown
)

// String returns a string representation of the direction.
func (d Direction) String() string {
	switch d {
	case DirectionUp:
		return "up"
	case DirectionDown:
		return "down"
	default:
		return "unknown"
	}
}

// MigrationStatus represents the status of a migration.
//
// This shows whether a migration has been applied and when.
type MigrationStatus struct {
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Applied     bool       `json:"applied"`
	AppliedAt   *time.Time `json:"applied_at,omitempty"`
}

// ErrNotSupported is returned when a migration doesn't support an operation.
type ErrNotSupported struct {
	Operation string
}

// Error returns the error message.
func (e ErrNotSupported) Error() string {
	return "operation not supported: " + e.Operation
}
