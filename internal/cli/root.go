// Package cli implements the command-line interface for envprof.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flags holds the command-line flags for the envprof CLI application.
type Flags struct {
	// File is the config file to use, in order of preference.
	File []string
	// Verbose enables verbose output.
	Verbose bool
}

// Execute runs the root command for the envprof CLI application.
func Execute() error {
	root := &cobra.Command{
		Use:           "envprof",
		Short:         "Manage env profiles in YAML/TOML/JSON with inheritance",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := &Flags{}

	defaultFiles := []string{
		"envprof.yaml",
		"envprof.yml",
		"envprof.toml",
	}

	if file := os.Getenv("ENVPROF_FILE"); file != "" {
		defaultFiles = []string{file}
	}

	root.PersistentFlags().
		StringSliceVarP(&flags.File, "file", "f", defaultFiles, "config file to use, in order of preference")
	root.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "enable verbose output")

	root.AddCommand(
		Get(flags),
		Export(flags),
		Profiles(flags), Convert(flags),
		Import(flags),
	)

	if err := root.Execute(); err != nil {
		return fmt.Errorf("envprof: %w", err)
	}

	return nil
}
