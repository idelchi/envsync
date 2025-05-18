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
		SilenceErrors: true,
		Version:       version,
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

	listCommand := List(flags)
	listCommand.GroupID = "core"

	exportCommand := Export(flags)
	exportCommand.GroupID = "core"

	convertCommand := Convert(flags)
	convertCommand.GroupID = "destructive"

	importCmd := Import(flags)
	importCmd.GroupID = "destructive"

	root.AddCommand(
		listCommand,
		exportCommand,
		convertCommand,
		importCmd,
	)

	root.AddGroup(
		&cobra.Group{
			ID:    "core",
			Title: "Core commands",
		},
		&cobra.Group{
			ID:    "destructive",
			Title: "Destructive commands that reformat the profiles file",
		})

	if err := root.Execute(); err != nil {
		return fmt.Errorf("envprof: %w", err)
	}

	return nil
}
