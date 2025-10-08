package migrations

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Migration_20231201_002 adds default data to collections
type Migration_20231201_002 struct{}

// Version returns the unique version identifier for this migration
func (m *Migration_20231201_002) Version() string {
	return "20231201_002"
}

// Description returns a human-readable description of what this migration does
func (m *Migration_20231201_002) Description() string {
	return "Add default admin user and system configuration"
}

// Up executes the migration
func (m *Migration_20231201_002) Up(ctx context.Context, db *mongo.Database) error {
	// Add default admin user
	usersCollection := db.Collection("users")

	adminUser := bson.D{
		{"_id", "admin-001"},
		{"username", "admin"},
		{"email", "admin@example.com"},
		{"role", "admin"},
		{"status", "active"},
		{"created_at", time.Now().UTC()},
		{"updated_at", time.Now().UTC()},
	}

	_, err := usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	// Add system configuration
	configCollection := db.Collection("system_config")

	configs := []interface{}{
		bson.D{
			{"_id", "app_settings"},
			{"max_upload_size", 10485760}, // 10MB
			{"allowed_file_types", []string{"jpg", "png", "pdf", "docx"}},
			{"maintenance_mode", false},
			{"created_at", time.Now().UTC()},
		},
		bson.D{
			{"_id", "email_settings"},
			{"smtp_host", "localhost"},
			{"smtp_port", 587},
			{"from_email", "noreply@example.com"},
			{"created_at", time.Now().UTC()},
		},
	}

	_, err = configCollection.InsertMany(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to insert system configuration: %w", err)
	}

	fmt.Println("Added default admin user and system configuration")
	return nil
}

// Down rolls back the migration
func (m *Migration_20231201_002) Down(ctx context.Context, db *mongo.Database) error {
	// Remove the admin user
	usersCollection := db.Collection("users")
	_, err := usersCollection.DeleteOne(ctx, bson.D{{"_id", "admin-001"}})
	if err != nil {
		return fmt.Errorf("failed to remove admin user: %w", err)
	}

	// Remove system configuration
	configCollection := db.Collection("system_config")
	_, err = configCollection.DeleteMany(ctx, bson.D{
		{"_id", bson.D{{"$in", []string{"app_settings", "email_settings"}}}},
	})
	if err != nil {
		return fmt.Errorf("failed to remove system configuration: %w", err)
	}

	fmt.Println("Removed default admin user and system configuration")
	return nil
}
