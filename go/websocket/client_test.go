package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/k4k3ru-hub/dydx/go/websocket/subscriptions"
	k4websocket "github.com/k4k3ru-hub/websocket/go"
)

type testSessionHandler struct{}

func (*testSessionHandler) HandleMessage(SessionContext, []byte) {}
func (*testSessionHandler) HandleClose(SessionContext)           {}

// TestNewClientComposesSubscriptionsAndPreservesOptions verifies constructor composition.
//
// Version:
//   - 2026-10-29: Include caller-owned compression options.
func TestNewClientComposesSubscriptionsAndPreservesOptions(t *testing.T) {
	option := &ClientOption{
		EndpointURL:       "wss://example.test/ws",
		EnableCompression: true,
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
	if client.Trades() == nil {
		t.Fatal("expected composed trades client")
	}
	if client.Markets() == nil {
		t.Fatal("expected composed markets client")
	}
	if option.HTTPHeader.Get("X-Test") != "original" {
		t.Fatalf("caller-owned option was modified: %+v", option)
	}
	if !option.EnableCompression {
		t.Fatal("caller-owned compression setting was modified")
	}
}

// TestClientCompressedOrderBook verifies option forwarding and subscription I/O.
//
// Version:
//   - 2026-10-29: Added.
func TestClientCompressedOrderBook(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload := `{"type":"subscribed","channel":"v4_orderbook","id":"BTC-USD","contents":{"bids":[{"price":"100","size":"2"}],"asks":[{"price":"101","size":"3"}]}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Sec-WebSocket-Extensions"), "permessage-deflate") {
			http.Error(w, "permessage-deflate required", http.StatusUpgradeRequired)
			return
		}
		upgrader := gorillaWebsocket.Upgrader{EnableCompression: true}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				t.Errorf("server close failed: %v", err)
			}
		}()
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Errorf("server read deadline failed: %v", err)
			return
		}
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("subscription receive failed: %v", err)
			return
		}
		var request struct {
			Type    string `json:"type"`
			Channel string `json:"channel"`
			ID      string `json:"id"`
			Batched bool   `json:"batched"`
		}
		if err := json.Unmarshal(raw, &request); err != nil {
			t.Errorf("subscription decode failed: %v", err)
			return
		}
		if request.Type != "subscribe" || request.Channel != "v4_orderbook" || request.ID != "BTC-USD" || !request.Batched {
			t.Errorf("unexpected subscription: %+v", request)
			return
		}
		conn.EnableWriteCompression(true)
		if err := conn.WriteMessage(gorillaWebsocket.TextMessage, []byte(payload)); err != nil {
			t.Errorf("snapshot send failed: %v", err)
		}
	}))
	defer server.Close()
	option := DefaultClientOption()
	option.EndpointURL = "ws" + strings.TrimPrefix(server.URL, "http")
	option.EnableCompression = true
	handler := &compressionTestHandler{messages: make(chan string, 1)}
	client, err := NewClient(ctx, DefaultEndpointURL, handler, option)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("client close failed: %v", err)
		}
	}()
	if err := client.OrderBook().Subscribe(ctx, subscriptions.OrderBookParams{Market: "BTC-USD", Batched: true}); err != nil {
		t.Fatal(err)
	}
	select {
	case received := <-handler.messages:
		if received != payload {
			t.Fatalf("decoded snapshot = %q, want %q", received, payload)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for decoded snapshot")
	}
}

type compressionTestHandler struct {
	messages chan string
}

// HandleMessage captures a decoded dYdX message.
//
// Version:
//   - 2026-10-29: Added.
func (h *compressionTestHandler) HandleMessage(_ SessionContext, payload []byte) {
	h.messages <- string(payload)
}

// HandleClose accepts completion of a test session.
//
// Version:
//   - 2026-10-29: Added.
func (*compressionTestHandler) HandleClose(SessionContext) {}

func TestNewClientRequiresHandler(t *testing.T) {
	if _, err := NewClient(context.Background(), DefaultEndpointURL, nil, nil); err == nil {
		t.Fatal("expected an error")
	}
}
