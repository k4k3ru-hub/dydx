// Package transport defines the dYdX Indexer REST execution contract.
package transport

import (
	"context"
	"net/http"
	"net/url"
)

type Request struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
}

type Executor interface {
	Do(ctx context.Context, request Request) ([]byte, error)
}
