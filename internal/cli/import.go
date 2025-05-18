package cli

import (
	"fmt"
	"maps"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/env"
)

// Import returns the cobra command for importing a dotenv file into a profile.
func Import(flags *Flags) *cobra.Command {
	var keep bool

	cmd := &cobra.Command{
		Use:   "import <profile> <dotenv> ",
		Short: "Import a dotenv file into a profile",
		Long: heredoc.Doc(`
			Import a dotenv file into a profile.
			The dotenv file will be merged with the profile's variables.
			If the '--keep' flag is set, the existing variables in the profile will be kept in case of conflicts.

			Be aware, that this will not preserve comments or order of the keys.
			`),
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(2), //nolint:mnd	// The command takes 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			store, err := load(flags)
			if err != nil {
				return err
			}

			prof := args[0]
			dotenvPath := args[1]

			dotenv, err := env.FromDotEnv(dotenvPath)
			if err != nil {
				return fmt.Errorf("loading dotenv: %w", err)
			}

			vars, err := store.RawEnv(prof)
			if err != nil {
				return fmt.Errorf("getting vars: %w", err)
			}

			dotenvAny := profile.ToRaw(dotenv)

			if keep {
				maps.Copy(dotenvAny, *vars)

				*vars = dotenvAny
			} else {
				maps.Copy(*vars, dotenvAny)
			}

			if err := store.Save(); err != nil {
				return fmt.Errorf("saving profile: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&keep, "keep", "k", false, "Keep existing variables in the profile in case of conflicts")

	return cmd
}
