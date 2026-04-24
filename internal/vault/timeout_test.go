package vault

import (
	"context"
	"testing"
	"time"
)

func TestDefaultTimeoutOptions(t *testing.T) {
	opts := DefaultTimeoutOptions()
	if opts.ListTimeout != 15*time.Second {
		t.Errorf("expected ListTimeout 15s, got %s", opts.ListTimeout)
	}
	if opts.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout 10s, got %s", opts.ReadTimeout)
	}
	if opts.TotalTimeout != 0 {
		t.Errorf("expected TotalTimeout 0 (disabled), got %s", opts.TotalTimeout)
	}
}

func TestTimeoutOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    TimeoutOptions
		wantErr bool
	}{
		{"valid defaults", DefaultTimeoutOptions(), false},
		{"all zero", TimeoutOptions{}, false},
		{"negative list", TimeoutOptions{ListTimeout: -1 * time.Second}, true},
		{"negative read", TimeoutOptions{ReadTimeout: -1 * time.Second}, true},
		{"negative total", TimeoutOptions{TotalTimeout: -1 * time.Second}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestWithListTimeout_Applied(t *testing.T) {
	opts := TimeoutOptions{ListTimeout: 5 * time.Second}
	ctx, cancel := opts.WithListTimeout(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 5*time.Second || time.Until(deadline) <= 0 {
		t.Errorf("unexpected deadline: %v", deadline)
	}
}

func TestWithListTimeout_ZeroNoDeadline(t *testing.T) {
	opts := TimeoutOptions{ListTimeout: 0}
	ctx, cancel := opts.WithListTimeout(context.Background())
	defer cancel()
	_, ok := ctx.Deadline()
	if ok {
		t.Error("expected no deadline when ListTimeout is zero")
	}
}

func TestWithTotalTimeout_Applied(t *testing.T) {
	opts := TimeoutOptions{TotalTimeout: 30 * time.Second}
	ctx, cancel := opts.WithTotalTimeout(context.Background())
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set for total timeout")
	}
}

func TestWithReadTimeout_ZeroNoDeadline(t *testing.T) {
	opts := TimeoutOptions{ReadTimeout: 0}
	ctx, cancel := opts.WithReadTimeout(context.Background())
	defer cancel()
	_, ok := ctx.Deadline()
	if ok {
		t.Error("expected no deadline when ReadTimeout is zero")
	}
}

func TestWithReadTimeout_Applied(t *testing.T) {
	opts := TimeoutOptions{ReadTimeout: 10 * time.Second}
	ctx, cancel := opts.WithReadTimeout(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 10*time.Second || time.Until(deadline) <= 0 {
		t.Errorf("unexpected deadline: %v", deadline)
	}
}
