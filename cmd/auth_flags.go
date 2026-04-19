package cmd

import (
	"github.com/spf13/cobra"
	vaultclient "vaultdiff/internal/vault"
)

// authFlags holds parsed authentication flag values.
type authFlags struct {
	token    string
	roleID   string
	secretID string
}

// registerAuthFlags adds auth-related flags to the given command.
func registerAuthFlags(cmd *cobra.Command, f *authFlags) {
	cmd.Flags().StringVar(&f.token, "token", "", "Vault token (overrides VAULT_TOKEN)")
	cmd.Flags().StringVar(&f.roleID, "role-id", "", "AppRole role ID (overrides VAULT_ROLE_ID)")
	cmd.Flags().StringVar(&f.secretID, "secret-id", "", "AppRole secret ID (overrides VAULT_SECRET_ID)")
}

// resolveAuthConfig builds an AuthConfig from flags, falling back to env vars.
func resolveAuthConfig(f authFlags) vaultclient.AuthConfig {
	cfg := vaultclient.AuthConfigFromEnv()
	if f.token != "" {
		cfg.Method = vaultclient.AuthToken
		cfg.Token = f.token
	}
	if f.roleID != "" {
		cfg.RoleID = f.roleID
		cfg.Method = vaultclient.AuthAppRole
	}
	if f.secretID != "" {
		cfg.SecretID = f.secretID
		cfg.Method = vaultclient.AuthAppRole
	}
	return cfg
}
