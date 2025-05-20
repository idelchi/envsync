package cli

import (
	"fmt"

	"github.com/idelchi/envprof/internal/profile"
	"github.com/idelchi/godyl/pkg/path/files"
)

// load loads the profile store from the specified file and fallbacks.
func load(paths []string) (profile.Profiles, error) {
	file, ok := files.New("", paths...).Exists()
	if !ok {
		//nolint:err113	// Occasional dynamic errors are fine.
		return nil, fmt.Errorf("profile file not found: searched for %v", paths)
	}

	profiles, err := profile.New(file)
	if err != nil {
		return nil, err //nolint:wrapcheck	// Error does not need additional wrapping.
	}

	store, err := profiles.Load()
	if err != nil {
		return nil, fmt.Errorf("loading profile from %s: %w", file.String(), err)
	}

	return store.Profiles, nil
}
