package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k4k3ru-hub/dydx/go/websocket/protocol"
)

// TradesClient manages public trade subscriptions.
type TradesClient struct {
	executor Executor
}

// TradesParams configures a public trade subscription.
type TradesParams struct {
	Market  string
	Batched bool
}

// Key returns the transport key for the public trade subscription.
//
// Version:
//   - 2026-08-17: Added.
func (p TradesParams) Key() (string, error) {
	market := strings.TrimSpace(p.Market)
	if market == "" {
		return "", fmt.Errorf("failed to build trades subscription key: market=empty")
	}
	return protocol.ChannelTrades + ":" + market, nil
}

// NewTradesClient creates a public trade subscription client.
//
// Version:
//   - 2026-08-17: Added.
func NewTradesClient(executor Executor) (*TradesClient, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create trades client: executor=null")
	}
	return &TradesClient{executor: executor}, nil
}

// Subscribe subscribes to public trade updates.
//
// Version:
//   - 2026-08-17: Added.
func (c *TradesClient) Subscribe(ctx context.Context, params TradesParams) error {
	return c.execute(ctx, protocol.MessageTypeSubscribe, params, true)
}

// Unsubscribe unsubscribes from public trade updates.
//
// Version:
//   - 2026-08-17: Added.
func (c *TradesClient) Unsubscribe(ctx context.Context, params TradesParams) error {
	return c.execute(ctx, protocol.MessageTypeUnsubscribe, params, false)
}

func (c *TradesClient) execute(ctx context.Context, messageType string, params TradesParams, subscribe bool) error {
	if c == nil {
		return fmt.Errorf("failed to %s trades: client=null", messageType)
	}
	if c.executor == nil {
		return fmt.Errorf("failed to %s trades: executor=null", messageType)
	}
	key, err := params.Key()
	if err != nil {
		return fmt.Errorf("failed to %s trades: %w", messageType, err)
	}

	request := protocol.Request{
		Type:    messageType,
		Channel: protocol.ChannelTrades,
		ID:      strings.TrimSpace(params.Market),
	}
	if subscribe {
		request.Batched = &params.Batched
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to %s trades: failed to encode request: %w", messageType, err)
	}

	if subscribe {
		err = c.executor.Subscribe(ctx, key, payload)
	} else {
		err = c.executor.Unsubscribe(ctx, key, payload)
	}
	if err != nil {
		return fmt.Errorf("failed to %s trades: %w", messageType, err)
	}
	return nil
}
