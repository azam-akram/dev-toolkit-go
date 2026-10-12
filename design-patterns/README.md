# Design Pattern

For a backend/Go interview, I'd distinguish between classic GoF design patterns and the patterns interviewers commonly expect you to recognize in distributed systems / production backend code.

```
1. Resilience patterns — very commonly asked
Pattern	Problem it solves	Typical example
Circuit Breaker	Downstream is failing repeatedly	Stop calling Service B
Retry	Temporary/transient failures	Retry a timeout
Timeout	Dependency takes too long	Don't wait forever for Service B
Bulkhead	One dependency consumes all resources	Separate worker pools per service
Rate Limiter	Too much incoming traffic	Max 100 req/sec
Backpressure	Producer is faster than consumer	Queue fills → slow/reject producer
Fallback	Dependency unavailable	Return cached/default response
Exponential Backoff	Avoid retry storms	1s → 2s → 4s → 8s
Idempotency	Safe retries	Prevent duplicate payment/order
```

These are particularly valuable because they often work together:

```
                  HTTP Request
                       │
                       ▼
                ┌─────────────┐
                │ Rate Limiter│
                └──────┬──────┘
                       │
                       ▼
                 ┌───────────┐
                 │  Timeout  │
                 └─────┬─────┘
                       │
                       ▼
                 ┌───────────┐
                 │  Retry    │
                 └─────┬─────┘
                       │
                       ▼
               ┌──────────────┐
               │Circuit Breaker│
               └──────┬───────┘
                      │
                      ▼
                  Service B
```

An interviewer might actually ask:
"Service B is intermittently failing. How would you make your HTTP client resilient?"

That's where you discuss timeout + retry + exponential backoff + circuit breaker + bulkhead, rather than treating each pattern independently.

2. Concurrency patterns — especially important in Go
These are probably more important for a Go interview than Abstract Factory.
Worker Pool

```
             Jobs
              │
        ┌─────┴─────┐
        ▼           ▼
     Worker 1    Worker 2
        │           │
        ▼           ▼
     Worker 3    Worker 4
```

Useful for:
- controlled concurrency
- background processing
- consuming queues
- preventing thousands of goroutines from hitting a dependency
Typical interview question:
"You receive 100,000 jobs. How do you process them without creating 100,000 concurrent requests?"

Worker pool.

## Producer / Consumer

```
Producer
   │
   ▼
┌─────────┐
│ Channel │
└────┬────┘
     │
     ├──► Consumer 1
     ├──► Consumer 2
     └──► Consumer 3
```

## Semaphore / Concurrency Limiter

Suppose Service B can safely handle only 20 concurrent requests.

```
sem := make(chan struct{}, 20)

sem <- struct{}{}        // acquire
defer func() {
    <-sem                 // release
}()

callServiceB()
```

## Fan-out / Fan-in
Very common system-design coding problem.

```
                  ┌── Service A
                  │
Request ──────────┼── Service B
                  │
                  └── Service C
                         │
                         ▼
                     Aggregate
```

In Go:

```
go callA()
go callB()
go callC()
```

Then collect the results.
Interviewers may ask:
"Three services can be called independently. How can you reduce overall latency?"

Fan-out/fan-in.

## Caching patterns

Very common in backend interviews.
Cache-Aside
Probably the most important one to know.

```
          Request
             │
             ▼
         ┌───────┐
         │ Cache │
         └───┬───┘
             │
       hit ──┤── miss
             │     │
             │     ▼
             │   Database
             │     │
             │     ▼
             │   Cache
             │
             ▼
          Response
```

Typical question:
"How would you reduce database load?"

Cache-aside.

## Write-Through
```
Application
     │
     ▼
   Cache
     │
     ▼
 Database
```

Application writes to cache, and cache writes to DB.



