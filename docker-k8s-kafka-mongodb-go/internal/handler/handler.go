package handler

import (
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/model"
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Request %v", req)

	switch req.URL.Path {
	case "/healthcheck":
		switch req.Method {
		case http.MethodGet:
			healthcheck(w)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/students":
		switch req.Method {
		case http.MethodPost:
			createStudent(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	default:
		fmt.Printf("Endpoint not supported: %s", req.URL.Path)
	}
}

func healthcheck(w http.ResponseWriter) {
	resp := map[string]string{
		"Status": "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("Failed to encode healthcheck response: %v\n", err)
	}
}

func createStudent(w http.ResponseWriter, req *http.Request) {
	reqBodyStr, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("error reading HTTP request body: %v\n", err)
	}

	var sReq model.StudentRequest
	if utils.StringToStruct(string(reqBodyStr), &sReq); err != nil {
		fmt.Printf("error unmarshalling request body: %v\n", err)
	}

	fmt.Printf("Student Request: %v\n", sReq)

	var sResp model.StudentResponse
	sResp.Created = time.Now()
	sResp.CreatedBy = sReq.CreatedBy
	sResp.HostName, _ = os.Hostname()
	sResp.Students = sReq.Students
	sResp.Version = "v1" // TODO

	fmt.Printf("Student Response: %v\n", sResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(sResp)
}
