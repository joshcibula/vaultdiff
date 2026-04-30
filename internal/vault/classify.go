package vault

import "strings"

// ClassifyOptions controls how secrets are classified into named tiers or categories.
type ClassifyOptions struct {
	Enabled    bool
	Rules      []ClassifyRule
	DefaultTag string
}

// ClassifyRule maps a path or key pattern to a classification tag.
type ClassifyRule struct {
	PathPrefix string
	KeyPrefix  string
	Tag        string
}

// DefaultClassifyOptions returns a disabled ClassifyOptions.
func DefaultClassifyOptions() ClassifyOptions {
	return ClassifyOptions{
		Enabled:    false,
		DefaultTag: "unclassified",
	}
}

// ClassifySecrets annotates each secret's keys with a "_class" metadata key
// derived from the matching rule, or the default tag if no rule matches.
// The input map is not mutated; a new map is returned.
func ClassifySecrets(secrets map[string]map[string]string, opts ClassifyOptions) map[string]map[string]string {
	if !opts.Enabled || len(secrets) == 0 {
		return secrets
	}

	out := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		copy := make(map[string]string, len(kv)+1)
		for k, v := range kv {
			copy[k] = v
		}
		copy["_class"] = resolveClass(path, kv, opts)
		out[path] = copy
	}
	return out
}

func resolveClass(path string, kv map[string]string, opts ClassifyOptions) string {
	for _, rule := range opts.Rules {
		if rule.PathPrefix != "" && strings.HasPrefix(path, rule.PathPrefix) {
			return rule.Tag
		}
		if rule.KeyPrefix != "" {
			for k := range kv {
				if strings.HasPrefix(k, rule.KeyPrefix) {
					return rule.Tag
				}
			}
		}
	}
	return opts.DefaultTag
}
