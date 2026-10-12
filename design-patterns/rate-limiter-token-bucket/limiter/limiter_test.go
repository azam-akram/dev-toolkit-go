package tokenbucket

import (
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{
		now: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

func newTestLimiter(
	t *testing.T,
	capacity int,
	refillRate float64,
	clock Clock,
) *TokenBucket {
	t.Helper()

	limiter, err := New(
		Config{
			Capacity:   capacity,
			RefillRate: refillRate,
		},
		clock,
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return limiter
}

func TestBurstCapacity(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 2, 1, clock)

	if !limiter.Allow() || !limiter.Allow() {
		t.Fatal("expected the initial two requests to be allowed")
	}

	if limiter.Allow() {
		t.Fatal("expected third request to be rejected")
	}
}

func TestRefill(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 2, 1, clock)

	limiter.Allow()
	limiter.Allow()

	if limiter.Allow() {
		t.Fatal("expected empty bucket")
	}

	clock.Advance(time.Second)

	if !limiter.Allow() {
		t.Fatal("expected one token after one second")
	}

	if limiter.Allow() {
		t.Fatal("expected bucket to be empty again")
	}
}

func TestRefillDoesNotExceedCapacity(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 2, 1, clock)

	clock.Advance(time.Hour)

	if !limiter.Allow() || !limiter.Allow() {
		t.Fatal("expected bucket to be full, up to capacity")
	}

	if limiter.Allow() {
		t.Fatal("bucket must not exceed capacity")
	}
}

func TestAllowNIsAtomic(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 3, 1, clock)

	if !limiter.AllowN(2) {
		t.Fatal("expected two tokens to be available")
	}

	if limiter.AllowN(2) {
		t.Fatal("expected insufficient tokens")
	}

	// The failed request must not have consumed its partial amount.
	if !limiter.Allow() {
		t.Fatal("expected one token to remain")
	}
}

func TestAllowNRejectsInvalidAmounts(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 3, 1, clock)

	for _, n := range []int{0, -1, 4} {
		if limiter.AllowN(n) {
			t.Errorf("AllowN(%d) = true; want false", n)
		}
	}
}

func TestInvalidConfig(t *testing.T) {
	tests := []Config{
		{Capacity: 0, RefillRate: 1},
		{Capacity: -1, RefillRate: 1},
		{Capacity: 1, RefillRate: 0},
		{Capacity: 1, RefillRate: -1},
		{Capacity: 1, RefillRate: math.NaN()},
		{Capacity: 1, RefillRate: math.Inf(1)},
	}

	for _, cfg := range tests {
		_, err := New(cfg, nil)
		if err != ErrInvalidConfig {
			t.Errorf("New(%+v) error = %v; want ErrInvalidConfig", cfg, err)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 1000, 1, clock)

	const workers = 100
	const attemptsPerWorker = 20

	var wg sync.WaitGroup
	var countMu sync.Mutex
	allowed := 0

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < attemptsPerWorker; j++ {
				if limiter.Allow() {
					countMu.Lock()
					allowed++
					countMu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if allowed != 1000 {
		t.Fatalf("allowed %d requests; want 1000", allowed)
	}
}

func TestHTTPMiddleware(t *testing.T) {
	clock := newFakeClock()
	limiter := newTestLimiter(t, 1, 1, clock)

	calls := 0
	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		calls++
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(limiter)(next)

	first := httptest.NewRecorder()
	handler.ServeHTTP(
		first,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)

	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d; want 200", first.Code)
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(
		second,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)

	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d; want 429", second.Code)
	}

	if calls != 1 {
		t.Fatalf("handler called %d times; want 1", calls)
	}
}
