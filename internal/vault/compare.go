package vault

import "strings"

// CompareOptions controls how two secret maps are compared.
type CompareOptions struct {
	IgnoreCase     bool
	IgnoreWhitespace bool
	IgnoreKeys     []string
}

// DefaultCompareOptions returns sensible defaults.
func DefaultCompareOptions() CompareOptions {
	return CompareOptions{
		IgnoreCase:       false,
		IgnoreWhitespace: false,
		IgnoreKeys:       nil,
	}
}

// CompareResult holds the outcome of comparing two secret values.
type CompareResult struct {
	Path     string
	Key      string
	LeftVal  string
	RightVal string
	Equal    bool
}

// CompareSecrets performs a field-level comparison of two secret maps,
// returning one CompareResult per (path, key) pair found in either side.
func CompareSecrets(left, right map[string]map[string]string, opts CompareOptions) []CompareResult {
	ignored := make(map[string]struct{}, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignored[k] = struct{}{}
	}

	paths := unionKeys(left, right)
	var results []CompareResult

	for _, path := range paths {
		lKV := left[path]
		rKV := right[path]
		keys := unionKeys(lKV, rKV)
		for _, key := range keys {
			if _, skip := ignored[key]; skip {
				continue
			}
			lv := normalize(lKV[key], opts)
			rv := normalize(rKV[key], opts)
			results = append(results, CompareResult{
				Path:     path,
				Key:      key,
				LeftVal:  lKV[key],
				RightVal: rKV[key],
				Equal:    lv == rv,
			})
		}
	}
	return results
}

func normalize(v string, opts CompareOptions) string {
	if opts.IgnoreWhitespace {
		v = strings.TrimSpace(v)
	}
	if opts.IgnoreCase {
		v = strings.ToLower(v)
	}
	return v
}

func unionKeys[V any](a, b map[string]V) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
