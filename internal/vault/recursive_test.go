package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSecrets_Flat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"keys": []interface{}{"db", "api"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	keys, err := client.ListSecrets(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("ListSecrets: %v", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
}

func TestListSecrets_Empty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	keys, err := client.ListSecrets(context.Background(), "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(keys))
	}
}
