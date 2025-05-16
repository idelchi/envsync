package cli

import (
	"fmt"

	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Profiles returns the cobra command for listing profile names.
func Profiles() *cobra.Command {
	return &cobra.Command{
		Use:   "profiles",
		Short: "List profile names",
		RunE: func(_ *cobra.Command, _ []string) error {
			s, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}
			for _, n := range s.ProfilesSorted() {
				fmt.Println(n)
			}
			return nil
		},
	}
}
