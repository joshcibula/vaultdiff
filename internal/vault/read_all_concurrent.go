package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// ReadAllConcurrent lists all secrets under mountPath and fetches them in parallel.
// It returns a merged map of path -> key/value secrets and a slice of non-fatal errors.
func ReadAllConcurrent(ctx context.Context, client *vaultapi.Client, mountPath string, concOpts ConcurrencyOptions) (map[string]map[string]string, []error) {
	paths, err := ListSecrets(ctx, client, mountPath)
	if err != nil {
		return nil, []error{fmt.Errorf("list secrets: %w", err)}
	}

	fetch := func(ctx context.Context, path string) (map[string]string, error) {
		return ReadSecrets(ctx, client, path)
	}

	results := FetchAllConcurrent(ctx, paths, concOpts, fetch)

	merged := make(map[string]map[string]string, len(results))
	var errs []error

	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, fmt.Errorf("read %s: %w", r.Path, r.Err))
			continue
		}
		if len(r.Secrets) > 0 {
			merged[r.Path] = r.Secrets
		}
	}

	return merged, errs
}
