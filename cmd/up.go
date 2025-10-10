package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	upTargetVersion string
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations (or up to target version)",
	Long: `Run pending migrations in forward direction.

By default, all pending migrations are executed in version order.
You can optionally specify a target version to migrate up to a specific version.

Examples:
  mongo-essential up                    # Run all pending migrations
  mongo-essential up --target 20231201_002  # Run migrations up to specific version`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx := context.Background()

		if upTargetVersion != "" {
			fmt.Printf("Running migrations up to version: %s\n", upTargetVersion)
		} else {
			fmt.Println("Running all pending migrations...")
		}

		if err := engine.Up(ctx, upTargetVersion); err != nil {
			return fmt.Errorf("migration up failed: %w", err)
		}

		fmt.Println("✓ Migrations completed successfully!")
		return nil
	},
}

func setupUpCommand() {
	upCmd.Flags().StringVar(&upTargetVersion, "target", "", "Target version to migrate up to")
}
