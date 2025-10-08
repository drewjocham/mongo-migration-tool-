package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long: `Show the status of all migrations, including which ones have been applied
and when they were applied.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		status, err := engine.GetStatus(ctx)
		if err != nil {
			return err
		}

		if len(status) == 0 {
			fmt.Println("No migrations found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "VERSION\tSTATUS\tAPPLIED AT\tDESCRIPTION")
		fmt.Fprintln(w, "-------\t------\t----------\t-----------")

		for _, migration := range status {
			status := "PENDING"
			appliedAt := ""

			if migration.Applied {
				status = "APPLIED"
				if migration.AppliedAt != nil {
					appliedAt = migration.AppliedAt.Format(time.RFC3339)
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				migration.Version,
				status,
				appliedAt,
				migration.Description,
			)
		}

		return w.Flush()
	},
}
