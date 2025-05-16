package cli

import (
	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Remove returns the cobra command for deleting a variable.
func Remove() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <profile> <key>",
		Short: "Delete variable",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, a []string) error {
			s, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}
			if len(a) == 1 {
				return s.Delete(a[0])
			}

			return s.RemoveVar(a[0], a[1])
		},
	}
}
