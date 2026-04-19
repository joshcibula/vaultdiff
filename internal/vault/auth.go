package vault

import (
	"errors"
	"os"

	"github.com/hashicorp/vault/api"
)

// AuthMethod represents a supported Vault authentication method.
type AuthMethod string

const (
	AuthToken     AuthMethod = "token"
	AuthAppRole   AuthMethod = "approle"
)

// AuthConfig holds credentials for authenticating to Vault.
type AuthConfig struct {
	Method   AuthMethod
	Token    string
	RoleID   string
	SecretID string
}

// AuthConfigFromEnv builds an AuthConfig from environment variables.
func AuthConfigFromEnv() AuthConfig {
	cfg := AuthConfig{
		Method:   AuthToken,
		Token:    os.Getenv("VAULT_TOKEN"),
		RoleID:   os.Getenv("VAULT_ROLE_ID"),
		SecretID: os.Getenv("VAULT_SECRET_ID"),
	}
	if cfg.RoleID != "" && cfg.SecretID != "" {
		cfg.Method = AuthAppRole
	}
	return cfg
}

// Authenticate applies the given AuthConfig to the Vault client.
func Authenticate(client *api.Client, cfg AuthConfig) error {
	switch cfg.Method {
	case AuthToken:
		if cfg.Token == "" {
			return errors.New("vault token is required (set VAULT_TOKEN)")
		}
		client.SetToken(cfg.Token)
		return nil
	case AuthAppRole:
		return authenticateAppRole(client, cfg)
	default:
		return errors.New("unsupported auth method: " + string(cfg.Method))
	}
}

func authenticateAppRole(client *api.Client, cfg AuthConfig) error {
	data := map[string]interface{}{
		"role_id":   cfg.RoleID,
		"secret_id": cfg.SecretID,
	}
	secret, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		return err
	}
	if secret == nil || secret.Auth == nil {
		return errors.New("approle login returned no auth info")
	}
	client.SetToken(secret.Auth.ClientToken)
	return nil
}
