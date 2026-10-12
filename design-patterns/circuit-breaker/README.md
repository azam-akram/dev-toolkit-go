# Circuit Breaker design pattern

```
                 failure threshold reached
        ┌─────────────────────────────────────┐
        │                                     ▼
    ┌────────┐                           ┌────────┐
    │ CLOSED │                           │  OPEN  │
    └────────┘                           └────────┘
        ▲                                     │
        │ success                             │ 10 min elapsed
        │                                     ▼
        │                                ┌───────────┐
        └────────────────────────────────│ HALF_OPEN │
                                         └───────────┘
                                               │
                                  ┌────────────┴────────────┐
                                  │                         │
                               success                   failure
                                  │                         │
                                  ▼                         ▼
                              CLOSED                      OPEN
```

## One important production improvement

In a real system, I would usually distinguish what counts as a failure from the circuit breaker itself.
For example:
type FailureClassifier func(resp *http.Response, err error) bool
Then you can explicitly configure:

```
Network error     → failure
Timeout           → failure
HTTP 500          → failure
HTTP 502          → failure
HTTP 503          → failure
HTTP 504          → failure

HTTP 400          → not a failure
HTTP 401          → not a failure
HTTP 403          → not a failure
HTTP 404          → not a failure
```

## Concurrency consideration

The probeInFlight field is particularly important.
Imagine the circuit is OPEN and the recovery period expires.
Without protection, 100 goroutines could simultaneously see:

```
OPEN
  ↓
recovery period elapsed
  ↓
HALF_OPEN
```

and all 100 could send requests.
That defeats the purpose of a cautious HALF_OPEN probe.
With probeInFlight:

```
HALF_OPEN
                 │
       ┌─────────┼─────────┐
       │         │         │
    request 1 request 2 request 3
       │         │         │
       ▼         ▼         ▼
      ALLOW     BLOCK     BLOCK
       │
     probe
       │
   ┌───┴────┐
 success   failure
   │          │
   ▼          ▼
CLOSED       OPEN
```

Only the first request gets through.

## Final state machine

The whole implementation can be understood as this:


