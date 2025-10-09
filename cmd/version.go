package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information
var (
	appVersion = "dev"
	appCommit  = "none"
	appDate    = "unknown"
)

// SetVersion sets the version information from main
func SetVersion(version, commit, date string) {
	appVersion = version
	appCommit = commit
	appDate = date
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of mongo-essential.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mongo-essential version %s\n", appVersion)
		fmt.Printf("  commit: %s\n", appCommit)
		fmt.Printf("  built: %s\n", appDate)
	},
}
