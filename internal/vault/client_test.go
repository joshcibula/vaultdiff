package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockVaultServer(t *testing.T, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
}

func TestNewClient_DefaultAddress(t *testing.T) {
	client, err := NewClient("", "test-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client == nil {
		t.Fatal("expected client, got nil")
	}
}

func TestReadSecrets_KVv2(t *testing.T) {
	response := `{"data":{"data":{"foo":"bar","baz":"qux"}}}`
	server := mockVaultServer(t, response)
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("reading secrets: %v", err)
	}

	if secrets["foo"] != "bar" {
		t.Errorf("expected foo=bar, got foo=%s", secrets["foo"])
	}
	if secrets["baz"] != "qux" {
		t.Errorf("expected baz=qux, got baz=%s", secrets["baz"])
	}
}

func TestReadSecrets_NoSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`null`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	_, err = client.ReadSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
