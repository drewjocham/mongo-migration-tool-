package main

import (
	"log"
	"os"

	"github.com/jocham/mongo-essential/cmd"
)

// Version information set by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version information for CLI
	cmd.SetVersion(version, commit, date)
	
	if err := cmd.Execute(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
