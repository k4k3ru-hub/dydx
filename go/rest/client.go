// Package rest implements the dYdX Indexer REST client.
package rest

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/k4k3ru-hub/dydx/go/rest/endpoint"
	"github.com/k4k3ru-hub/dydx/go/rest/markets"
	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOption struct {
	BaseURL        string
	ConnectTimeout time.Duration
	HTTPClient     HTTPClient
}

type Client struct {
	baseURL       string
	httpClient    HTTPClient
	marketsClient *markets.Client
}

// DefaultClientOption returns the default REST client options.
func DefaultClientOption() *ClientOption {
	return &ClientOption{BaseURL: endpoint.DefaultBaseURL, ConnectTimeout: 3 * time.Second}
}

// NewClient creates a REST client and composes its API clients.
func NewClient(option *ClientOption) (*Client, error) {
	if option == nil {
		option = DefaultClientOption()
	}
	baseURL := strings.TrimRight(option.BaseURL, "/")
	if baseURL == "" {
		baseURL = endpoint.DefaultBaseURL
	}
	httpClient := option.HTTPClient
	if httpClient == nil {
		connectTimeout := option.ConnectTimeout
		if connectTimeout == 0 {
			connectTimeout = DefaultClientOption().ConnectTimeout
		}
		if connectTimeout < 0 {
			return nil, fmt.Errorf("failed to create rest client: connect_timeout=out_of_range")
		}
		httpClient = &http.Client{Transport: &http.Transport{DialContext: (&net.Dialer{Timeout: connectTimeout}).DialContext}}
	}
	client := &Client{baseURL: baseURL, httpClient: httpClient}
	marketsClient, err := markets.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest client: %w", err)
	}
	client.marketsClient = marketsClient
	return client, nil
}

// Markets returns the composed Markets API client.
func (c *Client) Markets() *markets.Client {
	if c == nil {
		return nil
	}
	return c.marketsClient
}

// Do executes a REST request.
func (c *Client) Do(ctx context.Context, request transport.Request) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to execute rest request: client=null")
	}
	if c.httpClient == nil {
		return nil, fmt.Errorf("failed to execute rest request: http_client=null")
	}
	if request.Method == "" {
		return nil, fmt.Errorf("failed to execute rest request: method=empty")
	}
	if request.Path == "" {
		return nil, fmt.Errorf("failed to execute rest request: path=empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, c.baseURL+"/"+strings.TrimLeft(request.Path, "/"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: failed to create HTTP request: %w", err)
	}
	httpRequest.URL.RawQuery = request.Query.Encode()
	httpRequest.Header = request.Header.Clone()
	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: %w", err)
	}
	if response == nil {
		return nil, fmt.Errorf("failed to execute rest request: response=null")
	}
	if response.Body == nil {
		return nil, fmt.Errorf("failed to execute rest request: response_body=null")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: failed to read response body: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, &ResponseError{StatusCode: response.StatusCode}
	}
	return body, nil
}

type ResponseError struct{ StatusCode int }

// Error returns the REST response error message.
func (e *ResponseError) Error() string {
	if e == nil {
		return "failed to execute rest request: response_error=null"
	}
	return fmt.Sprintf("failed to execute rest request: unexpected HTTP status: status_code=%d", e.StatusCode)
}
