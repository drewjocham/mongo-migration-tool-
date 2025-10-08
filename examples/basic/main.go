package main

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jocham/mongo-essential/config"
	"github.com/jocham/mongo-essential/migration"
)

func main() {
	// Load configuration from environment
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	// Connect to MongoDB
	clientOpts := options.Client().
		ApplyURI(cfg.GetConnectionString()).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MinPoolSize))

	if cfg.SSLEnabled {
		clientOpts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: cfg.SSLInsecure,
		})
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Create migration engine
	db := client.Database(cfg.Database)
	engine := migration.NewEngine(db, cfg.MigrationsCollection)

	// Register migrations
	engine.RegisterMany(
		&CreateUsersCollection{},
		&AddUserIndexes{},
	)

	// Get current migration status
	status, err := engine.GetStatus(ctx)
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	log.Println("Migration Status:")
	for _, s := range status {
		applied := "❌"
		if s.Applied {
			applied = "✅"
		}
		log.Printf("  %s %s - %s", applied, s.Version, s.Description)
	}

	// Run pending migrations
	log.Println("Running pending migrations...")
	if err := engine.Up(ctx, ""); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("All migrations completed successfully!")
}

// Example migration: Create users collection
type CreateUsersCollection struct{}

func (m *CreateUsersCollection) Version() string {
	return "20240101_001"
}

func (m *CreateUsersCollection) Description() string {
	return "Create users collection with basic structure"
}

func (m *CreateUsersCollection) Up(ctx context.Context, db *mongo.Database) error {
	log.Println("Creating users collection...")

	// Create the collection (MongoDB creates collections automatically on first insert,
	// but we can explicitly create it to set options)
	opts := options.CreateCollection()
	err := db.CreateCollection(ctx, "users", opts)
	if err != nil {
		// Collection might already exist, which is fine
		log.Printf("Collection creation note: %v", err)
	}

	// Insert a sample document to ensure the collection exists
	collection := db.Collection("users")
	_, err = collection.InsertOne(ctx, bson.M{
		"email":      "admin@drewjocham.com",
		"username":   "admin",
		"created_at": time.Now(),
		"status":     "active",
	})

	return err
}

func (m *CreateUsersCollection) Down(ctx context.Context, db *mongo.Database) error {
	log.Println("Dropping users collection...")
	return db.Collection("users").Drop(ctx)
}

// Example migration: Add indexes to users collection
type AddUserIndexes struct{}

func (m *AddUserIndexes) Version() string {
	return "20240101_002"
}

func (m *AddUserIndexes) Description() string {
	return "Add email and username indexes to users collection"
}

func (m *AddUserIndexes) Up(ctx context.Context, db *mongo.Database) error {
	log.Println("Adding indexes to users collection...")

	collection := db.Collection("users")

	// Create unique index on email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}
	
	// Create unique index on username
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_unique"),
	}
	
	// Create compound index on created_at and status
	compoundIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: -1}, {Key: "status", Value: 1}},
		Options: options.Index().SetName("created_at_status"),
	}

	// Create all indexes
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		emailIndex,
		usernameIndex,
		compoundIndex,
	})

	return err
}

func (m *AddUserIndexes) Down(ctx context.Context, db *mongo.Database) error {
	log.Println("Dropping indexes from users collection...")

	collection := db.Collection("users")

	// Drop specific indexes
	indexes := []string{"email_unique", "username_unique", "created_at_status"}
	for _, indexName := range indexes {
		_, err := collection.Indexes().DropOne(ctx, indexName)
		if err != nil {
			log.Printf("Warning: failed to drop index %s: %v", indexName, err)
		}
	}

	return nil
}
