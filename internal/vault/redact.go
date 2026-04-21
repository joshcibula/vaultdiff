package vault

import "strings"

// RedactOptions controls which secret paths or keys are fully redacted from output.
type RedactOptions struct {
	// Paths is a list of path prefixes whose secrets will be fully redacted.
	Paths []string
	// Keys is a list of secret key names that will be fully redacted.
	Keys []string
}

// DefaultRedactOptions returns a RedactOptions with no redactions applied.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{}
}

const redactedPlaceholder = "[REDACTED]"

// RedactSecrets returns a copy of secrets with any matching paths or keys
// replaced by the redacted placeholder.
func RedactSecrets(secrets map[string]map[string]string, opts RedactOptions) map[string]map[string]string {
	result := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		if matchesAnyPrefix(path, opts.Paths) {
			redacted := make(map[string]string, len(kv))
			for k := range kv {
				redacted[k] = redactedPlaceholder
			}
			result[path] = redacted
			continue
		}
		entry := make(map[string]string, len(kv))
		for k, v := range kv {
			if containsString(opts.Keys, k) {
				entry[k] = redactedPlaceholder
			} else {
				entry[k] = v
			}
		}
		result[path] = entry
	}
	return result
}

func matchesAnyPrefix(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
