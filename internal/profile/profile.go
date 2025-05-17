// Package profile provides profile and store management for environment variable sets.
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

// Profile holds metadata plus an env-var map.
type Profile struct {
	Env     env.Env  `toml:"-" yaml:"-" `
	Extends []string `toml:"extends,omitempty" yaml:"extends,omitempty"`

	rawEnv map[string]any `toml:"env,omitempty" yaml:"env,omitempty" `
}

// ToEnv serialises p.rawEnv into p.Env using Stringify.
// – Scalars pass through unchanged.
// – Non-scalars are JSON-minified and single-quoted (see Stringify).
func (p *Profile) ToEnv() error {
	if p.Env == nil {
		p.Env = make(env.Env, len(p.rawEnv)) // guarantee non-nil target
	}

	// Stable output: sort keys once.
	keys := make([]string, 0, len(p.rawEnv))
	for k := range p.rawEnv {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, k := range keys {
		v := p.rawEnv[k]

		str, err := Stringify(v)
		if err != nil {
			return fmt.Errorf("profile %q: %w", k, err) // key first, then wrapped err
		}
		p.Env.AddPair(k, str)
	}

	return nil
}

func ProfileFile(path string, fallbacks ...string) (file.File, error) {
	if path != "" {
		file := file.New(path)
		if !file.Exists() {
			return file, fmt.Errorf("profile %w: %q", os.ErrNotExist, path)
		}

		return file, nil
	}

	file, ok := files.New("", fallbacks...).Exists()
	if !ok {
		return file, fmt.Errorf("profile %w: %v", os.ErrNotExist, fallbacks)
	}

	return file, nil
}
