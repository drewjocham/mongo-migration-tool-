package examplemigrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAuditCollectionMigration creates a new audit collection with validation
type CreateAuditCollectionMigration struct{}

func (m *CreateAuditCollectionMigration) Version() string {
	return "20240101_003"
}

func (m *CreateAuditCollectionMigration) Description() string {
	return "Create audit collection with schema validation and indexes"
}

func (m *CreateAuditCollectionMigration) Up(ctx context.Context, db *mongo.Database) error {
	// Define JSON schema validation
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"user_id", "action", "timestamp", "resource_type", "resource_id"},
			"properties": bson.M{
				"user_id": bson.M{
					"bsonType":    "objectId",
					"description": "ID of the user who performed the action",
				},
				"action": bson.M{
					"bsonType":    "string",
					"enum":        []string{"create", "read", "update", "delete", "login", "logout"},
					"description": "Type of action performed",
				},
				"timestamp": bson.M{
					"bsonType":    "date",
					"description": "When the action occurred",
				},
				"resource_type": bson.M{
					"bsonType":    "string",
					"description": "Type of resource accessed",
				},
				"resource_id": bson.M{
					"bsonType":    "string",
					"description": "ID of the resource accessed",
				},
				"ip_address": bson.M{
					"bsonType":    "string",
					"description": "IP address of the client",
				},
				"user_agent": bson.M{
					"bsonType":    "string",
					"description": "User agent string",
				},
				"metadata": bson.M{
					"bsonType":    "object",
					"description": "Additional metadata about the action",
				},
			},
		},
	}

	// Create collection with validation
	opts := options.CreateCollection().SetValidator(validator)
	err := db.CreateCollection(ctx, "audit_logs", opts)
	if err != nil {
		return err
	}

	collection := db.Collection("audit_logs")

	// Create indexes for efficient querying
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().
				SetName("idx_audit_user_timestamp").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "resource_type", Value: 1},
				{Key: "resource_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().
				SetName("idx_audit_resource_timestamp").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "action", Value: 1},
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().
				SetName("idx_audit_action_timestamp").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "timestamp", Value: -1},
			},
			Options: options.Index().
				SetName("idx_audit_timestamp").
				SetBackground(true).
				SetExpireAfterSeconds(365 * 24 * 60 * 60), // Auto-delete after 1 year
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (m *CreateAuditCollectionMigration) Down(ctx context.Context, db *mongo.Database) error {
	// Drop the entire audit_logs collection
	return db.Collection("audit_logs").Drop(ctx)
}
