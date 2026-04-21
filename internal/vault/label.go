package vault

import "strings"

// LabelOptions controls how secret paths are labelled in output.
type LabelOptions struct {
	// Prefix is prepended to every path label (e.g. "prod/").
	Prefix string
	// StripPrefix removes a leading string from every path label before display.
	StripPrefix string
	// Alias maps an exact path to a human-friendly label.
	Alias map[string]string
}

// DefaultLabelOptions returns a no-op LabelOptions.
func DefaultLabelOptions() LabelOptions {
	return LabelOptions{
		Alias: make(map[string]string),
	}
}

// LabelSecrets returns a new map with keys (paths) relabelled according to opts.
// Values are copied unchanged.
func LabelSecrets(secrets map[string]map[string]string, opts LabelOptions) map[string]map[string]string {
	if opts.Prefix == "" && opts.StripPrefix == "" && len(opts.Alias) == 0 {
		return secrets
	}

	out := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		label := applyLabel(path, opts)
		copy := make(map[string]string, len(kv))
		for k, v := range kv {
			copy[k] = v
		}
		out[label] = copy
	}
	return out
}

func applyLabel(path string, opts LabelOptions) string {
	if alias, ok := opts.Alias[path]; ok {
		return alias
	}
	label := path
	if opts.StripPrefix != "" {
		label = strings.TrimPrefix(label, opts.StripPrefix)
	}
	if opts.Prefix != "" {
		label = opts.Prefix + label
	}
	return label
}
