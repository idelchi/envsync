package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/env"
)

// load loads the profile store from the specified file and fallbacks.
func load(flags *Flags) (profile.Profiles, error) {
	profiles, err := profile.New(flags.File...)
	if err != nil {
		return nil, fmt.Errorf("new profile: %w", err)
	}

	store, err := profiles.Load()
	if err != nil {
		return nil, fmt.Errorf("loading profile: %w", err)
	}

	return store.Profiles, nil
}

// shell represents a shell environment.
type shell struct {
	// Name is the name of the shell to spawn.
	Name string
}

// newShell creates a new shell instance with the specified name.
// It defaults to the current shell if no name is provided.
func newShell(name string) *shell {
	shell := &shell{}

	shell.set(name)

	return shell
}

// spawnShell spawns a new shell with the specified environment variables.
func (s *shell) Spawn(env []string) error {
	cmd := exec.Command(s.Name) //nolint:gosec	// Aware of the security implications of using exec.Command.
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("spawning shell %q: %w", s.Name, err)
	}

	return nil
}

// set tries to determine the current shell being used.
func (s *shell) set(shell string) {
	if shell != "" {
		s.Name = shell

		return
	}

	env := env.FromEnv()

	if shell := env.GetAny("SHELL", "STARSHIP_SHELL"); shell != "" {
		s.Name = shell

		return
	}

	switch runtime.GOOS {
	case "windows":
		s.Name = env.Get("ComSpec")
		if s.Name == "" {
			s.Name = "cmd.exe"
		}
	default:
		s.Name = "sh"
	}
}
