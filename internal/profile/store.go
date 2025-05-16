// Package profile provides profile and store management for environment variable sets.
package profile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"

	"github.com/idelchi/envsync/internal/dag"
	"github.com/idelchi/envsync/internal/unmarshal"
)

// priority defines the order of file formats to check for.
var priority = []string{"envsync.yaml", "envsync.yml", "envsync.toml", "envsync.json"}

// Profile holds metadata plus an env-var map.
type Profile struct {
	Extends unmarshal.SingleOrSliceType[string] `json:"extends,omitempty" yaml:"extends,omitempty" toml:"extends,omitempty"`
	Env     map[string]string                   `json:"env,omitempty"     yaml:"env,omitempty"     toml:"env,omitempty"`
}

// Store represents the entire file and codec information.
type Store struct {
	Path     string
	Ext      string
	Profiles map[string]*Profile
}

func findPath(flag string) (string, string, error) {
	if flag != "" {
		return flag, filepath.Ext(flag), nil
	}
	for _, f := range priority {
		if _, err := os.Stat(f); err == nil {
			return f, filepath.Ext(f), nil
		}
	}
	return priority[0], filepath.Ext(priority[0]), nil
}

// Load reads a profile store from the given file or default location.
func Load(flag string) (*Store, error) {
	path, ext, err := findPath(flag)
	if err != nil {
		return nil, err
	}
	s := &Store{Path: path, Ext: ext, Profiles: map[string]*Profile{}}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return s, nil // empty/new store
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	switch ext {
	case ".yaml", ".yml":
		err = yaml.UnmarshalWithOptions(raw, &s.Profiles, yaml.Strict())
	case ".json":
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		err = dec.Decode(&s.Profiles)
	case ".toml":
		md, err := toml.Decode(string(raw), &s.Profiles)
		if err != nil {
			return nil, err
		}
		if keys := md.Undecoded(); len(keys) > 0 {
			return nil, fmt.Errorf("toml: unknown fields: %v", keys)
		}
	default:
		err = fmt.Errorf("unsupported format %q", strings.TrimPrefix(ext, "."))
	}
	return s, err
}

// Save writes the store to disk.
func (s *Store) Save() error {
	var (
		out []byte
		err error
	)
	switch s.Ext {
	case ".yaml", ".yml":
		out, err = yaml.Marshal(s.Profiles)
	case ".json":
		out, err = json.MarshalIndent(s.Profiles, "", "  ")
	case ".toml":
		out, err = toml.Marshal(s.Profiles)
	default:
		err = fmt.Errorf("unsupported format %q", strings.TrimPrefix(s.Ext, "."))
	}
	if err != nil {
		return err
	}
	tmp := s.Path + ".tmp"
	if err = os.WriteFile(tmp, out, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.Path)
}

func (s *Store) ensure(name string) *Profile {
	p, ok := s.Profiles[name]
	if !ok {
		p = &Profile{Env: map[string]string{}}
		s.Profiles[name] = p
	}
	if p.Env == nil {
		p.Env = map[string]string{}
	}
	return p
}

// Create adds a new profile to the store.
func (s *Store) Create(name string) error {
	if _, exists := s.Profiles[name]; exists {
		return fmt.Errorf("profile %q exists", name)
	}
	s.Profiles[name] = &Profile{Env: map[string]string{}}
	return s.Save()
}

// Delete removes a profile from the store.
func (s *Store) Delete(name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return fmt.Errorf("unknown profile %s", name)
	}
	delete(s.Profiles, name)
	return s.Save()
}

// AddVar adds a new variable to a profile.
func (s *Store) AddVar(name, k, v string) error {
	p := s.ensure(name)
	if _, dup := p.Env[k]; dup {
		return fmt.Errorf("%s already set in %s", k, name)
	}
	p.Env[k] = v
	return s.Save()
}

// SetVar sets or updates a variable in a profile.
func (s *Store) SetVar(name, k, v string) error {
	s.ensure(name).Env[k] = v
	return s.Save()
}

// RemoveVar deletes a variable from a profile.
func (s *Store) RemoveVar(name, k string) error {
	p, ok := s.Profiles[name]
	if !ok {
		return fmt.Errorf("unknown profile %s", name)
	}
	if _, ok := p.Env[k]; !ok {
		return fmt.Errorf("%s not found in %s", k, name)
	}
	delete(p.Env, k)
	return s.Save()
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
