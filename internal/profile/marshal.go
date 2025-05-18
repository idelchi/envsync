package profile

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
)

// marshal encodes the store's profiles into the specified format.
func (s *Store) marshal() ([]byte, error) {
	switch s.Type {
	case YAML:
		data, err := yaml.MarshalWithOptions(s.Profiles)
		if err != nil {
			return nil, fmt.Errorf("yaml: encode: %w", err)
		}

		return data, nil
	case TOML:
		data, err := toml.Marshal(s.Profiles)
		if err != nil {
			return nil, fmt.Errorf("toml: encode: %w", err)
		}

		return data, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedFileType, s.Type)
	}
}

// unmarshal decodes the data into the store's profiles.
func (s *Store) unmarshal(data []byte) error {
	switch s.Type {
	case YAML:
		err := yaml.UnmarshalWithOptions(data, &s.Profiles, yaml.Strict())
		if err != nil {
			return fmt.Errorf("yaml: decode: %w", err)
		}

	case TOML:
		dec := toml.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()

		if err := dec.Decode(&s.Profiles); err != nil {
			var sm *toml.StrictMissingError
			if errors.As(err, &sm) {
				unknown := make([]string, len(sm.Errors))
				for i, de := range sm.Errors {
					unknown[i] = strings.Join(de.Key(), ".")
				}

				return fmt.Errorf( //nolint:err113		// Occasional dynamic errors are fine.
					"toml: unknown fields: %v",
					unknown,
				)
			}

			// not a strict-missing problem – propagate original error
			return fmt.Errorf("toml: decode: %w", err)
		}
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedFileType, s.Type)
	}

	return nil
}
