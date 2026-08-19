// Package historical_funding implements the dYdX historical-funding endpoint.
package historical_funding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/k4k3ru-hub/dydx/go/rest/endpoint"
	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type Client struct{ executor transport.Executor }

type Params struct {
	Market                    string
	EffectiveBeforeOrAt       string
	EffectiveBeforeOrAtHeight int64
	Limit                     int
}

type HistoricalFunding struct {
	Ticker            string `json:"ticker"`
	EffectiveAt       string `json:"effectiveAt"`
	EffectiveAtHeight string `json:"effectiveAtHeight"`
	Price             string `json:"price"`
	Rate              string `json:"rate"`
}

type Response struct {
	HistoricalFunding []HistoricalFunding `json:"historicalFunding"`
}

// NewClient creates a historical-funding operation client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create historical funding client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests settled historical funding rates for a perpetual market.
//
// Parameters:
//   - ctx: Context for the operation.
//   - params: Historical-funding request parameters.
//
// Returns:
//   - Historical funding rates returned by the Indexer.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, params Params) (*Response, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request historical funding: client=null")
	}
	market := strings.TrimSpace(params.Market)
	if market == "" {
		return nil, fmt.Errorf("failed to request historical funding: market=empty")
	}
	if params.EffectiveBeforeOrAtHeight < 0 {
		return nil, fmt.Errorf("failed to request historical funding: effective_before_or_at_height=out_of_range")
	}
	if params.Limit < 0 {
		return nil, fmt.Errorf("failed to request historical funding: limit=out_of_range")
	}

	query := url.Values{}
	if params.EffectiveBeforeOrAt != "" {
		query.Set("effectiveBeforeOrAt", params.EffectiveBeforeOrAt)
	}
	if params.EffectiveBeforeOrAtHeight != 0 {
		query.Set("effectiveBeforeOrAtHeight", strconv.FormatInt(params.EffectiveBeforeOrAtHeight, 10))
	}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}

	body, err := c.executor.Do(ctx, transport.Request{
		Method: http.MethodGet,
		Path:   endpoint.HistoricalFundingPathPrefix + url.PathEscape(market),
		Query:  query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request historical funding: %w", err)
	}
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request historical funding: failed to decode response body: %w", err)
	}
	return &result, nil
}
