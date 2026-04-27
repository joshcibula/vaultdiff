package vault

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested secret maps are flattened.
type FlattenOptions struct {
	Enabled   bool
	Separator string
	MaxDepth  int
}

// DefaultFlattenOptions returns sensible defaults for flattening.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Enabled:   false,
		Separator: ".",
		MaxDepth:  10,
	}
}

// FlattenSecrets takes a map of path -> (key -> value) and flattens any
// values that are themselves nested maps into dot-separated keys.
// Non-map values are passed through unchanged.
func FlattenSecrets(secrets map[string]map[string]string, opts FlattenOptions) map[string]map[string]string {
	if !opts.Enabled {
		return secrets
	}
	sep := opts.Separator
	if sep == "" {
		sep = "."
	}
	maxDepth := opts.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 10
	}

	result := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		flat := make(map[string]string)
		for k, v := range kv {
			flat[k] = v
		}
		result[path] = flat
	}
	return result
}

// flattenMap recursively flattens a nested map[string]interface{} into
// dot-separated string keys, writing results into out.
func flattenMap(prefix string, m map[string]interface{}, sep string, depth, maxDepth int, out map[string]string) {
	if depth > maxDepth {
		out[prefix] = fmt.Sprintf("%v", m)
		return
	}
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			flattenMap(key, child, sep, depth+1, maxDepth, out)
		default:
			out[key] = strings.TrimSpace(fmt.Sprintf("%v", v))
		}
	}
}
