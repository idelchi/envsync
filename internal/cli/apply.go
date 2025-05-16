package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/envsync/internal/profile"
	"github.com/spf13/cobra"
)

// Apply returns the cobra command for applying a profile.
func Apply() *cobra.Command {
	return &cobra.Command{
		Use:   "apply <profile> [outfile]",
		Short: "Print KEY=VAL lines & optionally write file",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, a []string) error {
			s, err := profile.Load(fileFlag)
			if err != nil {
				return err
			}
			vars, err := s.Vars(a[0])
			if err != nil {
				return err
			}
			var b strings.Builder
			for k, v := range vars {
				line := fmt.Sprintf("%s=%s\n", k, v)
				b.WriteString(line)
				os.Setenv(k, v)
			}
			out := b.String()
			fmt.Print(out)
			if len(a) == 2 {
				if err := os.WriteFile(a[1], []byte(out), 0o644); err != nil {
					return err
				}
				fmt.Println("written to", a[1])
			}
			return nil
		},
	}
}
