package vault

import "sort"

// GroupOptions controls how secrets are grouped.
type GroupOptions struct {
	Enabled   bool
	GroupBy   string // "mount", "prefix", or "depth"
	Depth     int    // used when GroupBy == "depth"
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Enabled: false,
		GroupBy: "mount",
		Depth:   1,
	}
}

// SecretGroup holds a named collection of secrets.
type SecretGroup struct {
	Name    string
	Secrets map[string]map[string]string
}

// GroupSecrets partitions a flat secrets map into named groups.
func GroupSecrets(secrets map[string]map[string]string, opts GroupOptions) []SecretGroup {
	if !opts.Enabled || len(secrets) == 0 {
		return []SecretGroup{{Name: "", Secrets: secrets}}
	}

	buckets := make(map[string]map[string]map[string]string)

	for path, kv := range secrets {
		key := groupKey(path, opts)
		if buckets[key] == nil {
			buckets[key] = make(map[string]map[string]string)
		}
		buckets[key][path] = kv
	}

	groups := make([]SecretGroup, 0, len(buckets))
	for name, s := range buckets {
		groups = append(groups, SecretGroup{Name: name, Secrets: s})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups
}

// groupKey derives the bucket key for a path given the options.
func groupKey(path string, opts GroupOptions) string {
	switch opts.GroupBy {
	case "prefix":
		return pathSegment(path, 1)
	case "depth":
		return pathSegment(path, opts.Depth)
	default: // "mount"
		return pathSegment(path, 2)
	}
}

// pathSegment returns the first n slash-delimited segments of path.
func pathSegment(path string, n int) string {
	count := 0
	for i, c := range path {
		if c == '/' {
			count++
			if count == n {
				return path[:i]
			}
		}
	}
	return path
}
