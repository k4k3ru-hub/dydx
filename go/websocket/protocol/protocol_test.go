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
	if message.Contents.Bids[0].Price != "100" || message.Contents.Asks[0].Offset != "9" || len(message.Contents.Batches) != 1 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestOrderBookMessageDecodesBatchedContentsInReceiveOrder(t *testing.T) {
	payload := []byte(`{"type":"channel_data","channel":"v4_orderbook","id":"BTC-USD","message_id":8,"contents":[{"bids":[["100","0"]],"asks":[]},{"bids":[],"asks":[["99","2"]]}]}`)
	var message OrderBookMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Batches) != 2 {
		t.Fatalf("batches = %d, want 2", len(message.Contents.Batches))
	}
	if message.Contents.Batches[0].Bids[0].Price != "100" || message.Contents.Batches[1].Asks[0].Price != "99" {
		t.Fatalf("unexpected batches: %+v", message.Contents.Batches)
	}
	if len(message.Contents.Bids) != 1 || len(message.Contents.Asks) != 1 {
		t.Fatalf("unexpected flattened contents: %+v", message.Contents)
	}
}

func TestOrderBookContentsRejectsInvalidValues(t *testing.T) {
	for _, payload := range []string{"", "null", "1", `"orderbook"`, "[{"} {
		var contents OrderBookContents
		if err := contents.UnmarshalJSON([]byte(payload)); err == nil {
			t.Fatalf("UnmarshalJSON(%q) error = nil", payload)
		}
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

func TestMarketsInitialMessageDecodesFundingOpenInterestAndOraclePrice(t *testing.T) {
	payload := []byte(`{"type":"subscribed","channel":"v4_markets","message_id":1,"contents":{"markets":{"BTC-USD":{"clobPairId":"0","ticker":"BTC-USD","status":"ACTIVE","oraclePrice":"60000.25","priceChange24H":"125.5","volume24H":"1000000","trades24H":321,"nextFundingRate":"0.0000125","initialMarginFraction":"0.05","maintenanceMarginFraction":"0.03","openInterest":"123.456789","atomicResolution":-10,"quantumConversionExponent":-9,"tickSize":"1","stepSize":"0.0001","stepBaseQuantums":1000000,"subticksPerTick":100000,"marketType":"CROSS","openInterestLowerCap":"1000000","openInterestUpperCap":"10000000","baseOpenInterest":"10","defaultFundingRate1H":"0"}}}}`)
	var message MarketsInitialMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	market := message.Contents.Markets["BTC-USD"]
	if market.NextFundingRate != "0.0000125" || market.OpenInterest != "123.456789" || market.OraclePrice != "60000.25" || market.DefaultFundingRate1H != "0" {
		t.Fatalf("unexpected initial market: %+v", market)
	}
}

func TestMarketsUpdateMessageDecodesFullPayload(t *testing.T) {
	payload := []byte(`{"type":"channel_data","channel":"v4_markets","message_id":2,"contents":{"trading":{"BTC-USD":{"id":"0","clobPairId":"0","ticker":"BTC-USD","status":"ACTIVE","baseAsset":"BTC","quoteAsset":"USD","marketId":0,"priceChange24H":"125.5","volume24H":"1000000","trades24H":321,"nextFundingRate":"0.000015","openInterest":"124.000001","baseOpenInterest":"10","basePositionSize":"1","incrementalPositionSize":"0.1","maxPositionSize":"1000","initialMarginFraction":"0.05","maintenanceMarginFraction":"0.03","atomicResolution":-10,"quantumConversionExponent":-9,"stepBaseQuantums":1000000,"subticksPerTick":100000}},"oraclePrices":{"BTC-USD":{"marketId":0,"oraclePrice":"60001.75","effectiveAt":"2026-08-27T00:00:00.000Z","effectiveAtHeight":"123456"}}}}`)
	var message MarketsUpdateMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	market := message.Contents.Trading["BTC-USD"]
	oracle := message.Contents.OraclePrices["BTC-USD"]
	if market.NextFundingRate != "0.000015" || market.OpenInterest != "124.000001" || market.BaseAsset != "BTC" {
		t.Fatalf("unexpected trading update: %+v", market)
	}
	if oracle.OraclePrice != "60001.75" || oracle.EffectiveAt != "2026-08-27T00:00:00.000Z" || oracle.EffectiveAtHeight != "123456" {
		t.Fatalf("unexpected oracle update: %+v", oracle)
	}
	if len(message.Contents.Batches) != 1 {
		t.Fatalf("batches = %d, want 1", len(message.Contents.Batches))
	}
}

func TestMarketsUpdateMessageDecodesBatchedContents(t *testing.T) {
	payload := []byte(`{"contents":[{"trading":{"BTC-USD":{"nextFundingRate":"0.00001","openInterest":"100"}}},{"trading":{"BTC-USD":{"nextFundingRate":"0.00002","openInterest":"101"}},"oraclePrices":{"BTC-USD":{"oraclePrice":"60002","effectiveAt":"2026-08-27T00:01:00.000Z","effectiveAtHeight":"123457"}}}]}`)
	var message MarketsUpdateMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if len(message.Contents.Batches) != 2 {
		t.Fatalf("batches = %d, want 2", len(message.Contents.Batches))
	}
	if message.Contents.Trading["BTC-USD"].OpenInterest != "101" || message.Contents.OraclePrices["BTC-USD"].OraclePrice != "60002" {
		t.Fatalf("unexpected flattened update: %+v", message.Contents)
	}
}

func TestMarketsUpdateContentsRejectsInvalidValues(t *testing.T) {
	for _, payload := range []string{"", "null", "1", `"markets"`, "[{"} {
		var contents MarketsUpdateContents
		if err := contents.UnmarshalJSON([]byte(payload)); err == nil {
			t.Fatalf("UnmarshalJSON(%q) error = nil", payload)
		}
	}
}
