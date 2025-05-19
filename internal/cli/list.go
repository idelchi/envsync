package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// List returns the cobra command for listing profiles and their variables.
//
//nolint:forbidigo	// Command print out to the console.
func List(files *[]string) *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "list <profile> [key]",
		Short: "List profiles and their variables",
		Long: heredoc.Doc(`
		Calling this function with:
			- no arguments lists all available profiles (alphabetically sorted).
			- with a profile name lists all variables for that profile.
			- with a profile name and a key lists the value of that key for that profile.
		`),
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(2), //nolint:mnd	// The command takes up to 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			profiles, err := load(*files)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				for _, profile := range profiles.Names() {
					fmt.Println(profile)
				}

				return nil
			}

			prof := args[0]

			vars, err := profiles.Environment(prof)
			if err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if len(args) > 1 {
				if !vars.Env.Exists(args[1]) {
					//nolint:err113	// Occasional dynamic errors are fine.
					return fmt.Errorf("key %q not found in profile %q", args[1], prof)
				}
				fmt.Println(vars.Format(args[1], verbose, false))
			} else {
				fmt.Println(vars.FormatAll("", verbose))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show from which source each variable is inherited")

	return cmd
}
