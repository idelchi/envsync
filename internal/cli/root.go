// Package cli implements the command-line interface for envsync.
package cli

import "github.com/spf13/cobra"

var fileFlag string

// Execute runs the root command for the envsync CLI application.
func Execute() error {
	root := &cobra.Command{
		Use:   "envsync",
		Short: "Manage env profiles in YAML/TOML/JSON with inheritance",
	}
	root.PersistentFlags().StringVarP(&fileFlag, "file", "f", "", "config file to use")

	root.AddCommand(
		Add(), Get(), List(),
		Apply(), Export(),
		Remove(),
		Profiles(), Convert(),
	)

	return root.Execute()
}
