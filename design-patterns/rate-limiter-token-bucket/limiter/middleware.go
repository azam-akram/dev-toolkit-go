package tokenbucket

import "net/http"

// Middleware rejects requests when no token is available.
func Middleware(limiter Limiter) func(http.Handler) http.Handler {
	if limiter == nil {
		panic("tokenbucket: nil limiter")
	}

	return func(next http.Handler) http.Handler {
		if next == nil {
			panic("tokenbucket: nil next handler")
		}

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if !limiter.Allow() {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("Cache-Control", "no-store")
				w.Header().Set("Retry-After", "1")
				http.Error(
					w,
					"rate limit exceeded",
					http.StatusTooManyRequests,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
