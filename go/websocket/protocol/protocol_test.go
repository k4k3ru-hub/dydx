package protocol

import (
	"encoding/json"
	"testing"
)

func TestOrderBookMessageAcceptsSnapshotAndIncrementalLevels(t *testing.T) {
	payload := []byte(`{"type":"channel_data","channel":"v4_orderbook","id":"BTC-USD","message_id":7,"contents":{"bids":[{"price":"100","size":"2"}],"asks":[["101","3","9"]]}}`)
	var message OrderBookMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if message.Contents.Bids[0].Price != "100" || message.Contents.Asks[0].Offset != "9" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestPriceLevelRejectsInvalidTuple(t *testing.T) {
	var level PriceLevel
	if err := json.Unmarshal([]byte(`["100"]`), &level); err == nil {
		t.Fatal("expected an error")
	}
}

func TestTradesMessageDecodesSnapshotAndUpdateFields(t *testing.T) {
	payload := []byte(`{"type":"channel_data","channel":"v4_trades","id":"BTC-USD","message_id":8,"contents":{"trades":[{"id":"trade-1","createdAtHeight":"123","createdAt":"2026-08-17T00:00:00.000Z","side":"BUY","price":"100","size":"2","type":"LIMIT"}]}}`)
	var message TradesMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Trades) != 1 {
		t.Fatalf("unexpected trades: %+v", message.Contents.Trades)
	}
	trade := message.Contents.Trades[0]
	if trade.ID != "trade-1" || trade.CreatedAtHeight != "123" || trade.CreatedAt != "2026-08-17T00:00:00.000Z" || trade.Side != "BUY" || trade.Price != "100" || trade.Size != "2" || trade.Type != "LIMIT" {
		t.Fatalf("unexpected trade: %+v", trade)
	}
}

func TestTradesMessageDecodesBatchedContentsInReceiveOrder(t *testing.T) {
	payload := []byte(`{"type":"channel_data","channel":"v4_trades","id":"BTC-USD","message_id":123,"contents":[{"trades":[{"id":"trade-1","createdAtHeight":"123","createdAt":"2026-08-17T00:00:00.000Z","side":"BUY","price":"100","size":"2","type":"LIMIT"}]},{"trades":[{"id":"trade-2","createdAtHeight":"124","createdAt":"2026-08-17T00:00:01.000Z","side":"SELL","price":"101","size":"3","type":"MARKET"}]}]}`)
	var message TradesMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Trades) != 2 {
		t.Fatalf("unexpected trades: %+v", message.Contents.Trades)
	}
	first := message.Contents.Trades[0]
	if first.ID != "trade-1" || first.CreatedAtHeight != "123" || first.CreatedAt != "2026-08-17T00:00:00.000Z" || first.Side != "BUY" || first.Price != "100" || first.Size != "2" || first.Type != "LIMIT" {
		t.Fatalf("unexpected first trade: %+v", first)
	}
	second := message.Contents.Trades[1]
	if second.ID != "trade-2" || second.CreatedAtHeight != "124" || second.CreatedAt != "2026-08-17T00:00:01.000Z" || second.Side != "SELL" || second.Price != "101" || second.Size != "3" || second.Type != "MARKET" {
		t.Fatalf("unexpected second trade: %+v", second)
	}
}

func TestTradesMessageDecodesEmptyBatchedContents(t *testing.T) {
	var message TradesMessage
	if err := json.Unmarshal([]byte(`{"contents":[]}`), &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Trades) != 0 {
		t.Fatalf("unexpected trades: %+v", message.Contents.Trades)
	}
}

func TestTradesMessageDecodesBatchContainingEmptyTrades(t *testing.T) {
	var message TradesMessage
	if err := json.Unmarshal([]byte(`{"contents":[{"trades":[]},{"trades":[{"id":"trade-1"}]}]}`), &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Trades) != 1 || message.Contents.Trades[0].ID != "trade-1" {
		t.Fatalf("unexpected trades: %+v", message.Contents.Trades)
	}
}

func TestTradesContentsRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{name: "empty", payload: ""},
		{name: "invalid JSON", payload: `[{`},
		{name: "null", payload: `null`},
		{name: "number", payload: `1`},
		{name: "string", payload: `"trades"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var contents TradesContents
			if err := contents.UnmarshalJSON([]byte(test.payload)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestTradesMessageRejectsInvalidContents(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{name: "null", payload: `{"contents":null}`},
		{name: "number", payload: `{"contents":1}`},
		{name: "string", payload: `{"contents":"trades"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var message TradesMessage
			if err := json.Unmarshal([]byte(test.payload), &message); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestTradesMessageDecodeReuseClearsPreviousTrades(t *testing.T) {
	var message TradesMessage
	if err := json.Unmarshal([]byte(`{"contents":{"trades":[{"id":"old"}]}}`), &message); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(`{"contents":[{"trades":[{"id":"new"}]}]}`), &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Trades) != 1 || message.Contents.Trades[0].ID != "new" {
		t.Fatalf("unexpected trades: %+v", message.Contents.Trades)
	}
}
