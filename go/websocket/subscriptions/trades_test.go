package subscriptions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/k4k3ru-hub/dydx/go/websocket/protocol"
)

func TestTradesSubscribeBuildsDYdXRequest(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewTradesClient(executor)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Subscribe(context.Background(), TradesParams{Market: " BTC-USD ", Batched: true}); err != nil {
		t.Fatal(err)
	}
	var request protocol.Request
	if err := json.Unmarshal(executor.payload, &request); err != nil {
		t.Fatal(err)
	}
	if executor.key != "v4_trades:BTC-USD" || request.Type != protocol.MessageTypeSubscribe || request.Channel != protocol.ChannelTrades || request.ID != "BTC-USD" || request.Batched == nil || !*request.Batched {
		t.Fatalf("unexpected subscription: key=%s request=%+v", executor.key, request)
	}
}

func TestTradesUnsubscribeOmitsBatched(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewTradesClient(executor)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Unsubscribe(context.Background(), TradesParams{Market: "BTC-USD", Batched: true}); err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(executor.payload, &request); err != nil {
		t.Fatal(err)
	}
	if !executor.unsubscribed || request["type"] != protocol.MessageTypeUnsubscribe {
		t.Fatalf("unexpected unsubscribe: %+v", request)
	}
	if _, ok := request["batched"]; ok {
		t.Fatal("unsubscribe must omit batched")
	}
}

func TestTradesParamsRequiresMarket(t *testing.T) {
	if _, err := (TradesParams{Market: " "}).Key(); err == nil {
		t.Fatal("expected an error")
	}
}
