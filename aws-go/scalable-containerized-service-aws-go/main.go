package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Response struct to structure the API response
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// handler function for the API endpoint
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message:   "Welcome to the Go App on ECS!",
		Timestamp: time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// healthCheck handler to verify app health
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthCheck) // health check endpoint
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
