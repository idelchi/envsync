package profile

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/env"
)

// Inheritance carries a mapping from environment variable names to where they are inherited from.
type Inheritance map[string]string

// Env represents the environment variables and their inheritance.
type Env struct {
	// Env holds the environment variables.
	Env env.Env
	// Inheritance holds the mapping of environment variable names to their source.
	Inheritance Inheritance
}

// Format returns the formatted value of a single variable.
func (e Env) Format(key string, verbose, withKey bool) string {
	val := e.Env.Get(key)
	if withKey {
		val = fmt.Sprintf("%v=%v", key, val)
	}

	if verbose {
		if src := e.Inheritance[key]; src != "" {
			return fmt.Sprintf("%-45v (inherited from %q)", val, src)
		}
	}

	return val
}

// FormatAll returns all variables, one per line.
func (e Env) FormatAll(prefix string, verbose bool) string {
	out := []string{}

	for k := range e.Env {
		s := e.Format(k, verbose, true)
		if prefix != "" {
			s = fmt.Sprintf("%s %s", prefix, s)
		}

		out = append(out, s)
	}

	return strings.Join(out, "\n")
}
