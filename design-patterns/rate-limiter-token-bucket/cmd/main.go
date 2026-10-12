package main

import (
	tokenbucket "brain-dump/design-patterns/rate-limiter-token-bucket/limiter"
	"log"
	"net/http"
)

func main() {
	limiter, err := tokenbucket.New(
		tokenbucket.Config{
			Capacity:   10,
			RefillRate: 2,
		},
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Write([]byte("Hello!\n"))
	})

	mux.HandleFunc("GET /health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
	})

	handler := tokenbucket.Middleware(limiter)(mux)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
