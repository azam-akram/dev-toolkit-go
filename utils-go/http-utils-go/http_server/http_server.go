package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	LISTENING_PORT = 8989
	MAX_COUNTER    = 4
)

type Profile struct {
	Name    string   `json:"name"`
	Hobbies []string `json:"hobbies"`
}

var (
	counter = 0
	handler Handler
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type HTTPServer struct{}

func NewHTTPServer() Handler {
	if handler == nil {
		handler = &HTTPServer{}
	}
	return handler
}

func (h HTTPServer) Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("HTTPServer::Handler counter = %d\n", counter)

	if counter < MAX_COUNTER {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		counter++
		return
	}

	fmt.Printf("HTTPServer: recovered from problem after %d failed attempts\n", counter)

	profile := &Profile{
		Name:    "User",
		Hobbies: []string{"Sports", "Walk"},
	}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	fmt.Printf("Server is listening at port %d\n", LISTENING_PORT)

	h := NewHTTPServer()
	http.HandleFunc("/", h.Handle)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", LISTENING_PORT), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
