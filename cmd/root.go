// Package cmd provides the command-line interface for mongo-migrate.
// mongo-migrate is a comprehensive MongoDB migration and database analysis tool
// with AI-powered insights, similar to Liquibase/Flyway for MongoDB.
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jocham/mongo-essential/config"
	"github.com/jocham/mongo-essential/migration"
)

var (
	configFile string
	cfg        *config.Config
	db         *mongo.Database
	engine     *migration.Engine
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mongo-essential",
	Short: "Essential MongoDB toolkit with migrations and AI-powered analysis",
	Long: `A MongoDB migration tool that provides version control for your database schema.
	
Features:
- Version-controlled migrations with up/down support
- Migration status tracking
- Rollback capabilities
- Force migration marking
- Integration with existing Go projects`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if configFile != "" {
			cfg, err = config.Load(configFile)
		} else {
			cfg, err = config.Load(".env", ".env.local")
		}
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()

		clientOpts := options.Client().
			ApplyURI(cfg.GetConnectionString()).
			SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
			SetMinPoolSize(uint64(cfg.MinPoolSize)).
			SetMaxConnIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second).
			SetServerSelectionTimeout(time.Duration(cfg.Timeout) * time.Second).
			SetConnectTimeout(time.Duration(cfg.Timeout) * time.Second)

		if cfg.SSLEnabled {
			clientOpts.SetTLSConfig(&tls.Config{
				InsecureSkipVerify: cfg.SSLInsecure,
			})
		}

		client, err := mongo.Connect(ctx, clientOpts)
		if err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %w", err)
		}

		if err = client.Ping(ctx, nil); err != nil {
			return fmt.Errorf("failed to ping MongoDB: %w", err)
		}

		db = client.Database(cfg.Database)
		engine = migration.NewEngine(db, cfg.MigrationsCollection)

		if err := loadMigrations(); err != nil {
			return fmt.Errorf("failed to load migrations: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is .env)")

	// Add subcommands
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(forceCmd)
	rootCmd.AddCommand(versionCmd)
}

// loadMigrations loads migration files from the migrations directory
// This is a placeholder - in a real implementation you would scan the directory
// and dynamically load Go files or use a plugin system
func loadMigrations() error {
	// For now, we'll register migrations manually
	// In a production system, you might want to:
	// 1. Use go:embed to embed migration files
	// 2. Use a plugin system to load migration .so files
	// 3. Use a scripting language like Lua or JavaScript for migrations
	// 4. Parse Go files and use go/build to compile them

	fmt.Printf("Looking for migrations in: %s\n", cfg.MigrationsPath)

	// Check if migrations directory exists
	if _, err := os.Stat(cfg.MigrationsPath); os.IsNotExist(err) {
		fmt.Printf("Migrations directory does not exist: %s\n", cfg.MigrationsPath)
		return nil
	}

	return nil
}
