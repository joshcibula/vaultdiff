package vault

import "strings"

// TruncateOptions controls how secret values are truncated before display.
type TruncateOptions struct {
	// Enabled turns truncation on. When false, values are returned as-is.
	Enabled bool

	// MaxLength is the maximum number of runes to retain per value.
	// Values longer than this are cut and suffixed with Ellipsis.
	MaxLength int

	// Ellipsis is appended to truncated values. Defaults to "...".
	Ellipsis string

	// SkipKeys lists key names whose values should never be truncated.
	SkipKeys []string
}

// DefaultTruncateOptions returns a safe default configuration.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		Enabled:   false,
		MaxLength: 64,
		Ellipsis:  "...",
	}
}

// TruncateSecrets returns a copy of secrets with values truncated according to opts.
// The original map is never mutated.
func TruncateSecrets(secrets map[string]map[string]string, opts TruncateOptions) map[string]map[string]string {
	if !opts.Enabled {
		return secrets
	}

	if opts.Ellipsis == "" {
		opts.Ellipsis = "..."
	}

	result := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		copy := make(map[string]string, len(kv))
		for k, v := range kv {
			if containsString(opts.SkipKeys, k) {
				copy[k] = v
			} else {
				copy[k] = truncateValue(v, opts.MaxLength, opts.Ellipsis)
			}
		}
		result[path] = copy
	}
	return result
}

// truncateValue cuts s to maxLen runes and appends ellipsis if it was longer.
func truncateValue(s string, maxLen int, ellipsis string) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return strings.TrimRight(string(runes[:maxLen]), " ") + ellipsis
}
