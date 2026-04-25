package vault

import "strings"

// AnnotateOptions controls how secret paths and keys are annotated.
type AnnotateOptions struct {
	Enabled    bool
	TagKey     string            // key name injected into each secret map
	TagValue   string            // static value to inject (e.g. environment name)
	PathPrefix string            // prefix to strip before storing as tag value
	CustomTags map[string]string // additional key→value pairs injected into every secret
}

// DefaultAnnotateOptions returns safe no-op defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		Enabled:    false,
		TagKey:     "_vaultdiff_source",
		CustomTags: map[string]string{},
	}
}

// AnnotateSecrets injects metadata tags into each secret map.
// The original input is not mutated; a new map is returned.
func AnnotateSecrets(secrets map[string]map[string]string, opts AnnotateOptions) map[string]map[string]string {
	if !opts.Enabled {
		return secrets
	}

	out := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		copy := make(map[string]string, len(kv)+1+len(opts.CustomTags))
		for k, v := range kv {
			copy[k] = v
		}

		// Inject primary tag
		if opts.TagKey != "" {
			tagVal := opts.TagValue
			if tagVal == "" {
				tagVal = strings.TrimPrefix(path, opts.PathPrefix)
			}
			copy[opts.TagKey] = tagVal
		}

		// Inject custom tags
		for k, v := range opts.CustomTags {
			copy[k] = v
		}

		out[path] = copy
	}
	return out
}
