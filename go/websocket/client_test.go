package websocket

import (
	"context"
	"net/http"
	"testing"

	k4websocket "github.com/k4k3ru-hub/websocket/go"
)

type testSessionHandler struct{}

func (*testSessionHandler) HandleMessage(SessionContext, []byte) {}
func (*testSessionHandler) HandleClose(SessionContext)           {}

func TestNewClientComposesOrderBookAndPreservesOptions(t *testing.T) {
	option := &ClientOption{
		EndpointURL: "wss://example.test/ws",
		HTTPHeader: http.Header{
			"X-Test": {"original"},
		},
		SessionOption: &k4websocket.SessionOption{},
	}
	client, err := NewClient(context.Background(), DefaultEndpointURL, &testSessionHandler{}, option)
	if err != nil {
		t.Fatal(err)
	}
	if client.OrderBook() == nil {
		t.Fatal("expected composed order book client")
	}
	if option.HTTPHeader.Get("X-Test") != "original" {
		t.Fatalf("caller-owned option was modified: %+v", option)
	}
}

func TestNewClientRequiresHandler(t *testing.T) {
	if _, err := NewClient(context.Background(), DefaultEndpointURL, nil, nil); err == nil {
		t.Fatal("expected an error")
	}
}
