package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/pkg/env"
)

// Shell returns the cobra command for entering a scoped shell with the environment active.
//
//nolint:forbidigo	// Command print out to the console.
func Shell(flags *Flags) *cobra.Command {
	env := env.FromEnv()

	var (
		shell   = env.GetAny("SHELL", "STARSHIP_SHELL")
		inherit bool
	)

	cmd := &cobra.Command{
		Use:   "shell <profile>",
		Short: "Spawn a new shell with the profile's environment",
		Long: heredoc.Doc(`
      Spawn a new shell with the profile's environment.
      The shell will be spawned with the environment variables set to the profile's values.

      Use the --shell flag to specify the shell to spawn, or it will default to the value of $SHELL or $STARSHIP_SHELL,
      or if those around found, to the value of $ComSpec or cmd.exe on Windows, and sh on all other platforms.

      Use the --inherit flag to inherit the environment variables from the parent shell.
    `),
		Example: heredoc.Doc(`
      # Spawn a new shell with the profile's environment
      $ godyl shell dev
      # Spawn a new shell with the profile's environment and inherit the environment variables from the parent shell
      $ godyl shell dev --inherit
      # Spawn a new shell with the profile's environment and use zsh as the shell
      $ godyl shell dev --shell zsh
      `),
		Aliases: []string{"sh"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if active := env.Get("ENVPROF_ACTIVE_SHELL"); env.Exists("ENVPROF_ACTIVE_SHELL") {
				//nolint:err113	// Occasional dynamic errors are fine.
				return fmt.Errorf("already inside profile %q, exit first", active)
			}

			profiles, err := load(flags)
			if err != nil {
				return err
			}

			terminal := newShell(shell)

			prof := args[0]

			vars, err := profiles.Environment(prof)
			if err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if err = vars.Env.AddPair("ENVPROF_ACTIVE_SHELL", prof); err != nil {
				return err //nolint:wrapcheck	// Error does not need additional wrapping.
			}

			if inherit {
				vars.Env.Merge(env)
			}

			fmt.Printf("Entering shell %q with profile %q...\n", shell, prof)

			if err := terminal.Spawn(vars.Env.AsSlice()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&shell, "shell", "s", shell, "Shell to spawn (default: $SHELL or $STARSHIP_SHELL)")
	cmd.Flags().BoolVarP(&inherit, "inherit", "i", false, "Inherit environment variables from the parent shell")

	return cmd
}
