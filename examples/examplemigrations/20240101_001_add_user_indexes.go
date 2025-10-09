package examplemigrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddUserIndexesMigration adds indexes to the users collection
type AddUserIndexesMigration struct{}

func (m *AddUserIndexesMigration) Version() string {
	return "20240101_001"
}

func (m *AddUserIndexesMigration) Description() string {
	return "Add indexes to users collection for email and created_at fields"
}

func (m *AddUserIndexesMigration) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// Create index for email field (unique)
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().
			SetName("idx_users_email_unique").
			SetUnique(true).
			SetBackground(true),
	}

	// Create index for created_at field
	createdAtIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().
			SetName("idx_users_created_at").
			SetBackground(true),
	}

	// Create compound index for status and created_at
	compoundIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().
			SetName("idx_users_status_created_at").
			SetBackground(true),
	}

	// Create all indexes
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		emailIndexModel,
		createdAtIndexModel,
		compoundIndexModel,
	})

	return err
}

func (m *AddUserIndexesMigration) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// Drop the indexes we created
	indexNames := []string{
		"idx_users_email_unique",
		"idx_users_created_at",
		"idx_users_status_created_at",
	}

	for _, indexName := range indexNames {
		// Ignore errors when dropping indexes - they might not exist
		_, _ = collection.Indexes().DropOne(ctx, indexName)
	}

	return nil
}
