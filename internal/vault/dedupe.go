package vault

import "strings"

// DedupeOptions controls how duplicate secret keys are resolved across paths.
type DedupeOptions struct {
	Enabled         bool
	CaseSensitive   bool
	PreferLongerPath bool
}

// DefaultDedupeOptions returns sensible defaults for deduplication.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		Enabled:         false,
		CaseSensitive:   true,
		PreferLongerPath: false,
	}
}

// DedupeSecrets removes duplicate keys across paths, keeping the first
// occurrence unless PreferLongerPath is set, in which case the entry
// from the longer path wins.
func DedupeSecrets(secrets map[string]map[string]string, opts DedupeOptions) map[string]map[string]string {
	if !opts.Enabled {
		return secrets
	}

	// Track which path "owns" each key.
	type owner struct {
		path  string
		value string
	}
	seen := make(map[string]owner)

	for path, kv := range secrets {
		for k, v := range kv {
			normKey := k
			if !opts.CaseSensitive {
				normKey = strings.ToLower(k)
			}
			if existing, ok := seen[normKey]; !ok {
				seen[normKey] = owner{path: path, value: v}
			} else if opts.PreferLongerPath && len(path) > len(existing.path) {
				seen[normKey] = owner{path: path, value: v}
			}
		}
	}

	// Rebuild secrets, dropping duplicate keys from non-owning paths.
	result := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		filtered := make(map[string]string)
		for k, v := range kv {
			normKey := k
			if !opts.CaseSensitive {
				normKey = strings.ToLower(k)
			}
			if o, ok := seen[normKey]; ok && o.path == path {
				filtered[k] = v
			}
		}
		if len(filtered) > 0 {
			result[path] = filtered
		}
	}
	return result
}
