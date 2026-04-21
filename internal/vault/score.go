package vault

import (
	"fmt"
	"math"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// ScoreOptions controls how a diff similarity score is computed.
type ScoreOptions struct {
	// WeightModified is the penalty weight for modified keys (default 1.0).
	WeightModified float64
	// WeightOnlyInLeft is the penalty weight for keys only in left (default 1.0).
	WeightOnlyInLeft float64
	// WeightOnlyInRight is the penalty weight for keys only in right (default 1.0).
	WeightOnlyInRight float64
}

// DefaultScoreOptions returns sensible default weights.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		WeightModified:    1.0,
		WeightOnlyInLeft:  1.0,
		WeightOnlyInRight: 1.0,
	}
}

// DiffScore holds the computed similarity score for a diff result set.
type DiffScore struct {
	// Total number of keys considered across both sides.
	TotalKeys int
	// Penalty is the weighted sum of differences.
	Penalty float64
	// Similarity is a value in [0.0, 1.0] where 1.0 means identical.
	Similarity float64
	// Label is a human-readable description of the similarity.
	Label string
}

// String returns a formatted summary of the score.
func (s DiffScore) String() string {
	return fmt.Sprintf("similarity=%.2f (%s), total_keys=%d, penalty=%.2f",
		s.Similarity, s.Label, s.TotalKeys, s.Penalty)
}

// ScoreDiff computes a similarity score from a slice of diff.Result entries.
func ScoreDiff(results []diff.Result, opts ScoreOptions) DiffScore {
	if len(results) == 0 {
		return DiffScore{Similarity: 1.0, Label: "identical"}
	}

	totalKeys := 0
	penalty := 0.0

	for _, r := range results {
		totalKeys += len(r.Modified) + len(r.OnlyInLeft) + len(r.OnlyInRight)
		// unchanged keys contribute to total without penalty
		totalKeys += len(r.Unchanged)
		penalty += float64(len(r.Modified)) * opts.WeightModified
		penalty += float64(len(r.OnlyInLeft)) * opts.WeightOnlyInLeft
		penalty += float64(len(r.OnlyInRight)) * opts.WeightOnlyInRight
	}

	var similarity float64
	if totalKeys == 0 {
		similarity = 1.0
	} else {
		similarity = math.Max(0, 1.0-penalty/float64(totalKeys))
	}

	return DiffScore{
		TotalKeys:  totalKeys,
		Penalty:    penalty,
		Similarity: similarity,
		Label:      scoreLabel(similarity),
	}
}

func scoreLabel(s float64) string {
	switch {
	case s >= 0.95:
		return "identical"
	case s >= 0.75:
		return "similar"
	case s >= 0.50:
		return "diverged"
	default:
		return "very different"
	}
}
