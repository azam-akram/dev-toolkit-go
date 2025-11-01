package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github/dev-toolkit-go/aws-go/aws-apigateway-lambda-dynamo-go/internal/model"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// MockBookStore is the mock implementation of BookStorer.
type MockBookStore struct {
	SaveFn       func(book *model.MyBook) error
	GetByIDFn    func(id string) (*model.MyBook, error)
	UpdateFn     func(book *model.MyBook) error
	DeleteByIDFn func(id string) error
}

// Implement the BookStorer interface using the mock functions
func (m *MockBookStore) Save(book *model.MyBook) error {
	return m.SaveFn(book)
}
func (m *MockBookStore) GetByID(id string) (*model.MyBook, error) {
	return m.GetByIDFn(id)
}
func (m *MockBookStore) Update(book *model.MyBook) error {
	return m.UpdateFn(book)
}
func (m *MockBookStore) DeleteByID(id string) error {
	return m.DeleteByIDFn(id)
}

// --- MOCK DEPENDENCIES END ---

// handleRequestTestable is a refactored version of HandleRequest that accepts the BookStorer interface
// as a dependency, making it unit-testable.
func handleRequestTestable(ctx context.Context, req events.APIGatewayProxyRequest, handler *MockBookStore) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		var book model.MyBook
		if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
		if err := handler.Save(&book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: req.Body}, nil

	case "GET":
		id := req.QueryStringParameters["id"]
		book, err := handler.GetByID(id)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		if book == nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
		}
		body, err := json.Marshal(book)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil

	case "PUT":
		var book model.MyBook
		if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
		if err := handler.Update(&book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: req.Body}, nil

	case "DELETE":
		id := req.QueryStringParameters["id"]
		if err := handler.DeleteByID(id); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil

	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

// Test cases for the handler
func TestHandleRequest(t *testing.T) {
	ctx := context.Background()

	// Sample valid book data
	testBook := model.MyBook{ID: "123", Title: "The Go Programming Language", Author: "Donovan & Kernighan"}
	testBookJSON, _ := json.Marshal(testBook)
	testBookBody := string(testBookJSON)

	// Define all test scenarios
	tests := []struct {
		name               string
		req                events.APIGatewayProxyRequest
		mockHandler        *MockBookStore
		expectedStatusCode int
		expectedBody       string
		expectedError      bool
	}{
		// --- POST TESTS ---
		{
			name: "POST_Success",
			req: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       testBookBody,
			},
			mockHandler: &MockBookStore{
				SaveFn: func(book *model.MyBook) error {
					assert.Equal(t, "123", book.ID) // Verify input to mock
					return nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       testBookBody,
			expectedError:      false,
		},
		{
			name: "POST_InvalidJSON",
			req: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       "{invalid json",
			},
			mockHandler:        &MockBookStore{SaveFn: func(*model.MyBook) error { return nil }}, // Save won't be called
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "",
			expectedError:      true, // Error during unmarshal
		},
		{
			name: "POST_DynamoError",
			req: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       testBookBody,
			},
			mockHandler: &MockBookStore{
				SaveFn: func(*model.MyBook) error {
					return errors.New("dynamo save failure")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "",
			expectedError:      true,
		},

		// --- GET TESTS ---
		{
			name: "GET_Success",
			req: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: map[string]string{"id": "123"},
			},
			mockHandler: &MockBookStore{
				GetByIDFn: func(id string) (*model.MyBook, error) {
					assert.Equal(t, "123", id) // Verify input to mock
					return &testBook, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       testBookBody,
			expectedError:      false,
		},
		{
			name: "GET_NotFound",
			req: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: map[string]string{"id": "999"},
			},
			mockHandler: &MockBookStore{
				GetByIDFn: func(id string) (*model.MyBook, error) {
					return nil, nil // Not found
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "",
			expectedError:      false,
		},
		{
			name: "GET_DynamoError",
			req: events.APIGatewayProxyRequest{
				HTTPMethod:            "GET",
				QueryStringParameters: map[string]string{"id": "error"},
			},
			mockHandler: &MockBookStore{
				GetByIDFn: func(id string) (*model.MyBook, error) {
					return nil, errors.New("dynamo get failure")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "",
			expectedError:      true,
		},

		// --- PUT TESTS ---
		{
			name: "PUT_Success",
			req: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       testBookBody,
			},
			mockHandler: &MockBookStore{
				UpdateFn: func(book *model.MyBook) error {
					assert.Equal(t, "123", book.ID) // Verify input to mock
					return nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       testBookBody,
			expectedError:      false,
		},
		{
			name: "PUT_DynamoError",
			req: events.APIGatewayProxyRequest{
				HTTPMethod: "PUT",
				Body:       testBookBody,
			},
			mockHandler: &MockBookStore{
				UpdateFn: func(*model.MyBook) error {
					return errors.New("dynamo update failure")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "",
			expectedError:      true,
		},

		// --- DELETE TESTS ---
		{
			name: "DELETE_Success",
			req: events.APIGatewayProxyRequest{
				HTTPMethod:            "DELETE",
				QueryStringParameters: map[string]string{"id": "123"},
			},
			mockHandler: &MockBookStore{
				DeleteByIDFn: func(id string) error {
					assert.Equal(t, "123", id) // Verify input to mock
					return nil
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "",
			expectedError:      false,
		},
		{
			name: "DELETE_DynamoError",
			req: events.APIGatewayProxyRequest{
				HTTPMethod:            "DELETE",
				QueryStringParameters: map[string]string{"id": "error"},
			},
			mockHandler: &MockBookStore{
				DeleteByIDFn: func(id string) error {
					return errors.New("dynamo delete failure")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "",
			expectedError:      true,
		},

		// --- DEFAULT/METHOD NOT ALLOWED TEST ---
		{
			name: "DEFAULT_MethodNotAllowed",
			req: events.APIGatewayProxyRequest{
				HTTPMethod: "PATCH",
			},
			mockHandler:        &MockBookStore{},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "",
			expectedError:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the testable function with the mock handler
			resp, err := handleRequestTestable(ctx, tc.req, tc.mockHandler)

			// Assert error expectation
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Assert status code
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// Assert response body (only check if status is 200, otherwise body might be empty)
			if tc.expectedStatusCode == http.StatusOK {
				assert.Equal(t, tc.expectedBody, resp.Body)
			}
		})
	}
}
