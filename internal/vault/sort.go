package vault

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering direction for secret sorting.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortField defines which field to sort secrets by.
type SortField string

const (
	SortByPath SortField = "path"
	SortByKey  SortField = "key"
	SortByValue SortField = "value"
)

// DefaultSortOptions returns conservative defaults (path, ascending).
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Enabled: false,
		Field:   SortByPath,
		Order:   SortAsc,
	}
}

// SortOptions controls how secrets are sorted before diffing or rendering.
type SortOptions struct {
	Enabled bool
	Field   SortField
	Order   SortOrder
}

// SortSecrets returns a new map with keys sorted according to opts.
// Each path's key-value pairs may also be sorted when Field is SortByKey or SortByValue.
func SortSecrets(secrets map[string]map[string]string, opts SortOptions) map[string]map[string]string {
	if !opts.Enabled || len(secrets) == 0 {
		return secrets
	}

	// Collect and sort paths.
	paths := make([]string, 0, len(secrets))
	for p := range secrets {
		paths = append(paths, p)
	}

	sortStrings(paths, opts.Order)

	result := make(map[string]map[string]string, len(secrets))
	for _, p := range paths {
		orig := secrets[p]
		if opts.Field == SortByKey || opts.Field == SortByValue {
			result[p] = sortKV(orig, opts.Field, opts.Order)
		} else {
			// Copy as-is; path ordering is captured in the path slice above.
			copy := make(map[string]string, len(orig))
			for k, v := range orig {
				copy[k] = v
			}
			result[p] = copy
		}
	}
	return result
}

func sortStrings(ss []string, order SortOrder) {
	sort.Slice(ss, func(i, j int) bool {
		cmp := strings.Compare(ss[i], ss[j])
		if order == SortDesc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func sortKV(kv map[string]string, field SortField, order SortOrder) map[string]string {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	if field == SortByValue {
		sort.Slice(keys, func(i, j int) bool {
			cmp := strings.Compare(kv[keys[i]], kv[keys[j]])
			if order == SortDesc {
				return cmp > 0
			}
			return cmp < 0
		})
	} else {
		sortStrings(keys, order)
	}
	out := make(map[string]string, len(kv))
	for _, k := range keys {
		out[k] = kv[k]
	}
	return out
}
