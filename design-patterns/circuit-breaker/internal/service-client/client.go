package httpclient

import (
	"context"
	circuitbreaker "dev-toolkit-go/design-patterns/circuit-breaker/internal/circuit-breaker"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrDownstreamFailure = errors.New("downstream service failure")

type ServiceClient struct {
	httpClient *http.Client

	breakers map[string]*circuitbreaker.CircuitBreaker

	baseURLs map[string]string
}

const (
	ServiceB = "service-b"
	ServiceC = "service-c"
)

func NewServiceClient(httpClient *http.Client) (*ServiceClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	breakerConfig := circuitbreaker.Config{
		FailureThreshold: 5,
		FailureWindow:    10 * time.Minute,
		OpenDuration:     30 * time.Second,
	}

	serviceBBreaker, err := circuitbreaker.NewCircuitBreaker(breakerConfig)
	if err != nil {
		return nil, err
	}

	serviceCBreaker, err := circuitbreaker.NewCircuitBreaker(breakerConfig)
	if err != nil {
		return nil, err
	}

	return &ServiceClient{
		httpClient: httpClient,

		breakers: map[string]*circuitbreaker.CircuitBreaker{
			ServiceB: serviceBBreaker,
			ServiceC: serviceCBreaker,
		},

		baseURLs: map[string]string{
			ServiceB: "https://service-b.internal",
			ServiceC: "https://service-c.internal",
		},
	}, nil
}

func (c *ServiceClient) Do(
	ctx context.Context,
	service string,
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {

	breaker, ok := c.breakers[service]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s", service)
	}

	baseURL, ok := c.baseURLs[service]
	if !ok {
		return nil, fmt.Errorf("no base URL configured for service: %s", service)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		baseURL+path,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	var response *http.Response

	err = breaker.Execute(ctx, func(ctx context.Context) error {
		resp, err := c.httpClient.Do(req)

		if err != nil {
			// Network failure, timeout, connection failure, etc.
			return fmt.Errorf("%w: %v", ErrDownstreamFailure, err)
		}

		response = resp

		// 5xx means the downstream service failed.
		if resp.StatusCode >= 500 {
			// We should close the response body because we're
			// not returning the response to the caller.
			resp.Body.Close()

			return fmt.Errorf(
				"%w: service=%s status=%d",
				ErrDownstreamFailure,
				service,
				resp.StatusCode,
			)
		}

		// 4xx is generally an application/client error.
		// It should NOT normally trip the circuit breaker.
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ServiceClient) CallServiceB(
	ctx context.Context,
) (*http.Response, error) {

	return c.Do(
		ctx,
		ServiceB,
		http.MethodGet,
		"/api/resource",
		nil,
	)
}

func (c *ServiceClient) CallServiceC(
	ctx context.Context,
) (*http.Response, error) {

	return c.Do(
		ctx,
		ServiceC,
		http.MethodGet,
		"/api/resource",
		nil,
	)
}
