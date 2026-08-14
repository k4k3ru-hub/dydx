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
