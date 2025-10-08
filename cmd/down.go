package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down [target-version]",
	Short: "Roll back migrations (down to target version)",
	Long: `Roll back migrations in reverse order, down to a specific target version.
	
Examples:
  mongo-migrate down 20231201_001     # Roll back down to version 20231201_001
  
Note: Rolling back all migrations is not supported by default. 
You must specify a target version to prevent accidental data loss.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		target := args[0]

		fmt.Printf("Rolling back migrations down to target version: %s\n", target)

		return engine.Down(ctx, target)
	},
}
