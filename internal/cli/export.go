package cli

import (
	"fmt"
	"sort"

	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Export returns the cobra command for exporting a profile's variables.
func Export() *cobra.Command {
	return &cobra.Command{
		Use:   "export <profile>",
		Short: "Emit export KEY=VAL lines",
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
				fmt.Printf("export %s=%q\n", k, vars[k])
			}
			return nil
		},
	}
}
