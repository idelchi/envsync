package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Export defines the command for exporting a profile's variables.
// It emits '<prefix> KEY=VAL' lines or writes them out as a dotenv file.
// By default, 'prefix' is set to "export".
//
//nolint:forbidigo	// Command print out to the console.
func Export(files *[]string) *cobra.Command {
	prefix := "export "

	cmd := &cobra.Command{
		Use:   "export <profile> [file]",
		Short: "Emit 'export KEY=VAL' lines or write out as dotenv file",
		Long: heredoc.Doc(`
			Emit 'export KEY=VAL' lines or write out as dotenv file.
			If no file is provided, the output will be printed to stdout as:

			export KEY1=VAL1
			export KEY2=VAL2
			export KEY3=VAL3

			The --prefix flag can be used to set an alternative prefix instead of 'export '.

			If a file is provided, the output will be written to that file as:

			KEY1=VAL1
			KEY2=VAL2
			KEY3=VAL3

			The --prefix flag is ignored when writing to a file.
			`),
		Example: heredoc.Doc(`
			# Emit 'export KEY=VAL' lines
			$ envprof export dev

			# Emit '$env:KEY=VAL' lines with a custom prefix
			$ envprof export dev --prefix "$env:"

			# Write out as dotenv file
			$ envprof export dev .env
		`),

		Aliases: []string{"x"},
		Args:    cobra.RangeArgs(1, 2), //nolint:mnd	// The command takes 1 or 2 arguments as documented.
		RunE: func(_ *cobra.Command, args []string) error {
			profiles, err := load(*files)
			if err != nil {
				return err
			}

			prof := args[0]

			vars, err := profiles.Environment(prof)
			if err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if len(args) == 1 {
				envs := vars.FormatAll(prefix, false)

				fmt.Println(envs)

				return nil
			}

			envs := vars.Env.AsSlice()

			dotenv := file.New(args[1])
			envs = append([]string{fmt.Sprintf("# Active profile: %q", prof)}, envs...)
			envs = append(envs, "")
			if err := dotenv.Write([]byte(strings.Join(envs, "\n"))); err != nil {
				return fmt.Errorf("write dotenv: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&prefix, "prefix", "p", prefix, "Prefix for the export command")

	return cmd
}
