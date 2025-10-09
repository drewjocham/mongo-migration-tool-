package examplemigrations

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TransformUserDataMigration demonstrates data transformation operations
type TransformUserDataMigration struct{}

func (m *TransformUserDataMigration) Version() string {
	return "20240101_002"
}

func (m *TransformUserDataMigration) Description() string {
	return "Transform user data: normalize email case, add full_name field, and update timestamps"
}

func (m *TransformUserDataMigration) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// Find all users and transform their data
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			return err
		}

		// Prepare update operations
		update := bson.M{"$set": bson.M{}}

		// Normalize email to lowercase
		if email, exists := user["email"].(string); exists {
			update["$set"].(bson.M)["email"] = strings.ToLower(email)
		}

		// Create full_name from first_name and last_name
		firstName, hasFirst := user["first_name"].(string)
		lastName, hasLast := user["last_name"].(string)
		if hasFirst || hasLast {
			fullName := strings.TrimSpace(firstName + " " + lastName)
			if fullName != "" {
				update["$set"].(bson.M)["full_name"] = fullName
			}
		}

		// Add updated_at timestamp if it doesn't exist
		if _, hasUpdated := user["updated_at"]; !hasUpdated {
			update["$set"].(bson.M)["updated_at"] = time.Now()
		}

		// Only update if we have changes to make
		if len(update["$set"].(bson.M)) > 0 {
			filter := bson.M{"_id": user["_id"]}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	return cursor.Err()
}

func (m *TransformUserDataMigration) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// Remove the full_name field that we added
	update := bson.M{
		"$unset": bson.M{
			"full_name": "",
		},
	}

	_, err := collection.UpdateMany(ctx, bson.D{}, update)
	return err
}
