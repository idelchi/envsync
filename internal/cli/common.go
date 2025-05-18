package cli

import (
	"fmt"

	"github.com/idelchi/envprof/internal/profile"
)

// load loads the profile store from the specified file and fallbacks.
func load(flags *Flags) (*profile.Store, error) {
	profiles, err := profile.New(flags.File...)
	if err != nil {
		return nil, fmt.Errorf("new profile: %w", err)
	}

	store, err := profiles.Load()
	if err != nil {
		return nil, fmt.Errorf("loading profile: %w", err)
	}

	return store, nil
}
