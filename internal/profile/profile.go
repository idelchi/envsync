package profile

import (
	"fmt"
	"os"
	"slices"

	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

// Profiles is a map of profile names to their metadata.
type Profiles map[string]*Profile

// RawEnv is a map of environment variable names to their non-stringified values.
type RawEnv map[string]any

// Profile holds metadata plus an env-var map.
type Profile struct {
	Env     env.Env  `toml:"-"                 yaml:"-"`
	RawEnv  RawEnv   `toml:"env,omitempty"     yaml:"env,omitempty"`
	Extends []string `toml:"extends,omitempty" yaml:"extends,omitempty"`
}

// newProfile creates a new profile with an empty env-var map.
func newProfile() *Profile {
	return &Profile{
		Env:    make(env.Env),
		RawEnv: make(RawEnv),
	}
}

// ToEnv serialises p.rawEnv into p.Env using Stringify.
// – Scalars pass through unchanged.
// – Non-scalars are JSON-minified and single-quoted (see Stringify).
func (p *Profile) ToEnv() error {
	p.Env = make(env.Env, len(p.RawEnv))

	keys := make([]string, 0, len(p.RawEnv))
	for k := range p.RawEnv {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, key := range keys {
		v := p.RawEnv[key]

		str, err := Stringify(v)
		if err != nil {
			return fmt.Errorf("profile %q: %w", key, err)
		}

		if err := p.Env.AddPair(key, str); err != nil {
			return fmt.Errorf("profile %q: %w", key, err)
		}
	}

	return nil
}

// File returns the first file found in the given paths.
func File(paths ...string) (file.File, error) {
	file, ok := files.New("", paths...).Exists()
	if !ok {
		return file, fmt.Errorf("profile %w: %v", os.ErrNotExist, paths)
	}

	return file, nil
}
