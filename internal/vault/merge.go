package vault

// MergeOptions controls how two secret maps are merged.
type MergeOptions struct {
	// PreferLeft means left-side values win on conflict; otherwise right wins.
	PreferLeft bool
	// SkipEmpty skips keys whose value is an empty string.
	SkipEmpty bool
}

// DefaultMergeOptions returns sensible defaults: right-side wins, empty values included.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		PreferLeft: false,
		SkipEmpty:  false,
	}
}

// MergeSecrets merges two maps of secret paths → key/value pairs into one.
// Conflicts are resolved according to opts.
func MergeSecrets(
	left, right map[string]map[string]string,
	opts MergeOptions,
) map[string]map[string]string {
	result := make(map[string]map[string]string)

	// Copy left side first.
	for path, kv := range left {
		result[path] = copyKV(kv, opts.SkipEmpty)
	}

	// Merge right side.
	for path, kv := range right {
		if _, exists := result[path]; !exists {
			result[path] = copyKV(kv, opts.SkipEmpty)
			continue
		}
		for k, v := range kv {
			if opts.SkipEmpty && v == "" {
				continue
			}
			if _, conflict := result[path][k]; conflict && opts.PreferLeft {
				continue
			}
			result[path][k] = v
		}
	}

	return result
}

func copyKV(src map[string]string, skipEmpty bool) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		if skipEmpty && v == "" {
			continue
		}
		dst[k] = v
	}
	return dst
}
