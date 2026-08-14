// Package perpetual_markets implements the dYdX perpetual-markets endpoint.
package perpetual_markets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k4k3ru-hub/dydx/go/rest/endpoint"
	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type Client struct{ executor transport.Executor }

type Params struct {
	Market string
	Limit  int
}

type PerpetualMarkets struct {
	Markets map[string]PerpetualMarket `json:"markets"`
}

type PerpetualMarket struct {
	ClobPairID                string `json:"clobPairId"`
	Ticker                    string `json:"ticker"`
	Status                    string `json:"status"`
	OraclePrice               string `json:"oraclePrice"`
	PriceChange24H            string `json:"priceChange24H"`
	Volume24H                 string `json:"volume24H"`
	Trades24H                 int    `json:"trades24H"`
	NextFundingRate           string `json:"nextFundingRate"`
	InitialMarginFraction     string `json:"initialMarginFraction"`
	MaintenanceMarginFraction string `json:"maintenanceMarginFraction"`
	OpenInterest              string `json:"openInterest"`
	AtomicResolution          int    `json:"atomicResolution"`
	QuantumConversionExponent int    `json:"quantumConversionExponent"`
	TickSize                  string `json:"tickSize"`
	StepSize                  string `json:"stepSize"`
	StepBaseQuantums          int64  `json:"stepBaseQuantums"`
	SubticksPerTick           int64  `json:"subticksPerTick"`
	MarketType                string `json:"marketType"`
	OpenInterestLowerCap      string `json:"openInterestLowerCap"`
	OpenInterestUpperCap      string `json:"openInterestUpperCap"`
	BaseOpenInterest          string `json:"baseOpenInterest"`
	DefaultFundingRate1H      string `json:"defaultFundingRate1H"`
}

// NewClient creates a perpetual-markets operation client.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create perpetual markets client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests perpetual-market metadata.
func (c *Client) Send(ctx context.Context, params Params) (*PerpetualMarkets, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request perpetual markets: client=null")
	}
	query := url.Values{}
	if params.Market != "" {
		query.Set("ticker", params.Market)
	}
	if params.Limit < 0 {
		return nil, fmt.Errorf("failed to request perpetual markets: limit=out_of_range")
	}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.PerpetualMarketsPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request perpetual markets: %w", err)
	}
	var result PerpetualMarkets
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request perpetual markets: failed to decode response body: %w", err)
	}
	return &result, nil
}
