package vault

import (
	"fmt"
	"time"
)

// BaselineOptions controls how a baseline comparison is performed.
type BaselineOptions struct {
	// Enabled activates baseline comparison mode.
	Enabled bool
	// SnapshotPath is the path to a previously saved snapshot used as the baseline.
	SnapshotPath string
	// MaxAgeDays rejects snapshots older than this many days (0 = no limit).
	MaxAgeDays int
}

// DefaultBaselineOptions returns sensible defaults.
func DefaultBaselineOptions() BaselineOptions {
	return BaselineOptions{
		Enabled:     false,
		SnapshotPath: "",
		MaxAgeDays:  0,
	}
}

// BaselineError is returned when a baseline snapshot cannot be used.
type BaselineError struct {
	Msg string
}

func (e *BaselineError) Error() string {
	return fmt.Sprintf("baseline error: %s", e.Msg)
}

// LoadBaseline loads a snapshot from disk and validates it against opts.
// Returns the snapshot's secrets map or an error.
func LoadBaseline(opts BaselineOptions) (map[string]map[string]string, error) {
	if !opts.Enabled {
		return nil, &BaselineError{Msg: "baseline mode is not enabled"}
	}
	if opts.SnapshotPath == "" {
		return nil, &BaselineError{Msg: "snapshot path is required for baseline comparison"}
	}

	snap, err := LoadSnapshot(opts.SnapshotPath)
	if err != nil {
		return nil, fmt.Errorf("loading baseline snapshot: %w", err)
	}

	if opts.MaxAgeDays > 0 {
		maxAge := time.Duration(opts.MaxAgeDays) * 24 * time.Hour
		age := time.Since(snap.CreatedAt)
		if age > maxAge {
			return nil, &BaselineError{
				Msg: fmt.Sprintf("snapshot is %.1f days old, exceeds max age of %d days",
					age.Hours()/24, opts.MaxAgeDays),
			}
		}
	}

	return snap.Secrets, nil
}

// IsBaselineError reports whether err is a BaselineError.
func IsBaselineError(err error) bool {
	_, ok := err.(*BaselineError)
	return ok
}
