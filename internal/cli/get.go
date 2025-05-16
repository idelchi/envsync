package cli

import (
	"fmt"

	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Get returns the cobra command for retrieving a variable's value.
func Get() *cobra.Command {
	return &cobra.Command{
		Use:   "get <profile> <key>",
		Short: "Print variable value (after inheritance)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, a []string) error {
			s, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}
			vars, err := s.Vars(a[0])
			if err != nil {
				return err
			}
			v, ok := vars[a[1]]
			if !ok {
				return fmt.Errorf("%s not found", a[1])
			}
			fmt.Println(v)
			return nil
		},
	}
}
