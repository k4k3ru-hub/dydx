package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k4k3ru-hub/dydx/go/websocket/protocol"
)

type OrderBookClient struct {
	executor Executor
}

type OrderBookParams struct {
	Market  string
	Batched bool
}

func (p OrderBookParams) Key() (string, error) {
	market := strings.TrimSpace(p.Market)
	if market == "" {
		return "", fmt.Errorf("failed to build order book subscription key: market=empty")
	}
	return protocol.ChannelOrderBook + ":" + market, nil
}

func NewOrderBookClient(executor Executor) (*OrderBookClient, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create order book client: executor=null")
	}
	return &OrderBookClient{executor: executor}, nil
}

func (c *OrderBookClient) Subscribe(ctx context.Context, params OrderBookParams) error {
	return c.execute(ctx, protocol.MessageTypeSubscribe, params, true)
}

func (c *OrderBookClient) Unsubscribe(ctx context.Context, params OrderBookParams) error {
	return c.execute(ctx, protocol.MessageTypeUnsubscribe, params, false)
}

func (c *OrderBookClient) execute(ctx context.Context, messageType string, params OrderBookParams, subscribe bool) error {
	if c == nil {
		return fmt.Errorf("failed to %s order book: client=null", messageType)
	}
	if c.executor == nil {
		return fmt.Errorf("failed to %s order book: executor=null", messageType)
	}
	key, err := params.Key()
	if err != nil {
		return fmt.Errorf("failed to %s order book: %w", messageType, err)
	}

	request := protocol.Request{
		Type:    messageType,
		Channel: protocol.ChannelOrderBook,
		ID:      strings.TrimSpace(params.Market),
	}
	if subscribe {
		request.Batched = &params.Batched
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to %s order book: failed to encode request: %w", messageType, err)
	}

	if subscribe {
		err = c.executor.Subscribe(ctx, key, payload)
	} else {
		err = c.executor.Unsubscribe(ctx, key, payload)
	}
	if err != nil {
		return fmt.Errorf("failed to %s order book: %w", messageType, err)
	}
	return nil
}
