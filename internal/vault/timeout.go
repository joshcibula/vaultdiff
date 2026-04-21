package vault

import (
	"context"
	"fmt"
	"time"
)

// TimeoutOptions configures per-operation timeouts for Vault requests.
type TimeoutOptions struct {
	// ListTimeout is the timeout for list (recursive scan) operations.
	ListTimeout time.Duration
	// ReadTimeout is the timeout for individual secret read operations.
	ReadTimeout time.Duration
	// TotalTimeout is an optional hard cap on the entire diff run.
	TotalTimeout time.Duration
}

// DefaultTimeoutOptions returns sensible defaults.
func DefaultTimeoutOptions() TimeoutOptions {
	return TimeoutOptions{
		ListTimeout:  15 * time.Second,
		ReadTimeout:  10 * time.Second,
		TotalTimeout: 0, // disabled by default
	}
}

// Validate returns an error if any timeout value is negative.
func (o TimeoutOptions) Validate() error {
	if o.ListTimeout < 0 {
		return fmt.Errorf("list-timeout must be >= 0, got %s", o.ListTimeout)
	}
	if o.ReadTimeout < 0 {
		return fmt.Errorf("read-timeout must be >= 0, got %s", o.ReadTimeout)
	}
	if o.TotalTimeout < 0 {
		return fmt.Errorf("total-timeout must be >= 0, got %s", o.TotalTimeout)
	}
	return nil
}

// WithListTimeout returns a context with the list timeout applied,
// or the parent context unchanged if ListTimeout is zero.
func (o TimeoutOptions) WithListTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if o.ListTimeout == 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, o.ListTimeout)
}

// WithReadTimeout returns a context with the read timeout applied,
// or the parent context unchanged if ReadTimeout is zero.
func (o TimeoutOptions) WithReadTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if o.ReadTimeout == 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, o.ReadTimeout)
}

// WithTotalTimeout returns a context with the total timeout applied,
// or the parent context unchanged if TotalTimeout is zero.
func (o TimeoutOptions) WithTotalTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if o.TotalTimeout == 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, o.TotalTimeout)
}

// IsZero reports whether all timeout values are zero (i.e. no timeouts configured).
func (o TimeoutOptions) IsZero() bool {
	return o.ListTimeout == 0 && o.ReadTimeout == 0 && o.TotalTimeout == 0
}
