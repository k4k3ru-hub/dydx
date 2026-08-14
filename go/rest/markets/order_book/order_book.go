// Package order_book implements the dYdX perpetual order-book endpoint.
package order_book

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/k4k3ru-hub/dydx/go/rest/endpoint"
	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type Client struct{ executor transport.Executor }

type Params struct{ Market string }

type PriceLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

type OrderBook struct {
	Bids []PriceLevel `json:"bids"`
	Asks []PriceLevel `json:"asks"`
}

// NewClient creates a perpetual order-book operation client.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create order book client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests a perpetual-market order-book snapshot.
func (c *Client) Send(ctx context.Context, params Params) (*OrderBook, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request order book: client=null")
	}
	market := strings.TrimSpace(params.Market)
	if market == "" {
		return nil, fmt.Errorf("failed to request order book: market=empty")
	}
	path := endpoint.OrderBookPathPrefix + url.PathEscape(market)
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: path})
	if err != nil {
		return nil, fmt.Errorf("failed to request order book: %w", err)
	}
	var result OrderBook
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request order book: failed to decode response body: %w", err)
	}
	return &result, nil
}
