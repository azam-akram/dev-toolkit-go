package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	RETRY_COUNT                 = 3
	RETRY_MIN_WAIT_TIME_SECONDS = 5
	RETRY_MAX_WAIT_TIME_SECONDS = 15
)

// retryTransport wraps an http.RoundTripper and retries requests that fail
// or return a retryable status code, using exponential backoff with jitter.
type retryTransport struct {
	next        http.RoundTripper
	retryCount  int
	minWaitTime time.Duration
	maxWaitTime time.Duration
	shouldRetry func(resp *http.Response, err error) bool
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the body so it can be re-sent on retry.
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rt.retryCount; attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(newByteReader(bodyBytes))
		}

		resp, err = rt.next.RoundTrip(req)

		if !rt.shouldRetry(resp, err) || attempt == rt.retryCount {
			return resp, err
		}

		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(rt.backoff(attempt))
	}

	return resp, err
}

func (rt *retryTransport) backoff(attempt int) time.Duration {
	wait := rt.minWaitTime * time.Duration(math.Pow(2, float64(attempt)))
	if wait > rt.maxWaitTime {
		wait = rt.maxWaitTime
	}
	// Add jitter: random value in [0, wait/2)
	jitter := time.Duration(rand.Int63n(int64(wait) / 2))
	return wait + jitter
}

func newByteReader(b []byte) io.Reader {
	return &byteReader{data: b}
}

type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func newRetryableClient() *http.Client {
	return &http.Client{
		Transport: &retryTransport{
			next:        http.DefaultTransport,
			retryCount:  RETRY_COUNT,
			minWaitTime: RETRY_MIN_WAIT_TIME_SECONDS * time.Second,
			maxWaitTime: RETRY_MAX_WAIT_TIME_SECONDS * time.Second,
			shouldRetry: func(resp *http.Response, err error) bool {
				if err != nil {
					return true
				}
				return resp.StatusCode == http.StatusRequestTimeout ||
					resp.StatusCode >= http.StatusInternalServerError
			},
		},
	}
}

func main() {
	client := newRetryableClient()

	resp, err := client.Get("http://localhost:8989/")
	if err != nil {
		fmt.Println("Error occurred while making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error occurred while reading response body:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
