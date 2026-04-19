package vault

import "strings"

// FilterOptions controls which secrets are included in a diff.
type FilterOptions struct {
	// Prefix restricts secrets to those whose path starts with the given prefix.
	Prefix string
	// ExcludeKeys is a list of secret keys to omit from comparison.
	ExcludeKeys []string
}

// FilterSecrets returns a new map containing only the secrets that pass the
// filter criteria defined in opts.
func FilterSecrets(secrets map[string]map[string]interface{}, opts FilterOptions) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	for path, kv := range secrets {
		if opts.Prefix != "" && !strings.HasPrefix(path, opts.Prefix) {
			continue
		}

		filtered := make(map[string]interface{})
		for k, v := range kv {
			if !containsString(opts.ExcludeKeys, k) {
				filtered[k] = v
			}
		}

		if len(filtered) > 0 {
			result[path] = filtered
		}
	}

	return result
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
