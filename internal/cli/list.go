package cli

import (
	"fmt"
	"sort"

	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// List returns the cobra command for listing variables of a profile.
func List() *cobra.Command {
	return &cobra.Command{
		Use:   "list <profile>",
		Short: "List variables (inherited)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			s, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}
			vars, err := s.Vars(a[0])
			if err != nil {
				return err
			}
			keys := make([]string, 0, len(vars))
			for k := range vars {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%s=%s\n", k, vars[k])
			}
			return nil
		},
	}
}
