package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	RETRY_COUNT                 = 3
	RETRY_MIN_WAIT_TIME_SECONDS = 5
	RETRY_MAX_WAIT_TIME_SECONDS = 15
)

func main() {
	client := resty.New().
		SetRetryCount(RETRY_COUNT).
		SetRetryWaitTime(RETRY_MIN_WAIT_TIME_SECONDS * time.Second).
		SetRetryMaxWaitTime(RETRY_MAX_WAIT_TIME_SECONDS * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusRequestTimeout ||
					r.StatusCode() >= http.StatusInternalServerError
			},
		)

	resp, err := client.R().Get("http://localhost:8989/")
	if err != nil {
		fmt.Println("Error occurred while making request:", err)
		return
	}

	fmt.Println("Status:", resp.Status())
	fmt.Println("Response Body:", resp.String())
}
