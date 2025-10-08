package main

import (
	"log"
	"os"

	"github.com/jocham/mongo-essential/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
