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
func Execute(version string) error {
	root := &cobra.Command{
		Use:           "envprof",
		Short:         "Manage env profiles in YAML/TOML with inheritance",
		Version:       version,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			// Do not print usage after basic validation has been done.
			cmd.SilenceUsage = true
		},
	}

	root.SetVersionTemplate("{{ .Version }}\n")
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.Flags().SortFlags = false
	root.CompletionOptions.DisableDefaultCmd = true
	cobra.EnableCommandSorting = false

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
		List(flags),
		Export(flags),
		Shell(flags),
	)

	if err := root.Execute(); err != nil {
		return fmt.Errorf("envprof: %w", err)
	}

	return nil
}
