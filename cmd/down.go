package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	downTargetVersion string
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back migrations (down to target version)",
	Long: `Roll back applied migrations in reverse order.

You must specify a target version to rollback to. The target version itself
will remain applied (not rolled back).

Examples:
  mongo-essential down --target 20231201_001  # Rollback to version 20231201_001`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx := context.Background()

		if downTargetVersion == "" {
			return fmt.Errorf("target version is required for rollback")
		}

		fmt.Printf("Rolling back migrations to version: %s\n", downTargetVersion)

		if err := engine.Down(ctx, downTargetVersion); err != nil {
			return fmt.Errorf("migration down failed: %w", err)
		}

		fmt.Println("✓ Rollback completed successfully!")
		return nil
	},
}

func setupDownCommand() {
	downCmd.Flags().StringVar(&downTargetVersion, "target", "", "Target version to rollback to (required)")
	if err := downCmd.MarkFlagRequired("target"); err != nil {
		panic(fmt.Sprintf("failed to mark target flag as required: %v", err))
	}
}
