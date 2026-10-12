package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCircuitOpensAfterFailureThreshold(t *testing.T) {
	cb, err := NewCircuitBreaker(Config{
		FailureThreshold: 3,
		FailureWindow:    10 * time.Minute,
		OpenDuration:     30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	failure := errors.New("service unavailable")

	for i := 0; i < 3; i++ {
		err := cb.Execute(context.Background(), func(context.Context) error {
			return failure
		})

		if !errors.Is(err, failure) {
			t.Fatalf("expected failure, got %v", err)
		}
	}

	if got := cb.State(); got != StateOpen {
		t.Fatalf("expected OPEN, got %s", got)
	}

	err = cb.Execute(context.Background(), func(context.Context) error {
		t.Fatal("request should not have reached downstream")
		return nil
	})

	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestHalfOpenSuccessfulProbeClosesCircuit(t *testing.T) {
	cb, err := NewCircuitBreaker(Config{
		FailureThreshold: 1,
		FailureWindow:    10 * time.Minute,
		OpenDuration:     30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Force the circuit OPEN.
	_ = cb.Execute(context.Background(), func(context.Context) error {
		return errors.New("failure")
	})

	if cb.State() != StateOpen {
		t.Fatal("expected OPEN")
	}

	// Advance time beyond OpenDuration.
	currentTime := time.Now()

	cb.mu.Lock()
	cb.now = func() time.Time {
		return currentTime.Add(31 * time.Second)
	}
	cb.mu.Unlock()

	// Successful probe.
	err = cb.Execute(context.Background(), func(context.Context) error {
		return nil
	})

	if err != nil {
		t.Fatalf("expected successful probe, got %v", err)
	}

	if got := cb.State(); got != StateClosed {
		t.Fatalf("expected CLOSED, got %s", got)
	}
}

func TestHalfOpenFailedProbeReopensCircuit(t *testing.T) {
	cb, err := NewCircuitBreaker(Config{
		FailureThreshold: 1,
		FailureWindow:    10 * time.Minute,
		OpenDuration:     30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	// CLOSED -> OPEN.
	_ = cb.Execute(context.Background(), func(context.Context) error {
		return errors.New("failure")
	})

	currentTime := time.Now()

	cb.mu.Lock()
	cb.now = func() time.Time {
		return currentTime.Add(31 * time.Second)
	}
	cb.mu.Unlock()

	// HALF_OPEN -> OPEN.
	_ = cb.Execute(context.Background(), func(context.Context) error {
		return errors.New("still failing")
	})

	if got := cb.State(); got != StateOpen {
		t.Fatalf("expected OPEN, got %s", got)
	}
}

func TestFailuresOutsideWindowAreIgnored(t *testing.T) {
	cb, err := NewCircuitBreaker(Config{
		FailureThreshold: 3,
		FailureWindow:    10 * time.Minute,
		OpenDuration:     30 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	currentTime := time.Now()

	cb.mu.Lock()
	cb.now = func() time.Time {
		return currentTime
	}
	cb.mu.Unlock()

	failure := errors.New("failure")

	// First failure.
	_ = cb.Execute(context.Background(), func(context.Context) error {
		return failure
	})

	// Move forward 11 minutes.
	cb.mu.Lock()
	cb.now = func() time.Time {
		return currentTime.Add(11 * time.Minute)
	}
	cb.mu.Unlock()

	// Second failure should not count with the first one.
	_ = cb.Execute(context.Background(), func(context.Context) error {
		return failure
	})

	if got := cb.State(); got != StateClosed {
		t.Fatalf("expected CLOSED, got %s", got)
	}
}
