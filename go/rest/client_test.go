package rest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(request *http.Request) (*http.Response, error) { return f(request) }

func TestDoBuildsURL(t *testing.T) {
	client, err := NewClient(&ClientOption{
		BaseURL: "https://example.test/",
		HTTPClient: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.String() != "https://example.test/v4/test?market=BTC-USD" {
				t.Fatalf("unexpected URL: %s", request.URL)
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Do(context.Background(), transport.Request{Method: http.MethodGet, Path: "/v4/test", Query: map[string][]string{"market": {"BTC-USD"}}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoReturnsTypedResponseError(t *testing.T) {
	client, _ := NewClient(&ClientOption{HTTPClient: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusTooManyRequests, Body: io.NopCloser(strings.NewReader(`rate limited`))}, nil
	})})
	_, err := client.Do(context.Background(), transport.Request{Method: http.MethodGet, Path: "/v4/test"})
	var responseError *ResponseError
	if !errors.As(err, &responseError) || responseError.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unexpected error: %v", err)
	}
}
