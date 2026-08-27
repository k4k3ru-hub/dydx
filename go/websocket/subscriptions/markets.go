package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k4k3ru-hub/dydx/go/websocket/protocol"
)

// MarketsClient manages the all-markets subscription.
type MarketsClient struct {
	executor Executor
}

// MarketsParams configures an all-markets subscription.
type MarketsParams struct {
	Batched bool
}

// Key returns the transport key for the all-markets subscription.
//
// Version:
//   - 2026-08-27: Added.
func (MarketsParams) Key() (string, error) {
	return protocol.ChannelMarkets, nil
}

// NewMarketsClient creates an all-markets subscription client.
//
// Version:
//   - 2026-08-27: Added.
func NewMarketsClient(executor Executor) (*MarketsClient, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create markets client: executor=null")
	}
	return &MarketsClient{executor: executor}, nil
}

// Subscribe subscribes to all-market updates.
//
// Version:
//   - 2026-08-27: Added.
func (c *MarketsClient) Subscribe(ctx context.Context, params MarketsParams) error {
	return c.execute(ctx, protocol.MessageTypeSubscribe, params, true)
}

// Unsubscribe unsubscribes from all-market updates.
//
// Version:
//   - 2026-08-27: Added.
func (c *MarketsClient) Unsubscribe(ctx context.Context, params MarketsParams) error {
	return c.execute(ctx, protocol.MessageTypeUnsubscribe, params, false)
}

func (c *MarketsClient) execute(ctx context.Context, messageType string, params MarketsParams, subscribe bool) error {
	if c == nil {
		return fmt.Errorf("failed to %s markets: client=null", messageType)
	}
	if c.executor == nil {
		return fmt.Errorf("failed to %s markets: executor=null", messageType)
	}
	key, err := params.Key()
	if err != nil {
		return fmt.Errorf("failed to %s markets: %w", messageType, err)
	}

	request := protocol.Request{Type: messageType, Channel: protocol.ChannelMarkets}
	if subscribe {
		request.Batched = &params.Batched
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to %s markets: failed to encode request: %w", messageType, err)
	}
	if subscribe {
		err = c.executor.Subscribe(ctx, key, payload)
	} else {
		err = c.executor.Unsubscribe(ctx, key, payload)
	}
	if err != nil {
		return fmt.Errorf("failed to %s markets: %w", messageType, err)
	}
	return nil
}
