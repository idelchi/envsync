// Package cli implements the command-line interface for envprof.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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

	files := &[]string{
		"envprof.yaml",
		"envprof.yml",
		"envprof.toml",
	}

	if file := os.Getenv("ENVPROF_FILE"); file != "" {
		files = &[]string{file}
	}

	root.PersistentFlags().
		StringSliceVarP(files, "file", "f", *files, "config file to use, in order of preference")

	root.AddCommand(
		List(files),
		Export(files),
		Shell(files),
	)

	if err := root.Execute(); err != nil {
		return fmt.Errorf("envprof: %w", err)
	}

	return nil
}
