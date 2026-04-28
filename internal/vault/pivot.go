package vault

// PivotOptions controls how secrets are pivoted (transposed) so that
// secret keys become top-level paths and paths become inner keys.
type PivotOptions struct {
	Enabled    bool
	KeyField   string // which inner key to use as the new path segment
	PathPrefix string // optional prefix to prepend to pivoted paths
}

// DefaultPivotOptions returns conservative defaults (pivot disabled).
func DefaultPivotOptions() PivotOptions {
	return PivotOptions{
		Enabled:  false,
		KeyField: "",
	}
}

// PivotSecrets transposes a map[path]map[key]value into
// map[key]map[path]value so callers can compare how a single key
// varies across many paths.
//
// If opts.Enabled is false the input is returned unchanged.
// If opts.KeyField is non-empty only that key is pivoted; all others
// are dropped from the output.
func PivotSecrets(secrets map[string]map[string]string, opts PivotOptions) map[string]map[string]string {
	if !opts.Enabled {
		return secrets
	}

	result := make(map[string]map[string]string)

	for path, kv := range secrets {
		for key, value := range kv {
			if opts.KeyField != "" && key != opts.KeyField {
				continue
			}

			pivotPath := key
			if opts.PathPrefix != "" {
				pivotPath = opts.PathPrefix + "/" + key
			}

			if result[pivotPath] == nil {
				result[pivotPath] = make(map[string]string)
			}
			result[pivotPath][path] = value
		}
	}

	return result
}
