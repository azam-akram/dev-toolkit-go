package handler

import (
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/logger"
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/model"
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/utils"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	logger.Get().Info("Http request received")

	switch req.URL.Path {
	case "/healthcheck":
		switch req.Method {
		case http.MethodGet:
			healthcheck(w)
		default:
			logger.Get().Info("End point /healthcheck doesn't support this method", "path", req.URL.Path, "method", http.StatusMethodNotAllowed)
			errResponse := map[string]string{
				"error": "HTTP method not supported",
			}
			writeFailedResponse(w, errResponse, http.StatusMethodNotAllowed)
		}
	case "/students":
		switch req.Method {
		case http.MethodPost:
			createStudent(w, req)
		default:
			logger.Get().Info("End point /students doesn't support this method", "path", req.URL.Path, "method", http.StatusMethodNotAllowed)
			errResponse := map[string]string{
				"error": "HTTP method not supported",
			}
			writeFailedResponse(w, errResponse, http.StatusMethodNotAllowed)
		}
	default:
		logger.Get().Info("Endpoint not supported", "path", req.URL.Path)
		errResponse := map[string]string{
			"error": "HTTP method not supported",
		}
		writeFailedResponse(w, errResponse, http.StatusNotFound)
	}
}

func healthcheck(w http.ResponseWriter) {
	resp := map[string]string{
		"Status": "OK",
	}

	writeSuccessResponse(w, http.StatusOK, resp)
}

func createStudent(w http.ResponseWriter, req *http.Request) {
	version := os.Getenv("APP_VERSION")
	if len(version) == 0 {
		version = "v1"
	}

	reqBodyStr, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Get().Error("Error reading HTTP request body", "error", err)
		errResponse := map[string]string{
			"error": "Error reading HTTP request body",
		}
		writeFailedResponse(w, errResponse, http.StatusBadRequest)
	}

	var sReq model.StudentRequest
	if utils.StringToStruct(string(reqBodyStr), &sReq); err != nil {
		logger.Get().Error("Error unmarshalling request body", "error", err)
	}

	logger.Get().Info("Unmarshalled Student Object", "sReq", sReq)

	var sResp model.StudentResponse
	sResp.Created = time.Now()
	sResp.CreatedBy = sReq.CreatedBy
	sResp.HostName, _ = os.Hostname()
	sResp.Students = sReq.Students
	sResp.Version = version

	writeSuccessResponse(w, http.StatusCreated, sResp)
}

func writeSuccessResponse(w http.ResponseWriter, code int, resp any) {
	logger.Get().Info("Returning response", "resp", resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Get().Error("Failed to encode success response", "error", err)
		errResponse := map[string]string{
			"error": "Failed to encode success response",
		}
		writeFailedResponse(w, errResponse, http.StatusBadRequest)
	}
}

func writeFailedResponse(w http.ResponseWriter, resp any, code int) {
	logger.Get().Info("Returning failed response", "resp", resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Get().Error("Failed to encode success response", "error", err)
	}
}
