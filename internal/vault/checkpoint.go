package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// DefaultCheckpointOptions returns a CheckpointOptions with sensible defaults.
func DefaultCheckpointOptions() CheckpointOptions {
	return CheckpointOptions{
		Enabled: false,
		Dir:     ".vaultdiff",
		MaxAge:  24 * time.Hour,
	}
}

// CheckpointOptions controls checkpoint persistence behaviour.
type CheckpointOptions struct {
	Enabled bool
	Dir     string
	MaxAge  time.Duration
}

// Checkpoint records a named snapshot of secrets at a point in time.
type Checkpoint struct {
	Name      string                       `json:"name"`
	CreatedAt time.Time                    `json:"created_at"`
	Secrets   map[string]map[string]string `json:"secrets"`
}

// NewCheckpoint creates a Checkpoint with the current timestamp.
func NewCheckpoint(name string, secrets map[string]map[string]string) *Checkpoint {
	return &Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Secrets:   secrets,
	}
}

// SaveCheckpoint writes a checkpoint to dir/<name>.json.
func SaveCheckpoint(opts CheckpointOptions, cp *Checkpoint) error {
	if err := os.MkdirAll(opts.Dir, 0o700); err != nil {
		return fmt.Errorf("checkpoint: create dir: %w", err)
	}
	path := checkpointPath(opts.Dir, cp.Name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("checkpoint: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cp); err != nil {
		return fmt.Errorf("checkpoint: encode: %w", err)
	}
	return nil
}

// LoadCheckpoint reads a checkpoint from dir/<name>.json.
// It returns an error if the checkpoint is older than opts.MaxAge (when MaxAge > 0).
func LoadCheckpoint(opts CheckpointOptions, name string) (*Checkpoint, error) {
	path := checkpointPath(opts.Dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("checkpoint %q not found", name)
		}
		return nil, fmt.Errorf("checkpoint: read: %w", err)
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("checkpoint: decode: %w", err)
	}
	if opts.MaxAge > 0 && time.Since(cp.CreatedAt) > opts.MaxAge {
		return nil, fmt.Errorf("checkpoint %q expired (age %s > max %s)", name, time.Since(cp.CreatedAt).Round(time.Second), opts.MaxAge)
	}
	return &cp, nil
}

func checkpointPath(dir, name string) string {
	return fmt.Sprintf("%s/%s.json", dir, name)
}
