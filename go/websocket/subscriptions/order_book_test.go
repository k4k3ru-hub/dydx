package subscriptions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/k4k3ru-hub/dydx/go/websocket/protocol"
)

type recordingExecutor struct {
	key          string
	payload      []byte
	unsubscribed bool
}

func (e *recordingExecutor) Subscribe(_ context.Context, key string, payload []byte) error {
	e.key, e.payload = key, payload
	return nil
}

func (e *recordingExecutor) Unsubscribe(_ context.Context, key string, payload []byte) error {
	e.key, e.payload, e.unsubscribed = key, payload, true
	return nil
}

func TestSubscribeBuildsDYdXRequest(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewOrderBookClient(executor)
	if err := client.Subscribe(context.Background(), OrderBookParams{Market: "BTC-USD"}); err != nil {
		t.Fatal(err)
	}
	var request protocol.Request
	if err := json.Unmarshal(executor.payload, &request); err != nil {
		t.Fatal(err)
	}
	if executor.key != "v4_orderbook:BTC-USD" || request.Type != "subscribe" || request.Channel != "v4_orderbook" || request.ID != "BTC-USD" || request.Batched == nil || *request.Batched {
		t.Fatalf("unexpected subscription: key=%s request=%+v", executor.key, request)
	}
}

func TestUnsubscribeOmitsBatched(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewOrderBookClient(executor)
	if err := client.Unsubscribe(context.Background(), OrderBookParams{Market: "BTC-USD", Batched: true}); err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(executor.payload, &request); err != nil {
		t.Fatal(err)
	}
	if !executor.unsubscribed || request["type"] != "unsubscribe" {
		t.Fatalf("unexpected unsubscribe: %+v", request)
	}
	if _, ok := request["batched"]; ok {
		t.Fatal("unsubscribe must omit batched")
	}
}

func TestOrderBookParamsRequiresMarket(t *testing.T) {
	if _, err := (OrderBookParams{Market: " "}).Key(); err == nil {
		t.Fatal("expected an error")
	}
}
