package vault

import (
	"context"
	"fmt"
	"strings"
)

// ListSecrets returns all secret keys under a given path recursively.
func (c *Client) ListSecrets(ctx context.Context, path string) ([]string, error) {
	return c.listRecursive(ctx, path, "")
}

func (c *Client) listRecursive(ctx context.Context, basePath, prefix string) ([]string, error) {
	fullPath := basePath
	if prefix != "" {
		fullPath = strings.TrimRight(basePath, "/") + "/" + prefix
	}

	secret, err := c.logical.ListWithContext(ctx, "secret/metadata/"+fullPath)
	if err != nil {
		return nil, fmt.Errorf("listing %s: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, nil
	}

	var results []string
	for _, k := range keys {
		key, _ := k.(string)
		if strings.HasSuffix(key, "/") {
			// directory — recurse
			sub, err := c.listRecursive(ctx, basePath, strings.TrimRight(prefix+key, "/"))
			if err != nil {
				return nil, err
			}
			results = append(results, sub...)
		} else {
			if prefix != "" {
				results = append(results, prefix+"/"+key)
			} else {
				results = append(results, key)
			}
		}
	}
	return results, nil
}
