package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// forceCmd represents the force command
var forceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force mark a migration as applied without running it",
	Long: `Force mark a specific migration as applied in the database without
actually executing the migration logic.

This is useful in scenarios such as:
- A migration was applied manually outside of the migration tool
- You need to fix migration state after a partial failure
- You're migrating from another migration system

WARNING: Use this command with caution as it doesn't verify that the 
migration was actually applied correctly.

Examples:
  mongo-migrate force 20231201_001   # Mark version 20231201_001 as applied`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		version := args[0]
		ctx := context.Background()

		fmt.Printf("WARNING: You are about to force mark migration %s as applied.\n", version)
		fmt.Println("This will NOT execute the migration logic.")
		fmt.Print("Are you sure you want to continue? (y/N): ")

		var response string
		_, _ = fmt.Scanln(&response)

		if response != "y" && response != "Y" && response != "yes" && response != "YES" {
			fmt.Println("Operation cancelled.")
			return nil
		}

		return engine.Force(ctx, version)
	},
}
