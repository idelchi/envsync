// Package profile provides profile and store management for environment variable sets.
package profile

import (
	"fmt"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/pkg/dag"
	"github.com/idelchi/godyl/pkg/path/file"
)

type Type string

const (
	YAML Type = "yaml"
	TOML Type = "toml"
)

// Store represents the entire file and codec information.
type Store struct {
	File     file.File
	Profiles Profiles
	Type     Type
}

func New(path string, options ...string) (*Store, error) {
	file, err := ProfileFile(path, options...)
	if err != nil {
		return nil, err
	}

	var fileType Type
	switch ext := file.Extension(); ext {
	case "yaml", "yml":
		fileType = YAML
	case "toml":
		fileType = TOML
	default:
		return nil, fmt.Errorf("unsupported file extension %q", ext)
	}

	return &Store{File: file, Profiles: map[string]*Profile{}, Type: fileType}, nil
}

func (s *Store) marshal() ([]byte, error) {
	switch s.Type {
	case YAML:
		return yaml.MarshalWithOptions(s.Profiles)
	case TOML:
		return toml.Marshal(s.Profiles)
	default:
		return nil, fmt.Errorf("unsupported file type %q", s.Type)
	}
}

func (s *Store) unmarshal(data []byte) error {
	switch s.Type {
	case YAML:
		return yaml.UnmarshalWithOptions(data, s.Profiles, yaml.Strict())
	case TOML:
		md, err := toml.Decode(string(data), &s.Profiles)
		if err != nil {
			return err
		}

		if keys := md.Undecoded(); len(keys) > 0 {
			return fmt.Errorf("toml: unknown fields: %v", keys)
		}
	}

	return fmt.Errorf("unsupported file type %q", s.Type)
}

func (s *Store) Load() (*Store, error) {
	data, err := s.File.Read()
	if err != nil {
		return nil, err
	}

	if err := s.unmarshal(data); err != nil {
		return nil, err
	}

	for i, p := range s.Profiles {
		if err := p.ToEnv(); err != nil {
			return nil, fmt.Errorf("profile %q: %w", i, err) // key first, then wrapped err
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
		return err
	}

	return s.File.Write(out)
}

func (s *Store) Exists(name string) bool {
	_, ok := s.Profiles[name]
	return ok
}

// Create adds a new profile to the store.
func (s *Store) create(name string) {
	if _, exists := s.Profiles[name]; !exists {
		s.Profiles[name] = &Profile{Env: map[string]string{}}
	}
}

// Delete removes a profile from the store.
func (s *Store) Delete(name string) error {
	if !s.Exists(name) {
		return fmt.Errorf("unknown profile %q", name)
	}

	delete(s.Profiles, name)

	return nil
}

// SetVar sets or updates a variable in a profile.
func (s *Store) SetVar(name, k, v string) error {
	s.create(name)
	return s.Profiles[name].Env.AddPair(k, v)
}

// RemoveVar deletes a variable from a profile.
func (s *Store) RemoveVar(name, k string) error {
	if !s.Exists(name) {
		return fmt.Errorf("unknown profile %q", name)
	}

	env := s.Profiles[name].Env

	if !env.Exists(k) {
		return fmt.Errorf("variable %q not found in profile %q", k, name)
	}

	env.Delete(k)

	return nil
}

// Vars returns the merged environment variables for a profile, resolving dependencies.
func (s *Store) Vars(name string) (map[string]string, error) {
	// Build dependency DAG.
	nodes := make([]string, 0, len(s.Profiles))
	for n := range s.Profiles {
		nodes = append(nodes, n)
	}
	g, err := dag.Build(nodes, func(n string) []string { return s.Profiles[n].Extends })
	if err != nil {
		return nil, err
	}
	chain, err := g.Chain(name)
	if err != nil {
		return nil, err
	}
	merged := map[string]string{}
	for _, n := range chain {
		for k, v := range s.Profiles[n].Env {
			merged[k] = v
		}
	}
	return merged, nil
}

// ProfilesSorted returns the profile names in sorted order.
func (s *Store) ProfilesSorted() []string {
	out := make([]string, 0, len(s.Profiles))
	for k := range s.Profiles {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}
