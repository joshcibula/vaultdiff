package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of secrets at a given path.
type Snapshot struct {
	Path      string                       `json:"path"`
	Namespace string                       `json:"namespace,omitempty"`
	CapturedAt time.Time                   `json:"captured_at"`
	Secrets   map[string]map[string]string `json:"secrets"`
}

// NewSnapshot creates a Snapshot from the provided secrets map.
func NewSnapshot(path, namespace string, secrets map[string]map[string]string) *Snapshot {
	return &Snapshot{
		Path:       path,
		Namespace:  namespace,
		CapturedAt: time.Now().UTC(),
		Secrets:    secrets,
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at the given filepath.
func SaveSnapshot(s *Snapshot, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", filepath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from a JSON file at the given filepath.
func LoadSnapshot(filepath string) (*Snapshot, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", filepath, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
