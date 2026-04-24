package vault

import "fmt"

// PatchOperation represents a single patch action on a secret map.
type PatchOperation struct {
	Op    string // "set", "delete", "rename"
	Path  string
	Key   string
	Value string
	NewKey string
}

// DefaultPatchOptions returns a no-op PatchOptions.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{}
}

// PatchOptions controls how patches are applied.
type PatchOptions struct {
	Operations []PatchOperation
	DryRun     bool
}

// PatchSecrets applies a list of patch operations to the given secret map.
// It returns a new map with the operations applied and a list of applied op descriptions.
func PatchSecrets(secrets map[string]map[string]string, opts PatchOptions) (map[string]map[string]string, []string, error) {
	result := make(map[string]map[string]string, len(secrets))
	for path, kv := range secrets {
		copy := make(map[string]string, len(kv))
		for k, v := range kv {
			copy[k] = v
		}
		result[path] = copy
	}

	var applied []string

	for _, op := range opts.Operations {
		kv, ok := result[op.Path]
		if !ok {
			return nil, nil, fmt.Errorf("patch: path %q not found", op.Path)
		}
		switch op.Op {
		case "set":
			if !opts.DryRun {
				kv[op.Key] = op.Value
			}
			applied = append(applied, fmt.Sprintf("set %s/%s", op.Path, op.Key))
		case "delete":
			if !opts.DryRun {
				delete(kv, op.Key)
			}
			applied = append(applied, fmt.Sprintf("delete %s/%s", op.Path, op.Key))
		case "rename":
			if !opts.DryRun {
				if v, exists := kv[op.Key]; exists {
					kv[op.NewKey] = v
					delete(kv, op.Key)
				}
			}
			applied = append(applied, fmt.Sprintf("rename %s/%s -> %s", op.Path, op.Key, op.NewKey))
		default:
			return nil, nil, fmt.Errorf("patch: unknown operation %q", op.Op)
		}
	}

	return result, applied, nil
}
