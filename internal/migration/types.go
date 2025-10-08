package migration

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Migration represents a single database migration
type Migration interface {
	Version() string
	Description() string
	Up(ctx context.Context, db *mongo.Database) error
	Down(ctx context.Context, db *mongo.Database) error
}

// MigrationRecord represents a migration record stored in the database
type MigrationRecord struct {
	Version     string    `bson:"version"`
	Description string    `bson:"description"`
	AppliedAt   time.Time `bson:"applied_at"`
	Checksum    string    `bson:"checksum,omitempty"`
}

// Direction represents the migration direction
type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
)

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

// MigrationStatus represents the status of migrations
type MigrationStatus struct {
	Version     string
	Description string
	Applied     bool
	AppliedAt   *time.Time
}

// ErrNotSupported is returned when a migration doesn't support rollback
type ErrNotSupported struct {
	Operation string
}

func (e ErrNotSupported) Error() string {
	return "operation not supported: " + e.Operation
}
