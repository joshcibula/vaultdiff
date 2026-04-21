package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSnapshot_FieldsSet(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/foo": {"key": "value"},
	}
	before := time.Now().UTC()
	s := NewSnapshot("secret/", "dev", secrets)
	after := time.Now().UTC()

	if s.Path != "secret/" {
		t.Errorf("expected path %q, got %q", "secret/", s.Path)
	}
	if s.Namespace != "dev" {
		t.Errorf("expected namespace %q, got %q", "dev", s.Namespace)
	}
	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v out of expected range", s.CapturedAt)
	}
	if len(s.Secrets) != 1 {
		t.Errorf("expected 1 secret path, got %d", len(s.Secrets))
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/alpha": {"username": "admin", "password": "s3cr3t"},
		"secret/beta":  {"token": "abc123"},
	}
	orig := NewSnapshot("secret/", "staging", secrets)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := SaveSnapshot(orig, path); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if loaded.Path != orig.Path {
		t.Errorf("path mismatch: got %q want %q", loaded.Path, orig.Path)
	}
	if loaded.Namespace != orig.Namespace {
		t.Errorf("namespace mismatch: got %q want %q", loaded.Namespace, orig.Namespace)
	}
	if len(loaded.Secrets) != len(orig.Secrets) {
		t.Errorf("secrets length mismatch: got %d want %d", len(loaded.Secrets), len(orig.Secrets))
	}
	if loaded.Secrets["secret/alpha"]["password"] != "s3cr3t" {
		t.Errorf("unexpected secret value: %q", loaded.Secrets["secret/alpha"]["password"])
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error loading missing file, got nil")
	}
}

func TestSaveSnapshot_UnwritablePath(t *testing.T) {
	s := NewSnapshot("secret/", "", map[string]map[string]string{})
	err := SaveSnapshot(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Fatal("expected error saving to unwritable path, got nil")
	}
}

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "out.json")
	s := NewSnapshot("kv/", "prod", map[string]map[string]string{})

	if err := SaveSnapshot(s, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist after SaveSnapshot")
	}
}
