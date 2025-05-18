package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Get returns the cobra command for retrieving a variable's value.
//
//nolint:forbidigo	// Command print out to the console.
func Get(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "get <profile> [key]",
		Short: "Get all or a specific variable from a profile",
		Long: heredoc.Doc(`
		Get all or a specific variable from a profile.
		If no key is provided, all variables will be printed.
		`),
		Aliases: []string{"ls"},
		Args:    cobra.RangeArgs(1, 2), //nolint:mnd	// The command takes 1 or 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			store, err := load(flags)
			if err != nil {
				return err
			}

			prof := args[0]

			vars, err := store.Resolved(prof)
			if err != nil {
				return fmt.Errorf("get vars: %w", err)
			}

			if len(args) > 1 {
				if !vars.Env.Exists(args[1]) {
					//nolint:err113	// Occasional dynamic errors are fine.
					return fmt.Errorf("key %q not found in profile %q", args[1], prof)
				}
				fmt.Println(vars.Format(args[1], flags.Verbose, false))
			} else {
				fmt.Println(vars.FormatAll("", flags.Verbose))
			}

			return nil
		},
	}
}
