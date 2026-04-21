package vault

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestDefaultScoreOptions(t *testing.T) {
	opts := DefaultScoreOptions()
	if opts.WeightModified != 1.0 || opts.WeightOnlyInLeft != 1.0 || opts.WeightOnlyInRight != 1.0 {
		t.Fatalf("unexpected default weights: %+v", opts)
	}
}

func TestScoreDiff_EmptyResults(t *testing.T) {
	score := ScoreDiff(nil, DefaultScoreOptions())
	if score.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0, got %.2f", score.Similarity)
	}
	if score.Label != "identical" {
		t.Errorf("expected label 'identical', got %q", score.Label)
	}
}

func TestScoreDiff_Identical(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/a", Unchanged: map[string]string{"key": "val"}},
	}
	score := ScoreDiff(results, DefaultScoreOptions())
	if score.Similarity != 1.0 {
		t.Errorf("expected 1.0 similarity, got %.2f", score.Similarity)
	}
	if score.TotalKeys != 1 {
		t.Errorf("expected TotalKeys=1, got %d", score.TotalKeys)
	}
}

func TestScoreDiff_AllModified(t *testing.T) {
	results := []diff.Result{
		{
			Path:     "secret/a",
			Modified: map[string][2]string{"key": {"old", "new"}},
		},
	}
	score := ScoreDiff(results, DefaultScoreOptions())
	if score.Similarity != 0.0 {
		t.Errorf("expected 0.0 similarity, got %.2f", score.Similarity)
	}
	if score.Label != "very different" {
		t.Errorf("expected 'very different', got %q", score.Label)
	}
}

func TestScoreDiff_PartialDiff(t *testing.T) {
	results := []diff.Result{
		{
			Path:      "secret/a",
			Unchanged: map[string]string{"a": "1", "b": "2", "c": "3"},
			Modified:  map[string][2]string{"d": {"x", "y"}},
		},
	}
	score := ScoreDiff(results, DefaultScoreOptions())
	// 4 total keys, 1 penalty => similarity = 0.75
	if score.TotalKeys != 4 {
		t.Errorf("expected TotalKeys=4, got %d", score.TotalKeys)
	}
	if score.Similarity != 0.75 {
		t.Errorf("expected similarity 0.75, got %.2f", score.Similarity)
	}
	if score.Label != "similar" {
		t.Errorf("expected label 'similar', got %q", score.Label)
	}
}

func TestDiffScore_String(t *testing.T) {
	s := DiffScore{TotalKeys: 10, Penalty: 2.0, Similarity: 0.80, Label: "similar"}
	out := s.String()
	if out == "" {
		t.Error("expected non-empty string from DiffScore.String()")
	}
}

func TestScoreLabel_Boundaries(t *testing.T) {
	cases := []struct {
		score float64
		want  string
	}{
		{1.00, "identical"},
		{0.95, "identical"},
		{0.94, "similar"},
		{0.75, "similar"},
		{0.74, "diverged"},
		{0.50, "diverged"},
		{0.49, "very different"},
		{0.00, "very different"},
	}
	for _, tc := range cases {
		got := scoreLabel(tc.score)
		if got != tc.want {
			t.Errorf("scoreLabel(%.2f) = %q, want %q", tc.score, got, tc.want)
		}
	}
}
