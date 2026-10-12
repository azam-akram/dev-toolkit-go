package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrCircuitHalfOpen = errors.New("circuit breaker is half-open")
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

type CircuitBreaker struct {
	mu            sync.Mutex
	state         State
	config        Config
	failures      []time.Time
	probeInFlight bool
	openedAt      time.Time
	now           func() time.Time
}

func NewCircuitBreaker(config Config) (*CircuitBreaker, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &CircuitBreaker{
		state:  StateClosed,
		config: config,
		now:    time.Now,
	}, nil
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return cb.state
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	err := fn(ctx)

	cb.afterRequest(err)

	return err
}

func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := cb.now()

	switch cb.state {
	case StateClosed:
		cb.pruneFailures(now)
		return nil

	case StateOpen:
		// Has the recovery period elapsed?
		if now.Sub(cb.openedAt) < cb.config.OpenDuration {
			return ErrCircuitOpen
		}

		// Move to HALF_OPEN.
		cb.state = StateHalfOpen

		// Only one request is allowed as the probe.
		if cb.probeInFlight {
			return ErrCircuitHalfOpen
		}

		cb.probeInFlight = true
		return nil

	case StateHalfOpen:
		// Only one probe request should be allowed.
		if cb.probeInFlight {
			return ErrCircuitHalfOpen
		}

		cb.probeInFlight = true
		return nil

	default:
		return fmt.Errorf("unknown circuit breaker state: %d", cb.state)
	}
}

func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := cb.now()

	switch cb.state {
	case StateClosed:
		if err == nil {
			// Successful request.
			return
		}

		cb.recordFailure(now)

		if len(cb.failures) >= cb.config.FailureThreshold {
			cb.open(now)
		}

	case StateHalfOpen:
		cb.probeInFlight = false

		if err == nil {
			// Probe succeeded. Service appears healthy again.
			cb.close()
			return
		}

		// Probe failed. Service is still unhealthy.
		cb.open(now)

	case StateOpen:
		// Normally requests should not execute while OPEN.
		// This case is defensive.
	}
}

func (cb *CircuitBreaker) recordFailure(now time.Time) {
	cb.pruneFailures(now)
	cb.failures = append(cb.failures, now)
}

func (cb *CircuitBreaker) pruneFailures(now time.Time) {
	cutoff := now.Add(-cb.config.FailureWindow)

	// Find the first failure still inside the window.
	i := 0
	for i < len(cb.failures) && cb.failures[i].Before(cutoff) {
		i++
	}

	if i > 0 {
		cb.failures = cb.failures[i:]
	}
}

func (cb *CircuitBreaker) open(now time.Time) {
	cb.state = StateOpen
	cb.openedAt = now
	cb.probeInFlight = false
}

func (cb *CircuitBreaker) close() {
	cb.state = StateClosed
	cb.openedAt = time.Time{}
	cb.failures = cb.failures[:0]
	cb.probeInFlight = false
}
