package subscriptions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/k4k3ru-hub/dydx/go/websocket/protocol"
)

func TestMarketsSubscribeBuildsDYdXRequestWithoutID(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewMarketsClient(executor)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Subscribe(context.Background(), MarketsParams{Batched: true}); err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(executor.payload, &request); err != nil {
		t.Fatal(err)
	}
	if executor.key != protocol.ChannelMarkets || request["type"] != protocol.MessageTypeSubscribe || request["channel"] != protocol.ChannelMarkets || request["batched"] != true {
		t.Fatalf("unexpected subscription: key=%s request=%+v", executor.key, request)
	}
	if _, ok := request["id"]; ok {
		t.Fatal("markets subscription must omit id")
	}
}

func TestMarketsUnsubscribeOmitsIDAndBatched(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewMarketsClient(executor)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Unsubscribe(context.Background(), MarketsParams{Batched: true}); err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(executor.payload, &request); err != nil {
		t.Fatal(err)
	}
	if !executor.unsubscribed || request["type"] != protocol.MessageTypeUnsubscribe {
		t.Fatalf("unexpected unsubscribe: %+v", request)
	}
	if _, ok := request["id"]; ok {
		t.Fatal("markets unsubscribe must omit id")
	}
	if _, ok := request["batched"]; ok {
		t.Fatal("markets unsubscribe must omit batched")
	}
}
