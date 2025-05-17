package profile

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Stringify turns any Go value into the dotenv-ready string described earlier.
//   - Scalars → plain strings (`true`, `5432`, `foo`…).
//   - Slices / maps / structs → compact JSON wrapped in single quotes
//     (`'["a","b"]'`, `'{"k":"v"}'`).
func Stringify(v any) (string, error) {
	switch x := v.(type) {

	case nil:
		return "", nil

	case string:
		if needsQuotes(x) {
			return strconv.Quote(x), nil // adds double-quotes and escapes
		}
		return x, nil

	case bool:
		return strconv.FormatBool(x), nil

	// All signed ints
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", x), nil

	// Unsigned ints
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", x), nil

	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64), nil

	default:
		// Anything not matched above gets JSON-encoded.
		raw, err := json.Marshal(x)
		if err != nil {
			return "", err
		}
		return "'" + string(raw) + "'", nil
	}
}

// needsQuotes returns true if the string contains characters that
// are troublesome in a dotenv line (space, =, #, quotes, etc.).
func needsQuotes(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case unicode.IsSpace(r):
			return true
		case strings.ContainsRune(`"'#=`, r):
			return true
		}
	}
	return false
}
