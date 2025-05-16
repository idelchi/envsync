// Command envsync is a command-line tool for managing environment variables in profiles.
package main

import (
	"log"

	"github.com/idelchi/envsync/internal/cli"
)

// main is the entry point for the CLI application.
func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
