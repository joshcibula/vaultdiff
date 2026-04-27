package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCheckpoint_FieldsSet(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"key": "value"},
	}
	cp := NewCheckpoint("test", secrets)
	if cp.Name != "test" {
		t.Errorf("expected name 'test', got %q", cp.Name)
	}
	if cp.Secrets["secret/app"]["key"] != "value" {
		t.Error("expected secrets to be set")
	}
	if cp.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestSaveAndLoadCheckpoint_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	opts := CheckpointOptions{Enabled: true, Dir: dir, MaxAge: time.Hour}
	secrets := map[string]map[string]string{
		"secret/db": {"pass": "hunter2"},
	}
	cp := NewCheckpoint("roundtrip", secrets)
	if err := SaveCheckpoint(opts, cp); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}
	loaded, err := LoadCheckpoint(opts, "roundtrip")
	if err != nil {
		t.Fatalf("LoadCheckpoint: %v", err)
	}
	if loaded.Name != "roundtrip" {
		t.Errorf("expected name 'roundtrip', got %q", loaded.Name)
	}
	if loaded.Secrets["secret/db"]["pass"] != "hunter2" {
		t.Error("expected secret value to survive round-trip")
	}
}

func TestLoadCheckpoint_MissingFile(t *testing.T) {
	dir := t.TempDir()
	opts := DefaultCheckpointOptions()
	opts.Dir = dir
	_, err := LoadCheckpoint(opts, "nonexistent")
	if err == nil {
		t.Error("expected error for missing checkpoint")
	}
}

func TestLoadCheckpoint_Expired(t *testing.T) {
	dir := t.TempDir()
	opts := CheckpointOptions{Enabled: true, Dir: dir, MaxAge: time.Millisecond}
	cp := NewCheckpoint("old", nil)
	cp.CreatedAt = time.Now().Add(-time.Hour)
	if err := SaveCheckpoint(opts, cp); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}
	_, err := LoadCheckpoint(opts, "old")
	if err == nil {
		t.Error("expected expiry error")
	}
}

func TestSaveCheckpoint_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "checkpoints")
	opts := CheckpointOptions{Enabled: true, Dir: dir, MaxAge: time.Hour}
	cp := NewCheckpoint("init", nil)
	if err := SaveCheckpoint(opts, cp); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "init.json")); err != nil {
		t.Errorf("expected checkpoint file to exist: %v", err)
	}
}

func TestDefaultCheckpointOptions(t *testing.T) {
	opts := DefaultCheckpointOptions()
	if opts.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if opts.Dir == "" {
		t.Error("expected non-empty Dir")
	}
	if opts.MaxAge <= 0 {
		t.Error("expected positive MaxAge")
	}
}
