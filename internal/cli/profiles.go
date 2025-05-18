package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Profiles returns the cobra command for listing profile names.
//
//nolint:forbidigo	// Command print out to the console.
func Profiles(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:     "profiles",
		Short:   "List existing profile (alphabetically sorted)",
		Aliases: []string{"profs"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			store, err := load(flags)
			if err != nil {
				return err
			}

			for _, profile := range store.ProfilesSorted() {
				fmt.Println(profile)
			}

			return nil
		},
	}
}
