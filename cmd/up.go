package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [target-version]",
	Short: "Run all pending migrations (or up to target version)",
	Long: `Run all pending migrations in order, or up to a specific target version.
	
Examples:
  mongo-migrate up                    # Run all pending migrations
  mongo-migrate up 20231201_001       # Run migrations up to version 20231201_001`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		var target string
		if len(args) > 0 {
			target = args[0]
		}

		fmt.Printf("Running migrations up")
		if target != "" {
			fmt.Printf(" to target version: %s", target)
		}
		fmt.Println()

		return engine.Up(ctx, target)
	},
}
