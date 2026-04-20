package vault

import "strings"

// MaskOptions controls how secret values are masked in output.
type MaskOptions struct {
	Enabled    bool
	MaskString string
	RevealKeys []string
}

// DefaultMaskOptions returns sensible defaults.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Enabled:    false,
		MaskString: "***",
		RevealKeys: []string{},
	}
}

// MaskSecrets returns a copy of secrets with values replaced by the mask string,
// unless the key is in RevealKeys.
func MaskSecrets(secrets map[string]map[string]string, opts MaskOptions) map[string]map[string]string {
	if !opts.Enabled {
		return secrets
	}
	masked := make(map[string]map[string]string, len(secrets))
	for path, kvs := range secrets {
		maskedKVs := make(map[string]string, len(kvs))
		for k, v := range kvs {
			if containsString(opts.RevealKeys, k) {
				maskedKVs[k] = v
			} else {
				maskedKVs[k] = opts.MaskString
			}
		}
		masked[path] = maskedKVs
	}
	return masked
}

// MaskValue masks a single value unless its key is revealed.
func MaskValue(key, value string, opts MaskOptions) string {
	if !opts.Enabled {
		return value
	}
	for _, r := range opts.RevealKeys {
		if strings.EqualFold(r, key) {
			return value
		}
	}
	return opts.MaskString
}
