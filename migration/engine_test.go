package migration

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestMigration is a simple test migration
type TestMigration struct {
	version      string
	description  string
	upExecuted   bool
	downExecuted bool
}

func (m *TestMigration) Version() string {
	return m.version
}

func (m *TestMigration) Description() string {
	return m.description
}

func (m *TestMigration) Up(ctx context.Context, db *mongo.Database) error {
	m.upExecuted = true
	return nil
}

func (m *TestMigration) Down(ctx context.Context, db *mongo.Database) error {
	m.downExecuted = true
	return nil
}

func TestNewEngine(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("creates engine with correct parameters", func(mt *mtest.T) {
		engine := NewEngine(mt.DB, "test_migrations")

		if engine.db != mt.DB {
			t.Error("Engine database not set correctly")
		}

		if engine.migrationsCollection != "test_migrations" {
			t.Error("Engine migrations collection not set correctly")
		}

		if engine.migrations == nil {
			t.Error("Engine migrations map not initialized")
		}
	})
}

func TestRegister(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("registers migration correctly", func(mt *mtest.T) {
		engine := NewEngine(mt.DB, "test_migrations")
		migration := &TestMigration{
			version:     "20240101_001",
			description: "Test migration",
		}

		engine.Register(migration)

		if len(engine.migrations) != 1 {
			t.Errorf("Expected 1 migration, got %d", len(engine.migrations))
		}

		registered, exists := engine.migrations["20240101_001"]
		if !exists {
			t.Error("Migration not registered")
		}

		if registered != migration {
			t.Error("Registered migration is not the same instance")
		}
	})
}

func TestRegisterMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("registers multiple migrations correctly", func(mt *mtest.T) {
		engine := NewEngine(mt.DB, "test_migrations")

		migration1 := &TestMigration{version: "20240101_001", description: "Test migration 1"}
		migration2 := &TestMigration{version: "20240101_002", description: "Test migration 2"}
		migration3 := &TestMigration{version: "20240101_003", description: "Test migration 3"}

		engine.RegisterMany(migration1, migration2, migration3)

		if len(engine.migrations) != 3 {
			t.Errorf("Expected 3 migrations, got %d", len(engine.migrations))
		}

		for _, migration := range []*TestMigration{migration1, migration2, migration3} {
			registered, exists := engine.migrations[migration.Version()]
			if !exists {
				t.Errorf("Migration %s not registered", migration.Version())
			}
			if registered != migration {
				t.Errorf("Registered migration %s is not the same instance", migration.Version())
			}
		}
	})
}

func TestDirection(t *testing.T) {
	tests := []struct {
		direction Direction
		expected  string
	}{
		{DirectionUp, "up"},
		{DirectionDown, "down"},
		{Direction(999), "unknown"},
	}

	for _, test := range tests {
		if test.direction.String() != test.expected {
			t.Errorf("Direction %d should return %s, got %s",
				test.direction, test.expected, test.direction.String())
		}
	}
}

func TestMigrationStatus(t *testing.T) {
	status := MigrationStatus{
		Version:     "20240101_001",
		Description: "Test migration",
		Applied:     true,
		AppliedAt:   &time.Time{},
	}

	if status.Version != "20240101_001" {
		t.Error("Version not set correctly")
	}

	if status.Description != "Test migration" {
		t.Error("Description not set correctly")
	}

	if !status.Applied {
		t.Error("Applied status not set correctly")
	}

	if status.AppliedAt == nil {
		t.Error("AppliedAt should not be nil")
	}
}

func TestErrNotSupported(t *testing.T) {
	err := ErrNotSupported{Operation: "test operation"}
	expected := "operation not supported: test operation"

	if err.Error() != expected {
		t.Errorf("Expected error message %s, got %s", expected, err.Error())
	}
}
