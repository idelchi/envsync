package profile

import (
	"errors"
	"fmt"
	"slices"

	"github.com/idelchi/godyl/pkg/dag"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Type represents the type of the profile file.
type Type string

const (
	// YAML is the YAML file type.
	YAML Type = "yaml"
	// TOML is the TOML file type.
	TOML Type = "toml"
)

// Store represents the entire file and codec information.
type Store struct {
	// File is the file containing the profiles.
	File file.File
	// Profiles is a map of profile names to their metadata.
	Profiles Profiles
	// Type is the type of the file.
	Type Type
}

var (
	// ErrUnsupportedFileType is returned when the file type is not supported.
	ErrUnsupportedFileType = errors.New("unsupported file type")
	// ErrProfileNotFound is returned when a profile is not found in the store.
	ErrProfileNotFound = errors.New("profile not found")
)

// New creates a new store from the given path and fallbacks.
// If path is empty, it will look for the first file in fallbacks.
func New(path ...string) (*Store, error) {
	file, err := File(path...)
	if err != nil {
		return nil, err
	}

	store := &Store{
		File:     file,
		Profiles: map[string]*Profile{},
	}

	switch ext := file.Extension(); ext {
	case "yaml", "yml":
		store.Type = YAML
	case "toml":
		store.Type = TOML
	default:
		return nil, fmt.Errorf("%w: %q: %q", ErrUnsupportedFileType, file.Path(), ext)
	}

	return store, nil
}

// Load reads the file and unmarshals it into the store.
func (s *Store) Load() (*Store, error) {
	data, err := s.File.Read()
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", s.File.Path(), err)
	}

	if err = s.unmarshal(data); err != nil {
		return nil, fmt.Errorf("unmarshalling file %q: %w", s.File.Path(), err)
	}

	for name, profile := range s.Profiles {
		if err = profile.ToEnv(); err != nil {
			return nil, fmt.Errorf("profile %q: converting to env: %w", name, err)
		}
	}

	return s, err
}

// Save writes the store to disk.
func (s *Store) Save() error {
	var (
		out []byte
		err error
	)

	out, err = s.marshal()
	if err != nil {
		return fmt.Errorf("marshalling file %q: %w", s.File.Path(), err)
	}

	if err := s.File.Write(out); err != nil {
		return fmt.Errorf("writing file %q: %w", s.File.Path(), err)
	}

	return nil
}

// Exists checks if a profile exists in the store.
func (s *Store) Exists(name string) bool {
	_, ok := s.Profiles[name]

	return ok
}

// RawEnv returns the raw environment variables for a profile.
func (s *Store) RawEnv(name string) (*RawEnv, error) {
	profile, ok := s.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrProfileNotFound, name)
	}

	return &profile.RawEnv, nil
}

// Resolved returns the merged environment variables for a profile, resolving dependencies.
func (s *Store) Resolved(name string) (*Env, error) {
	if !s.Exists(name) {
		return nil, fmt.Errorf("%w: %q", ErrProfileNotFound, name)
	}

	// Build dependency DAG.
	nodes := make([]string, 0, len(s.Profiles))
	for n := range s.Profiles {
		nodes = append(nodes, n)
	}

	g, err := dag.Build(nodes, func(n string) []string { return s.Profiles[n].Extends })
	if err != nil {
		return nil, fmt.Errorf("dag: %w", err)
	}

	chain, err := g.Chain(name)
	if err != nil {
		return nil, fmt.Errorf("chain: %w", err)
	}

	env := &Env{
		Env:         s.Profiles[name].Env,
		Inheritance: make(Inheritance),
	}

	if len(chain) == 1 {
		return env, nil
	}

	slices.Reverse(chain)

	// Remove the first element, as it is the profile itself.
	chain = chain[1:]

	errs := []error{}

	for _, name := range chain {
		m := s.Profiles[name].Env

		for k, v := range m {
			if !env.Env.Exists(k) {
				errs = append(errs, env.Env.AddPair(k, v))
				env.Inheritance[k] = name
			}
		}
	}

	return env, errors.Join(errs...)
}

// ProfilesSorted returns the profile names in sorted order.
func (s *Store) ProfilesSorted() []string {
	out := make([]string, 0, len(s.Profiles))
	for k := range s.Profiles {
		out = append(out, k)
	}

	slices.Sort(out)

	return out
}
