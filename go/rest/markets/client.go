// Package markets composes dYdX Indexer market-data operations.
package markets

import (
	"fmt"

	"github.com/k4k3ru-hub/dydx/go/rest/markets/order_book"
	"github.com/k4k3ru-hub/dydx/go/rest/markets/perpetual_markets"
	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type Client struct {
	perpetualMarkets *perpetual_markets.Client
	orderBook        *order_book.Client
}

// NewClient creates a Markets API client.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create markets client: executor=null")
	}
	perpetualMarketsClient, err := perpetual_markets.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create markets client: %w", err)
	}
	orderBookClient, err := order_book.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create markets client: %w", err)
	}
	return &Client{perpetualMarkets: perpetualMarketsClient, orderBook: orderBookClient}, nil
}

// PerpetualMarkets returns the perpetual-markets operation client.
func (c *Client) PerpetualMarkets() *perpetual_markets.Client {
	if c == nil {
		return nil
	}
	return c.perpetualMarkets
}

// OrderBook returns the perpetual order-book operation client.
func (c *Client) OrderBook() *order_book.Client {
	if c == nil {
		return nil
	}
	return c.orderBook
}
