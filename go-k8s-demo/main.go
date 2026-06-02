package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message  string    `json:"message"`
	Hostname string    `json:"hostname"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	resp := Response{
		Message:  "Hello from Kubernetes!",
		Hostname: hostname,
		Version:  os.Getenv("APP_VERSION"),
		Time:     time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/healthcheck", healthcheck)

	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Listen error: %v", err)
	}
}
