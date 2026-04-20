package vault

import (
	"sync"
	"testing"
	"time"
)

func TestDefaultRateLimitOptions(t *testing.T) {
	opts := DefaultRateLimitOptions()
	if opts.RequestsPerSecond <= 0 {
		t.Errorf("expected positive RequestsPerSecond, got %f", opts.RequestsPerSecond)
	}
	if opts.Burst <= 0 {
		t.Errorf("expected positive Burst, got %f", opts.Burst)
	}
}

func TestNewRateLimiter_InitialTokens(t *testing.T) {
	opts := RateLimitOptions{RequestsPerSecond: 5, Burst: 10}
	rl := NewRateLimiter(opts)
	if rl.tokens != 10 {
		t.Errorf("expected initial tokens=10, got %f", rl.tokens)
	}
}

func TestTryAcquire_ConsumesToken(t *testing.T) {
	opts := RateLimitOptions{RequestsPerSecond: 1, Burst: 3}
	rl := NewRateLimiter(opts)

	for i := 0; i < 3; i++ {
		if !rl.TryAcquire() {
			t.Fatalf("expected TryAcquire to succeed on attempt %d", i+1)
		}
	}
	// Burst exhausted — next acquire should fail immediately.
	if rl.TryAcquire() {
		t.Error("expected TryAcquire to fail after burst exhausted")
	}
}

func TestTryAcquire_RefillsOverTime(t *testing.T) {
	opts := RateLimitOptions{RequestsPerSecond: 100, Burst: 1}
	rl := NewRateLimiter(opts)

	// Consume the single burst token.
	if !rl.TryAcquire() {
		t.Fatal("expected first TryAcquire to succeed")
	}
	// Wait long enough for refill (100 rps => 10ms per token).
	time.Sleep(20 * time.Millisecond)
	if !rl.TryAcquire() {
		t.Error("expected TryAcquire to succeed after refill window")
	}
}

func TestWait_DoesNotExceedBurst(t *testing.T) {
	opts := RateLimitOptions{RequestsPerSecond: 50, Burst: 5}
	rl := NewRateLimiter(opts)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.Wait()
		}()
	}
	wg.Wait()
	// All 5 burst tokens should be consumed quickly (< 100ms).
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Errorf("burst of 5 took too long: %v", elapsed)
	}
}
