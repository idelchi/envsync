package profile

import (
	"fmt"
	"slices"

	"github.com/idelchi/godyl/pkg/env"
)

// Env is a map of environment variable names to their non-stringified values.
type Env map[string]any

// Stringified serializes Env into env.Env using Stringify.
// – Scalars pass through unchanged.
// – Non-scalars are JSON-minified and single-quoted (see Stringify).
func (e Env) Stringified() (env.Env, error) {
	env := make(env.Env, len(e))

	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, key := range keys {
		str, err := Stringify(e[key])
		if err != nil {
			return nil, fmt.Errorf("profile %q: %w", key, err)
		}

		if err := env.AddPair(key, str); err != nil {
			return nil, fmt.Errorf("profile %q: %w", key, err)
		}
	}

	return env, nil
}
