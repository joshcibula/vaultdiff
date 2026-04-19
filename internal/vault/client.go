package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	v *vaultapi.Client
}

// NewClient creates a new Vault client using environment variables or provided address.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()

	if address != "" {
		cfg.Address = address
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	} else {
		cfg.Address = "http://127.0.0.1:8200"
	}

	if err := cfg.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("reading vault environment: %w", err)
	}

	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	} else if t := os.Getenv("VAULT_TOKEN"); t != "" {
		client.SetToken(t)
	}

	return &Client{v: client}, nil
}

// ReadSecrets reads KV v2 secrets at the given path and returns a flat map of key->value.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.v.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading path %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 fallback
		data = secret.Data
	}

	raw, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at path %q", path)
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
