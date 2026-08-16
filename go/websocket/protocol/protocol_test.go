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
	if trade.ID != "trade-1" || trade.CreatedAtHeight != "123" || trade.Side != "BUY" || trade.Price != "100" || trade.Size != "2" || trade.Type != "LIMIT" {
		t.Fatalf("unexpected trade: %+v", trade)
	}
}
