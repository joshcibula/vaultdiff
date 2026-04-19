package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestAuthConfigFromEnv_Token(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test-token")
	os.Unsetenv("VAULT_ROLE_ID")
	os.Unsetenv("VAULT_SECRET_ID")
	t.Cleanup(func() { os.Unsetenv("VAULT_TOKEN") })

	cfg := AuthConfigFromEnv()
	if cfg.Method != AuthToken {
		t.Errorf("expected token method, got %s", cfg.Method)
	}
	if cfg.Token != "test-token" {
		t.Errorf("expected token 'test-token', got %s", cfg.Token)
	}
}

func TestAuthConfigFromEnv_AppRole(t *testing.T) {
	os.Setenv("VAULT_ROLE_ID", "my-role")
	os.Setenv("VAULT_SECRET_ID", "my-secret")
	t.Cleanup(func() {
		os.Unsetenv("VAULT_ROLE_ID")
		os.Unsetenv("VAULT_SECRET_ID")
	})

	cfg := AuthConfigFromEnv()
	if cfg.Method != AuthAppRole {
		t.Errorf("expected approle method, got %s", cfg.Method)
	}
}

func TestAuthenticate_Token(t *testing.T) {
	client, _ := api.NewClient(api.DefaultConfig())
	cfg := AuthConfig{Method: AuthToken, Token: "s.abc123"}
	if err := Authenticate(client, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Token() != "s.abc123" {
		t.Errorf("expected token s.abc123, got %s", client.Token())
	}
}

func TestAuthenticate_MissingToken(t *testing.T) {
	client, _ := api.NewClient(api.DefaultConfig())
	cfg := AuthConfig{Method: AuthToken}
	if err := Authenticate(client, cfg); err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestAuthenticate_AppRole(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/auth/approle/login" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"auth": map[string]interface{}{"client_token": "approle-token"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cfg := api.DefaultConfig()
	cfg.Address = ts.URL
	client, _ := api.NewClient(cfg)

	auth := AuthConfig{Method: AuthAppRole, RoleID: "role", SecretID: "secret"}
	if err := Authenticate(client, auth); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Token() != "approle-token" {
		t.Errorf("expected approle-token, got %s", client.Token())
	}
}
