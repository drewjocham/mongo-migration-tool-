package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long: `Display the current status of all migrations, showing which have been applied
and which are pending.

This command shows:
- Migration version and description
- Applied status (✓ or ✗)  
- Timestamp when applied (if applicable)

Examples:
  mongo-essential status
  mongo-essential status --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fmt.Println("Migration Status")
		fmt.Println(strings.Repeat("=", 50))

		status, err := engine.GetStatus(ctx)
		if err != nil {
			return fmt.Errorf("failed to get migration status: %w", err)
		}

		if len(status) == 0 {
			fmt.Println("No migrations found")
			return nil
		}

		for _, s := range status {
			statusIcon := "✗"
			statusText := "Pending"
			appliedAt := ""

			if s.Applied {
				statusIcon = "✓"
				statusText = "Applied"
				if s.AppliedAt != nil {
					appliedAt = fmt.Sprintf(" (%s)", s.AppliedAt.Format("2006-01-02 15:04:05"))
				}
			}

			fmt.Printf("%s %s - %s %s%s\n",
				statusIcon,
				s.Version,
				s.Description,
				statusText,
				appliedAt)
		}

		return nil
	},
}
