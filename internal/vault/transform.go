package vault

import "strings"

// TransformOptions controls how secret values are transformed before comparison.
type TransformOptions struct {
	// TrimSpace removes leading and trailing whitespace from all values.
	TrimSpace bool
	// LowercaseKeys normalizes all secret keys to lowercase.
	LowercaseKeys bool
	// IgnoreKeys is a set of key names whose values are zeroed out before diffing.
	IgnoreKeys []string
}

// DefaultTransformOptions returns a TransformOptions with no transformations enabled.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// TransformSecrets applies the given TransformOptions to a map of path -> (key -> value) secrets.
// It returns a new map; the original is not modified.
func TransformSecrets(secrets map[string]map[string]string, opts TransformOptions) map[string]map[string]string {
	result := make(map[string]map[string]string, len(secrets))
	for path, kvs := range secrets {
		transformed := make(map[string]string, len(kvs))
		for k, v := range kvs {
			key := k
			if opts.LowercaseKeys {
				key = strings.ToLower(k)
			}
			val := v
			if opts.TrimSpace {
				val = strings.TrimSpace(v)
			}
			if containsString(opts.IgnoreKeys, k) || containsString(opts.IgnoreKeys, key) {
				val = ""
			}
			transformed[key] = val
		}
		result[path] = transformed
	}
	return result
}
