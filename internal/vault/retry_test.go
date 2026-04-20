package vault

import (
	"errors"
	"testing"
	"time"
)

func TestWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := WithRetry(DefaultRetryOptions(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_NonRetryableErrorStopsImmediately(t *testing.T) {
	sentinel := errors.New("fatal")
	calls := 0
	err := WithRetry(DefaultRetryOptions(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_RetriesUpToMaxAttempts(t *testing.T) {
	opts := RetryOptions{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
	calls := 0
	inner := errors.New("transient")
	err := WithRetry(opts, func() error {
		calls++
		return &RetryableError{Cause: inner}
	})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if !IsRetryable(err) {
		t.Fatalf("expected retryable error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWithRetry_SucceedsOnSecondAttempt(t *testing.T) {
	opts := RetryOptions{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
	calls := 0
	err := WithRetry(opts, func() error {
		calls++
		if calls < 2 {
			return &RetryableError{Cause: errors.New("not yet")}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestIsRetryable(t *testing.T) {
	if IsRetryable(errors.New("plain")) {
		t.Error("plain error should not be retryable")
	}
	if !IsRetryable(&RetryableError{Cause: errors.New("x")}) {
		t.Error("RetryableError should be retryable")
	}
}

func TestDefaultRetryOptions(t *testing.T) {
	opts := DefaultRetryOptions()
	if opts.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", opts.MaxAttempts)
	}
	if opts.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", opts.Multiplier)
	}
}
