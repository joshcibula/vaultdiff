package vault

import (
	"errors"
	"time"
)

// RetryOptions configures retry behaviour for Vault API calls.
type RetryOptions struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryOptions returns sensible defaults for retrying Vault requests.
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxAttempts:  3,
		InitialDelay: 250 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableError wraps an error to signal that the operation may be retried.
type RetryableError struct {
	Cause error
}

func (e *RetryableError) Error() string {
	return "retryable: " + e.Cause.Error()
}

func (e *RetryableError) Unwrap() error { return e.Cause }

// IsRetryable reports whether err should trigger a retry.
func IsRetryable(err error) bool {
	var re *RetryableError
	return errors.As(err, &re)
}

// WithRetry executes fn up to opts.MaxAttempts times, backing off between
// attempts. Only errors that satisfy IsRetryable are retried; all others are
// returned immediately.
func WithRetry(opts RetryOptions, fn func() error) error {
	if opts.MaxAttempts <= 0 {
		opts.MaxAttempts = 1
	}

	delay := opts.InitialDelay
	var lastErr error

	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		if !IsRetryable(err) {
			return err
		}
		lastErr = err
		if attempt < opts.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * opts.Multiplier)
			if delay > opts.MaxDelay {
				delay = opts.MaxDelay
			}
		}
	}
	return lastErr
}
