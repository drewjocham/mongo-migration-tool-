package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of mongo-migrate tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mongo-migrate version %s\n", version)
	},
}
