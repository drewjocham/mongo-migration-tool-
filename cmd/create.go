package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <description>",
	Short: "Create a new migration file",
	Long: `Create a new migration file with the given description.
	
The migration file will be created with a timestamp-based version number
and the provided description. The file will contain a template with
Up and Down functions.

Examples:
  mongo-migrate create "add user collection"
  mongo-migrate create "create index on products"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		description := args[0]

		// Generate version based on timestamp
		version := time.Now().Format("20060102_150405")

		// Clean description for filename
		cleanDesc := ""
		for _, char := range description {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
				cleanDesc += string(char)
			} else if char == ' ' || char == '-' || char == '_' {
				cleanDesc += "_"
			}
		}

		filename := fmt.Sprintf("%s_%s.go", version, cleanDesc)
		filepath := filepath.Join(cfg.MigrationsPath, filename)

		// Create migrations directory if it doesn't exist
		if err := os.MkdirAll(cfg.MigrationsPath, 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %w", err)
		}

		// Check if file already exists
		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			return fmt.Errorf("migration file already exists: %s", filepath)
		}

		// Generate migration template
		template := fmt.Sprintf(`package migrations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration_%s implements the Migration interface
type Migration_%s struct{}

// Version returns the unique version identifier for this migration
func (m *Migration_%s) Version() string {
	return "%s"
}

// Description returns a human-readable description of what this migration does
func (m *Migration_%s) Description() string {
	return "%s"
}

// Up executes the migration
func (m *Migration_%s) Up(ctx context.Context, db *mongo.Database) error {
	// TODO: Implement your migration logic here
	
	// Example: Create a collection
	// collection := db.Collection("your_collection_name")
	
	// Example: Create an index
	// indexModel := mongo.IndexModel{
	//     Keys: bson.D{{"field_name", 1}},
	//     Options: options.Index().SetUnique(true),
	// }
	// _, err := collection.Indexes().CreateOne(ctx, indexModel)
	// if err != nil {
	//     return fmt.Errorf("failed to create index: %%w", err)
	// }
	
	// Example: Insert documents
	// documents := []interface{}{
	//     bson.D{{"name", "example1"}, {"value", 1}},
	//     bson.D{{"name", "example2"}, {"value", 2}},
	// }
	// _, err := collection.InsertMany(ctx, documents)
	// if err != nil {
	//     return fmt.Errorf("failed to insert documents: %%w", err)
	// }
	
	return fmt.Errorf("migration not implemented yet")
}

// Down rolls back the migration (optional - can return ErrNotSupported)
func (m *Migration_%s) Down(ctx context.Context, db *mongo.Database) error {
	// TODO: Implement rollback logic here
	
	// Example: Drop collection
	// return db.Collection("your_collection_name").Drop(ctx)
	
	// Example: Drop index
	// collection := db.Collection("your_collection_name")
	// _, err := collection.Indexes().DropOne(ctx, "index_name")
	// return err
	
	return fmt.Errorf("rollback not implemented yet")
}
`, version, version, version, version, version, description, version, version)

		file, err := os.Create(filepath)
		if err != nil {
			return fmt.Errorf("failed to create migration file: %w", err)
		}
		defer file.Close()

		if _, err := file.WriteString(template); err != nil {
			return fmt.Errorf("failed to write migration template: %w", err)
		}

		fmt.Printf("Created migration file: %s\n", filepath)
		fmt.Printf("Migration version: %s\n", version)
		fmt.Println("\nNext steps:")
		fmt.Println("1. Edit the migration file and implement the Up() and Down() methods")
		fmt.Println("2. Register the migration in your main application")
		fmt.Println("3. Run 'mongo-migrate up' to apply the migration")

		return nil
	},
}
