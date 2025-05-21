// Package cli implements the command-line interface for envprof.
package cli

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Execute runs the root command for the envprof CLI application.
func Execute(version string) error {
	root := &cobra.Command{
		Use:   "envprof",
		Short: "Manage env profiles in YAML/TOML with inheritance",
		Long: heredoc.Doc(`
			Manage env profiles in YAML/TOML with inheritance.
			Profiles are loaded from a config file, which can be specified with the --file flag.

			The tool will by default search for the following files in the current directory:
			- envprof.yaml
			- envprof.yml
			- envprof.toml

			Profiles can be listed, exported, and used to spawn a new shell with the profile's environment.

			Profiles can be inherited from other profiles and dotenv files, allowing for some flexibility.
		`),
		Example: heredoc.Doc(`
			# List the variables for the 'dev' profile
			$ envprof list dev -v

			# Create a dotenv file from a given profile
			$ envprof export dev dev.env

			# Eval the profile in the current shell
			$ eval "$(envprof export dev)"

			# Enter a new shell with the profile's environment
			$ envprof shell dev --shell zsh
		`),
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
