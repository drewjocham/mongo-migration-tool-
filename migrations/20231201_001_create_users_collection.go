package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration_20231201_001 creates the users collection with indexes
type Migration_20231201_001 struct{}

// Version returns the unique version identifier for this migration
func (m *Migration_20231201_001) Version() string {
	return "20231201_001"
}

// Description returns a human-readable description of what this migration does
func (m *Migration_20231201_001) Description() string {
	return "Create users collection with email and username indexes"
}

// Up executes the migration
func (m *Migration_20231201_001) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// Create unique index on email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}

	// Create unique index on username
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index().SetUnique(true).SetName("username_unique"),
	}

	// Create compound index on created_at and status
	compoundIndex := mongo.IndexModel{
		Keys:    bson.D{{"created_at", -1}, {"status", 1}},
		Options: options.Index().SetName("created_at_status"),
	}

	// Create all indexes
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		emailIndex,
		usernameIndex,
		compoundIndex,
	})
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	fmt.Println("Created users collection with indexes: email_unique, username_unique, created_at_status")
	return nil
}

// Down rolls back the migration
func (m *Migration_20231201_001) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// Drop specific indexes
	indexes := []string{"email_unique", "username_unique", "created_at_status"}
	for _, indexName := range indexes {
		_, err := collection.Indexes().DropOne(ctx, indexName)
		if err != nil {
			// Continue dropping other indexes even if one fails
			fmt.Printf("Warning: failed to drop index %s: %v\n", indexName, err)
		}
	}

	// Optionally drop the entire collection
	// return collection.Drop(ctx)

	fmt.Println("Dropped indexes from users collection")
	return nil
}
