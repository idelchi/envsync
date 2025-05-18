// Command envprof is a command-line tool for managing environment variables in profiles.
package main

import (
	"fmt"

	"github.com/idelchi/envprof/internal/cli"
)

// main is the entry point for the CLI application.
func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
	}
}
