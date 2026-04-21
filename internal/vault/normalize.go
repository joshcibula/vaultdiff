package vault

import (
	"strings"
)

// NormalizeOptions controls how secret keys and values are normalized
// before comparison.
type NormalizeOptions struct {
	// TrimKeyPrefix removes a leading prefix from all secret keys.
	TrimKeyPrefix string

	// StripTrailingSlash removes trailing slashes from path keys.
	StripTrailingSlash bool

	// CollapseWhitespace replaces runs of whitespace in values with a single space.
	CollapseWhitespace bool
}

// DefaultNormalizeOptions returns NormalizeOptions with safe defaults (no-op).
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{}
}

// NormalizeSecrets applies normalization rules to a map of secret paths to
// key/value pairs. It returns a new map and does not mutate the input.
func NormalizeSecrets(secrets map[string]map[string]string, opts NormalizeOptions) map[string]map[string]string {
	result := make(map[string]map[string]string, len(secrets))

	for path, kv := range secrets {
		normPath := path
		if opts.StripTrailingSlash {
			normPath = strings.TrimRight(normPath, "/")
		}

		normKV := make(map[string]string, len(kv))
		for k, v := range kv {
			normKey := k
			if opts.TrimKeyPrefix != "" {
				normKey = strings.TrimPrefix(normKey, opts.TrimKeyPrefix)
			}

			normVal := v
			if opts.CollapseWhitespace {
				normVal = collapseWhitespace(normVal)
			}

			normKV[normKey] = normVal
		}
		result[normPath] = normKV
	}

	return result
}

// collapseWhitespace replaces consecutive whitespace characters with a single space
// and trims leading/trailing whitespace.
func collapseWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}
