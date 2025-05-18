package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profile"
)

// Convert returns the cobra command for converting the config format.
func Convert(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "convert <yaml|toml>",
		Short: "Convert the current config to another format",
		Long: heredoc.Doc(`
			Convert the current config to another format.
			Will write to a file with the same name as the current config file, but with the new extension.

			Be aware, that this will not preserve comments or order of the keys.
			`),
		Aliases: []string{"conv"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			store, err := load(flags)
			if err != nil {
				return err
			}

			extension := args[0]

			store.File = store.File.WithExtension(extension)
			store.Type = profile.Type(extension)

			if err := store.Save(); err != nil {
				return fmt.Errorf("saving profile: %w", err)
			}

			return nil
		},
	}
}
