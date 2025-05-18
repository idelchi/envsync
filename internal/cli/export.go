package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Export returns the cobra command for exporting a profile's variables.
//
//nolint:forbidigo	// Command print out to the console.
func Export(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "export <profile> [file]",
		Short: "Emit export KEY=VAL lines or write out as dotenv file",
		Long: heredoc.Doc(`
			Emit export KEY=VAL lines or write out as dotenv file.
			If no file is provided, the output will be printed to stdout as:

			$ export KEY=VAL
			$ export KEY2=VAL2
			$ export KEY3=VAL3

			If a file is provided, the output will be written to that file as:

			KEY=VAL
			KEY2=VAL2
			KEY3=VAL3
			`),
		Aliases: []string{"x"},
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

			if len(args) == 1 {
				envs := vars.FormatAll("export", false)

				fmt.Println(envs)

				return nil
			}

			envs := vars.Env.AsSlice()

			dotenv := file.New(args[1])
			envs = append([]string{fmt.Sprintf("# Active profile: %q", prof)}, envs...)
			if err := dotenv.Write([]byte(strings.Join(envs, "\n"))); err != nil {
				return fmt.Errorf("write dotenv: %w", err)
			}

			return nil
		},
	}
}
