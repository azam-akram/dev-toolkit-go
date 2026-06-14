package main

import (
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/handler"
	"fmt"

	"net/http"
)

type Student struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func main() {
	http.HandleFunc("/", handler.HandleRequest)

	fmt.Println("HTTP server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start http server: %v\n", err)
	}

	fmt.Println("HTTP server is listening on :8080")

}
