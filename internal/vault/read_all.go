package vault

import (
	"context"
	"fmt"
)

// ReadAll lists all secrets under basePath and reads each one,
// returning a map of relative key path -> secret data map.
func (c *Client) ReadAll(ctx context.Context, basePath string) (map[string]map[string]string, error) {
	keys, err := c.ListSecrets(ctx, basePath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	result := make(map[string]map[string]string, len(keys))
	for _, key := range keys {
		fullPath := basePath + "/" + key
		data, err := c.ReadSecrets(ctx, fullPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", fullPath, err)
		}
		result[key] = data
	}
	return result, nil
}
