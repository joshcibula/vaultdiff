package vault

import (
	"testing"
	"time"
)

func TestSecretCache_SetAndGet(t *testing.T) {
	cache := NewSecretCache(0)
	secrets := map[string]string{"key": "value"}

	cache.Set("secret/path", secrets)
	got, ok := cache.Get("secret/path")

	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got["key"] != "value" {
		t.Errorf("expected 'value', got %q", got["key"])
	}
}

func TestSecretCache_MissOnUnknownPath(t *testing.T) {
	cache := NewSecretCache(0)

	_, ok := cache.Get("nonexistent/path")
	if ok {
		t.Fatal("expected cache miss for unknown path")
	}
}

func TestSecretCache_ExpiredEntry(t *testing.T) {
	cache := NewSecretCache(10 * time.Millisecond)
	cache.Set("secret/path", map[string]string{"k": "v"})

	time.Sleep(20 * time.Millisecond)

	_, ok := cache.Get("secret/path")
	if ok {
		t.Fatal("expected cache miss for expired entry")
	}
}

func TestSecretCache_NoExpiry(t *testing.T) {
	cache := NewSecretCache(0)
	cache.Set("secret/path", map[string]string{"k": "v"})

	time.Sleep(5 * time.Millisecond)

	_, ok := cache.Get("secret/path")
	if !ok {
		t.Fatal("expected cache hit for entry with no TTL")
	}
}

func TestSecretCache_Invalidate(t *testing.T) {
	cache := NewSecretCache(0)
	cache.Set("secret/path", map[string]string{"k": "v"})
	cache.Invalidate("secret/path")

	_, ok := cache.Get("secret/path")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestSecretCache_Len(t *testing.T) {
	cache := NewSecretCache(0)
	if cache.Len() != 0 {
		t.Fatalf("expected length 0, got %d", cache.Len())
	}

	cache.Set("path/a", map[string]string{})
	cache.Set("path/b", map[string]string{})

	if cache.Len() != 2 {
		t.Fatalf("expected length 2, got %d", cache.Len())
	}

	cache.Invalidate("path/a")
	if cache.Len() != 1 {
		t.Fatalf("expected length 1 after invalidation, got %d", cache.Len())
	}
}
