# Circuit Breaker – Interview Review Questions

Mock code-review / pair-programming questions for the circuit breaker implementation in this directory.
Answer out loud, as you would in the room.

## Warm-up (design intent)

1. Walk me through the state machine in your own words. Why did you choose a *sliding window of failure timestamps* rather than a simple consecutive-failure counter? What does that cost you at high traffic?

   **Answer**

   *State machine*

   - **CLOSED**: all requests pass through. Each failure is timestamped. When the number of failures inside `FailureWindow` reaches `FailureThreshold`, the breaker moves to OPEN.
   - **OPEN**: requests are rejected immediately with `ErrCircuitOpen` and the downstream is not called. After `OpenDuration` has elapsed, the next request moves the breaker to HALF_OPEN. The transition is lazy, so there is no timer.
   - **HALF_OPEN**: exactly one probe request is allowed through. If it succeeds, the breaker goes to CLOSED and the failure history is cleared. If it fails, the breaker goes back to OPEN and the timer restarts. Other requests are rejected with `ErrCircuitHalfOpen` while the probe is in flight.

   *Why a sliding window and not a plain counter*

   A cumulative counter never forgets. Failures from hours ago add up with a failure from now and trip the breaker without any real outage. The window makes failures expire, so only recent failures count.

   *Why not a consecutive-failure counter*

   A consecutive counter also has no stale-failure problem, because any success resets it to zero. The cost is that a flaky downstream that fails 9 of every 10 calls never trips it: a single success wipes the streak, so the breaker misses a clearly unhealthy service. The sliding window tolerates interleaved successes. It measures failures *per unit time*, not failures *in a row*.

   *What my code actually does*

   In the `StateClosed` branch of `afterRequest`, a success just returns and does not clear `cb.failures`. So this is an "N failures within a time window" breaker, not a consecutive one. With threshold 5 and the sequence fail, fail, fail, fail, success, fail, the breaker trips on the 6th call, because 5 failures sit inside the window and the success didn't reset anything. That behaviour is deliberate, and I should be able to say so out loud.

   *Cost at high traffic*

   1. **A raw count doesn't scale with volume.** The threshold is an absolute number, not a rate. With 5 failures per 10 minutes at 10,000 RPS, a window holds 6,000,000 requests. A healthy 0.001% error rate means about 60 failures per window, so the breaker trips roughly every 50 seconds on a perfectly healthy service. The reverse also happens: at 1 RPS, 5 failures in 10 minutes is an 83% error rate, which is a real outage. The same config is far too sensitive on busy services and reasonable on quiet ones.
   2. **The fix is a failure *ratio* with a minimum sample size.** Trip when `failures / total >= X%` and `total >= minRequests` within the window. That means tracking successes too. At volume, use fixed time buckets (for example 10 buckets of 1s in a ring buffer) instead of a per-failure timestamp slice. That gives O(1) memory and O(buckets) pruning, with no slice shifting.
   3. **Lock contention.** Every request takes the single mutex twice (`beforeRequest` and `afterRequest`), so the breaker becomes a serialization point across all goroutines. Options are atomic counters, sharding, or sampling. The critical sections are tiny, so it is fine until a profile says otherwise.
   4. **Memory is bounded but not free.** `failures` can't grow past `FailureThreshold` entries because the breaker opens at the threshold. But `cb.failures[i:]` keeps the old backing array alive and `append` may reallocate. That is acceptable at this size, and a ring buffer would be cleaner.

2. In `Execute` you take the lock in `beforeRequest`, release it, run `fn`, then take it again in `afterRequest`. Why not hold the lock across the call? What are you giving up by releasing it?

   **Answer**

   *Why the lock must not be held across `fn`*

   `fn` is a network call that can take seconds. Holding the mutex across it would serialize every request to that downstream to concurrency 1, and `State()` would block too. The breaker would itself become the outage. A mutex should only protect the breaker's own bookkeeping, which takes microseconds.

   *What releasing it costs*

   The decision (`beforeRequest`) and the outcome (`afterRequest`) are no longer atomic. Between them, other goroutines can change the state. Concretely:

   - `afterRequest` switches on `cb.state` *as it is now*, not as it was when this request was admitted. A result can be applied to the wrong state (this is question 3).
   - In CLOSED there is no limit on in-flight requests. 1,000 requests can be admitted before the first failure is recorded, so the breaker reacts one round-trip late. That is inherent to the pattern, and a bulkhead or semaphore is the tool for bounding concurrency.

   *The fix*

   Keep the short critical sections but make the outcome carry its context. `beforeRequest` returns a **generation** number, which is bumped on every state transition. `afterRequest(gen, err)` ignores the result if `gen != cb.generation`. This is what `sony/gobreaker` does.

## Digging into the code

3. Say 50 requests are in flight while the circuit is CLOSED. The 5th failure trips it to OPEN. After 30s a probe goes out and the circuit is HALF_OPEN. Now one of the *old* in-flight requests finally returns. Walk me through what `afterRequest` does with it. Is that the behavior you want?

   **Answer**

   *What happens*

   `afterRequest` has no idea the request is old. It reads `cb.state`, which is now `StateHalfOpen`, and takes that branch:

   1. It sets `probeInFlight = false`, even though the *real* probe is still running. A second probe can now be admitted, so the "only one probe" guarantee is broken.
   2. If the stale request succeeded, it calls `close()`. The circuit closes on the evidence of a request that started before the outage was even detected, and the real probe's result is discarded or applied to a CLOSED state.
   3. If the stale request failed, it calls `open(now)`, which resets `openedAt` and cancels a perfectly good probe.

   So no, it isn't the behavior I want. HALF_OPEN is supposed to be decided only by the probe.

   *Fix*

   A generation counter, incremented in every transition (`open`, `close`, and the OPEN to HALF_OPEN move):

   ```go
   func (cb *CircuitBreaker) beforeRequest() (gen uint64, err error) { /* ... */ return cb.generation, nil }

   func (cb *CircuitBreaker) afterRequest(gen uint64, err error) {
       cb.mu.Lock()
       defer cb.mu.Unlock()
       if gen != cb.generation {
           return // result belongs to an earlier state, ignore it
       }
       // ... existing switch ...
   }
   ```

   Requests that finish in OPEN are already ignored by the defensive branch. The generation check makes the same rule hold for HALF_OPEN and for a re-opened and re-closed circuit.

4. Look at `beforeRequest`, the `StateOpen` case. You set `cb.state = StateHalfOpen` and then check `cb.probeInFlight`. Can that branch ever be true? What does that tell you about how you'd restructure it?

   **Answer**

   *Can it be true?*

   No, it is dead code. `probeInFlight` is only set to `true` while in HALF_OPEN. Every exit from HALF_OPEN goes through `open()` or `close()`, and both reset it to `false`. So by the time we are in OPEN and elapsed time is over, the flag is always `false`. The check was written defensively, but it hides the fact that I wasn't sure of my own invariant.

   *What it says about the structure*

   The OPEN-to-HALF_OPEN transition and the "admit one probe" logic are duplicated across two cases. I'd separate the two concerns: **first** resolve the current state (lazy transition), **then** make the admission decision once.

   ```go
   func (cb *CircuitBreaker) currentState(now time.Time) State {
       if cb.state == StateOpen && now.Sub(cb.openedAt) >= cb.config.OpenDuration {
           cb.state = StateHalfOpen
           cb.generation++
       }
       return cb.state
   }

   func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
       cb.mu.Lock()
       defer cb.mu.Unlock()

       switch cb.currentState(cb.now()) {
       case StateOpen:
           return 0, ErrCircuitOpen
       case StateHalfOpen:
           if cb.probeInFlight {
               return 0, ErrCircuitHalfOpen
           }
           cb.probeInFlight = true
       }
       return cb.generation, nil
   }
   ```

   `State()` should use the same helper. Today it reports OPEN even after `OpenDuration` has passed, which misleads dashboards.

5. What happens to the breaker if `fn` panics while the circuit is HALF_OPEN? How would you find that bug in production, and what's the fix?

   **Answer**

   *What happens*

   `afterRequest` is never called, so `probeInFlight` stays `true` and the state stays HALF_OPEN. Every later call hits `if cb.probeInFlight { return ErrCircuitHalfOpen }`. **The breaker is stuck permanently** and rejects 100% of traffic until the process restarts. The panic itself may be recovered by a middleware (for example `net/http`'s server recover), which makes it worse, because the process stays alive and stays broken.

   *How I'd find it in production*

   - Rejection metric for `ErrCircuitHalfOpen` is high and constant, with no transitions out of HALF_OPEN.
   - A "time in HALF_OPEN" gauge that keeps growing. A healthy probe takes at most one request timeout.
   - The downstream's own traffic drops to zero while it is healthy. The stack trace of the original panic is in the logs.

   *Fix*

   Guarantee that `afterRequest` always runs, with `defer` plus `recover`, and re-panic so the behaviour is unchanged for the caller:

   ```go
   func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) (err error) {
       gen, err := cb.beforeRequest()
       if err != nil {
           return err
       }

       defer func() {
           if r := recover(); r != nil {
               cb.afterRequest(gen, fmt.Errorf("panic: %v", r))
               panic(r)
           }
       }()

       err = fn(ctx)
       cb.afterRequest(gen, err)
       return err
   }
   ```

   A related case is a probe that *hangs* because `fn` ignores `ctx`. Same symptom, same stuck state. The mitigation is to store `probeStartedAt` and allow a new probe once it is older than a probe timeout, or to enforce a deadline around `fn` inside `Execute`.

## Behavior and correctness

6. In `client.go`, a caller cancels their `ctx` mid-request, say because the user closed the browser tab. How does your breaker account for that? Should a cancellation count as a downstream failure? What about a 429 or a 503 with `Retry-After`?

   **Answer**

   *Cancellation*

   Today any non-nil error counts. `httpClient.Do` returns a `context canceled` error, the client wraps it in `ErrDownstreamFailure`, and the breaker records a failure. So users closing tabs can trip the breaker on a perfectly healthy downstream, and a burst of cancellations (a client-side retry storm, a deploy of the caller) opens the circuit for everybody.

   A cancellation says nothing about the downstream's health, so it should be **neutral**:

   - In CLOSED: don't record a failure (and don't record a success).
   - In HALF_OPEN: release the probe slot (`probeInFlight = false`) *without* a verdict, so the next request becomes the probe.

   *How*

   Make "what counts as a failure" pluggable, as gobreaker does with `IsSuccessful`:

   ```go
   // in Config
   IsFailure func(err error) bool

   // default
   func defaultIsFailure(err error) bool {
       return err != nil && !errors.Is(err, context.Canceled)
   }
   ```

   `context.DeadlineExceeded` is ambiguous. If it is the *caller's* short deadline, it isn't the downstream's fault. If it is the HTTP client's own timeout (`url.Error` with `Timeout() == true`), the downstream is slow and it should count. Check `ctx.Err()` on the caller's context to tell them apart.

   *HTTP status codes*

   - **Other 4xx** (400, 401, 404, 422): the caller's mistake. Not a failure, which is what the code does today.
   - **429**: the downstream is telling us to slow down. That is a health signal. Count it as a failure, or better, handle it separately with backoff or a rate limiter, because tripping the breaker doesn't reduce load in a controlled way.
   - **503 / 502 / 504**: downstream is unavailable. Count as failure (the `>= 500` check already does). If `Retry-After` is present, honour it by using it as the open duration, so we don't probe before the service asked us to.
   - **408 / 5xx from a proxy**: depends on the topology; make it configurable rather than hard-coded.

7. You return `fmt.Errorf("%w: %v", ErrDownstreamFailure, err)`. What do callers lose with `%v` there? Also, in the 5xx branch you call `resp.Body.Close()` without reading the body. Does that have any effect on connection reuse?

   **Answer**

   *`%v` vs `%w`*

   `%v` only formats the underlying error as text, so the chain is cut. Callers lose `errors.Is(err, context.DeadlineExceeded)`, `errors.Is(err, context.Canceled)` and `errors.As(err, &netErr)` / `errors.As(err, &urlErr)`. That also kills the classification from question 6, because the breaker or the caller can no longer distinguish a timeout from connection refused from a cancellation. Since Go 1.20 a format can have several `%w`:

   ```go
   return fmt.Errorf("%w: %w", ErrDownstreamFailure, err)
   ```

   *Unread body and connection reuse*

   Yes, it matters. In HTTP/1.1, `net/http` can only return a connection to the keep-alive pool if the response body was read to EOF and closed. Closing without draining makes the transport tear the TCP connection down. During an outage, where almost every response is a 5xx, we would close and re-dial (TCP and TLS handshakes) on every failing call, adding load to a service that is already struggling. Drain a bounded amount first:

   ```go
   _, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4<<10))
   resp.Body.Close()
   ```

   *Smaller things in the same function*

   - The closure ignores its own `ctx` parameter and uses `req`, which was built with the outer `ctx`. It works but it is confusing, and the parameter shadows the outer name.
   - `response = resp` happens before the 5xx check. It is harmless because we return `nil, err`, but it reads as a bug.
   - The 5xx body is thrown away, so the caller can't see the downstream's error message. Returning a typed error that carries `StatusCode` (and maybe a truncated body) would help debugging.
   - On success the breaker records its verdict when headers arrive, so a failure while the caller reads the body isn't tracked.

## Production and testing

8. HALF_OPEN closes the circuit after a *single* successful probe. In a real system, what could go wrong with that, and how would you generalize it?

   **Answer**

   *What can go wrong*

   - **Flapping.** One lucky success closes the circuit, the full load hits a half-recovered service, it fails again, the circuit re-opens, and so on: CLOSED, OPEN, HALF_OPEN, CLOSED. Each cycle sends a burst of failing traffic.
   - **Thundering herd on a cold service.** Caches are cold, autoscaling or connection pools haven't caught up, and we go from 0% to 100% of traffic instantly.
   - **Weak evidence.** One success from a service that fails 50% of the time happens often.
   - **No backoff.** The next attempt is always exactly `OpenDuration` later. If many instances of the caller run the same logic they probe in lockstep.

   *How I'd generalize*

   - `HalfOpenMaxRequests`: allow up to N concurrent probes, not one.
   - `SuccessThreshold`: require N consecutive successful probes before closing. Any failure re-opens immediately.
   - Optionally ramp traffic gradually after closing (10%, 25%, 50%, 100%).
   - Exponential backoff of `OpenDuration` for consecutive re-opens (30s, 60s, 120s, capped), with jitter to desynchronize instances.

   The state would need a `successes` counter in HALF_OPEN, and the generation check from question 3 matters even more, because several requests are then in flight.

9. Your tests mutate `cb.now` by reaching into the struct and locking `cb.mu` manually. What's your opinion on that? What tests are *missing*? Specifically, how would you prove that only one probe is allowed under concurrent load, and what flag would you run the tests with?

   **Answer**

   *Opinion on the current approach*

   It works, because the test lives in the same package and the fake clock makes it deterministic and fast, with no `time.Sleep`. But it couples the tests to private fields, and the manual `cb.mu.Lock()` is unnecessary in a single-goroutine test. I would inject the clock properly, either as `Config.Clock` or a functional option, and use a small fake with `Advance(d)`:

   ```go
   type fakeClock struct{ t time.Time }
   func (f *fakeClock) Now() time.Time          { return f.t }
   func (f *fakeClock) Advance(d time.Duration) { f.t = f.t.Add(d) }
   ```

   *Missing tests*

   - **Concurrency**: only one probe admitted, run under `-race` (below).
   - **Stale result** (question 3): start a request in CLOSED, trip the breaker, move to HALF_OPEN, let the old request finish, assert the state is still HALF_OPEN and the probe slot is still taken.
   - **Panic in `fn`** (question 5): after a recovered panic the breaker is not stuck.
   - **Cancellation / failure classification** (question 6).
   - **Boundaries**: failure exactly at the window edge (`Before(cutoff)` means "exactly at the cutoff" still counts), and `now - openedAt == OpenDuration` exactly (allowed through).
   - **Success in CLOSED** doesn't reset failures. This documents the behaviour from question 1.
   - **`Config.Validate`** as a table-driven test, and `State().String()`.
   - **Client tests** with `httptest.Server`: 500 counts, 404 doesn't, a hung server hits the timeout, unknown service name returns an error, and the body is closed.

   *Proving one probe under concurrent load*

   ```go
   func TestOnlyOneProbeInHalfOpen(t *testing.T) {
       cb, clock := newOpenBreaker(t)          // helper: trips the breaker
       clock.Advance(31 * time.Second)

       var calls atomic.Int32
       release := make(chan struct{})

       const n = 100
       var wg sync.WaitGroup
       var rejected atomic.Int32
       start := make(chan struct{})

       for i := 0; i < n; i++ {
           wg.Add(1)
           go func() {
               defer wg.Done()
               <-start
               err := cb.Execute(context.Background(), func(context.Context) error {
                   calls.Add(1)
                   <-release // hold the probe open
                   return nil
               })
               if errors.Is(err, ErrCircuitHalfOpen) {
                   rejected.Add(1)
               }
           }()
       }

       close(start)
       // wait until the other 99 are rejected, then let the probe finish
       require.Eventually(t, func() bool { return rejected.Load() == n-1 }, time.Second, time.Millisecond)
       close(release)
       wg.Wait()

       if got := calls.Load(); got != 1 { t.Fatalf("expected 1 probe, got %d", got) }
       if cb.State() != StateClosed { t.Fatalf("expected CLOSED, got %s", cb.State()) }
   }
   ```

   *Flags*

   ```
   go test -race -count=100 -shuffle=on ./...
   ```

   `-race` catches unsynchronized access to `state` and `now`. `-count=100` re-runs the test to surface scheduling-dependent bugs, and `-shuffle=on` catches order dependence between tests. Concurrency bugs are probabilistic, so a single green run proves little.

10. If I put this in front of 200 downstream services at 10k RPS, what do you wish you had instrumented? How would an on-call engineer know it tripped, and why?

    **Answer**

    *Today, a trip is nearly invisible*

    The breaker has no name, no logging, no metrics and no callback. The only signal is callers seeing a bare `circuit breaker is open` error, with no idea which service, why, or until when. And `State()` reports OPEN lazily, so it can disagree with reality.

    *What I'd add*

    - **A name** on the breaker (`service-b`) and a label on everything below.
    - **An `OnStateChange(name, from, to, reason)` hook.** Used to log at WARN with structured fields: the failure count, the window, the last error that tripped it, and the open duration. Also used to emit a metric.
    - **Metrics**, per breaker:
      - `state` gauge (0/1/2).
      - `requests_total{outcome=success|failure|rejected_open|rejected_half_open|ignored}`.
      - `transitions_total{from,to}`.
      - time spent in OPEN and HALF_OPEN.
      - failure ratio in the current window.
      - downstream latency histogram.
    - **Richer errors.** `ErrCircuitOpen` should carry the service name and time until the next probe (`*OpenError{Service, RetryAfter}`), so callers and logs can say something useful.
    - **Tracing.** A span attribute such as `circuit_breaker.state` and `circuit_breaker.rejected=true` so a rejected call is visibly different from a slow call.

    *How on-call knows, and why*

    - **Alerts**: breaker OPEN for longer than N minutes, flapping (more than X transitions per 10 minutes), and rejection rate over a threshold. Each links to a dashboard.
    - **The why**: the transition log line holds the triggering error classes (timeouts vs connection refused vs 503). Combined with downstream latency and error-rate panels, on-call can see whether the downstream is down or the breaker is too sensitive (question 1).

    *Other things that matter at that scale*

    - **Config per service**, not one hard-coded config in `NewServiceClient`. A ratio plus minimum-volume threshold (question 1) is essential at 10k RPS.
    - **Per-instance breakers each have their own view.** With many replicas, some may be open while others aren't. Usually that is acceptable, and it should be understood.
    - **Metric cardinality**: label by service and state, never by path or error string.
    - **Lock contention** on the shared mutex under 10k RPS per breaker, and a profile before optimizing.
