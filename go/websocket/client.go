package websocket

import (
	"context"
	"fmt"
	"net/http"

	"github.com/k4k3ru-hub/dydx/go/websocket/subscriptions"
	k4websocket "github.com/k4k3ru-hub/websocket/go"
)

const DefaultEndpointURL = "wss://indexer.dydx.trade/v4/ws"

type ClientOption = k4websocket.ClientOption
type SessionHandler = k4websocket.SessionHandler
type SessionContext = k4websocket.SessionContext

type Client struct {
	wsClient        *k4websocket.Client
	orderBookClient *subscriptions.OrderBookClient
}

func DefaultClientOption() *ClientOption {
	option := k4websocket.DefaultClientOption()
	option.EndpointURL = DefaultEndpointURL
	return option
}

func NewClient(ctx context.Context, endpointURL string, handler SessionHandler, option *ClientOption) (*Client, error) {
	if handler == nil {
		return nil, fmt.Errorf("failed to create client: missing required parameter: session_handler=null")
	}
	cloned := cloneClientOption(option)
	if cloned.EndpointURL != "" {
		endpointURL = cloned.EndpointURL
	}
	if endpointURL == "" {
		endpointURL = DefaultEndpointURL
	}

	wsClient, err := k4websocket.NewClient(ctx, endpointURL, handler, cloned)
	if err != nil {
		return nil, err
	}
	client := &Client{wsClient: wsClient}
	orderBookClient, err := subscriptions.NewOrderBookClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}
	client.orderBookClient = orderBookClient
	return client, nil
}

func cloneClientOption(option *ClientOption) *ClientOption {
	if option == nil {
		return DefaultClientOption()
	}
	cloned := &ClientOption{
		EndpointURL: option.EndpointURL, ConnectTimeout: option.ConnectTimeout,
		HandshakeTimeout: option.HandshakeTimeout,
	}
	if option.HTTPHeader != nil {
		cloned.HTTPHeader = option.HTTPHeader.Clone()
	} else {
		cloned.HTTPHeader = make(http.Header)
	}
	if option.SessionOption != nil {
		cloned.SessionOption = option.SessionOption.Clone()
	} else {
		cloned.SessionOption = k4websocket.DefaultSessionOption()
	}
	return cloned
}

func (c *Client) OrderBook() *subscriptions.OrderBookClient {
	if c == nil {
		return nil
	}
	return c.orderBookClient
}

func (c *Client) Connect(ctx context.Context) error {
	if c == nil || c.wsClient == nil {
		return fmt.Errorf("failed to connect websocket: client=null")
	}
	return c.wsClient.Connect(ctx)
}

func (c *Client) Close() error {
	if c == nil || c.wsClient == nil {
		return fmt.Errorf("failed to close websocket: client=null")
	}
	return c.wsClient.Close()
}

func (c *Client) SessionContext() (SessionContext, error) {
	if c == nil || c.wsClient == nil {
		return nil, fmt.Errorf("failed to get websocket session context: client=null")
	}
	return c.wsClient.SessionContext()
}

func (c *Client) Subscribe(ctx context.Context, key string, payload []byte) error {
	if c == nil || c.wsClient == nil {
		return fmt.Errorf("failed to subscribe: client=null")
	}
	return c.wsClient.Subscribe(ctx, key, payload)
}

func (c *Client) Unsubscribe(ctx context.Context, key string, payload []byte) error {
	if c == nil || c.wsClient == nil {
		return fmt.Errorf("failed to unsubscribe: client=null")
	}
	return c.wsClient.Unsubscribe(ctx, key, payload)
}
