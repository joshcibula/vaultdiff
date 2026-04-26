package vault

import (
	"fmt"
	"sort"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// DriftOptions configures drift detection behaviour.
type DriftOptions struct {
	// Enabled controls whether drift detection is active.
	Enabled bool
	// Threshold is the minimum percentage of changed keys (0-100) that
	// constitutes a "significant" drift event.
	Threshold float64
	// IgnorePaths is a list of path prefixes to exclude from drift analysis.
	IgnorePaths []string
}

// DefaultDriftOptions returns sensible defaults for drift detection.
func DefaultDriftOptions() DriftOptions {
	return DriftOptions{
		Enabled:   false,
		Threshold: 10.0,
	}
}

// DriftReport summarises the drift detected between two secret sets.
type DriftReport struct {
	DetectedAt   time.Time
	TotalKeys    int
	ChangedKeys  int
	DriftPercent float64
	Significant  bool
	ChangedPaths []string
}

// String returns a human-readable summary of the drift report.
func (r DriftReport) String() string {
	return fmt.Sprintf(
		"drift report [%s]: %d/%d keys changed (%.1f%%) significant=%v",
		r.DetectedAt.Format(time.RFC3339),
		r.ChangedKeys, r.TotalKeys, r.DriftPercent, r.Significant,
	)
}

// DetectDrift analyses diff results and produces a DriftReport.
func DetectDrift(results []diff.Result, opts DriftOptions) DriftReport {
	report := DriftReport{
		DetectedAt: time.Now().UTC(),
	}

	if !opts.Enabled || len(results) == 0 {
		return report
	}

	pathSet := make(map[string]struct{})

	for _, r := range results {
		if isIgnoredPath(r.Path, opts.IgnorePaths) {
			continue
		}
		report.TotalKeys++
		if r.Change != diff.NoChange {
			report.ChangedKeys++
			pathSet[r.Path] = struct{}{}
		}
	}

	if report.TotalKeys > 0 {
		report.DriftPercent = float64(report.ChangedKeys) / float64(report.TotalKeys) * 100
	}

	report.Significant = report.DriftPercent >= opts.Threshold

	for p := range pathSet {
		report.ChangedPaths = append(report.ChangedPaths, p)
	}
	sort.Strings(report.ChangedPaths)

	return report
}

func isIgnoredPath(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}
	return false
}
