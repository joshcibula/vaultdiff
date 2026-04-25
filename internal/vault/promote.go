package vault

import "fmt"

// PromoteOptions controls how secrets are promoted from one path to another.
type PromoteOptions struct {
	Enabled    bool
	DryRun     bool
	Overwrite  bool
	PathPrefix string
}

// DefaultPromoteOptions returns safe defaults for promotion.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		Enabled:   false,
		DryRun:    true,
		Overwrite: false,
	}
}

// PromoteResult records the outcome of a single promotion operation.
type PromoteResult struct {
	Path    string
	Key     string
	Skipped bool
	Reason  string
}

// PromoteSecrets copies secrets from src into dst, respecting the given options.
// Existing keys in dst are only overwritten when opts.Overwrite is true.
// When opts.DryRun is true no mutations are applied.
func PromoteSecrets(
	src map[string]map[string]string,
	dst map[string]map[string]string,
	opts PromoteOptions,
) (map[string]map[string]string, []PromoteResult, error) {
	if !opts.Enabled {
		return dst, nil, nil
	}

	output := copyNestedMap(dst)
	var results []PromoteResult

	for path, kv := range src {
		targetPath := path
		if opts.PathPrefix != "" {
			targetPath = opts.PathPrefix + "/" + path
		}

		if _, exists := output[targetPath]; !exists {
			output[targetPath] = make(map[string]string)
		}

		for k, v := range kv {
			if _, conflict := output[targetPath][k]; conflict && !opts.Overwrite {
				results = append(results, PromoteResult{
					Path:    targetPath,
					Key:     k,
					Skipped: true,
					Reason:  fmt.Sprintf("key already exists in destination (overwrite=false)"),
				})
				continue
			}

			results = append(results, PromoteResult{Path: targetPath, Key: k})
			if !opts.DryRun {
				output[targetPath][k] = v
			}
		}
	}

	return output, results, nil
}

func copyNestedMap(m map[string]map[string]string) map[string]map[string]string {
	out := make(map[string]map[string]string, len(m))
	for path, kv := range m {
		copy := make(map[string]string, len(kv))
		for k, v := range kv {
			copy[k] = v
		}
		out[path] = copy
	}
	return out
}
