package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	MessageTypeSubscribe    = "subscribe"
	MessageTypeUnsubscribe  = "unsubscribe"
	MessageTypeConnected    = "connected"
	MessageTypeSubscribed   = "subscribed"
	MessageTypeUnsubscribed = "unsubscribed"
	MessageTypeChannelData  = "channel_data"

	ChannelOrderBook = "v4_orderbook"
)

// Request is a dYdX Indexer WebSocket subscription request.
type Request struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	ID      string `json:"id"`
	Batched *bool  `json:"batched,omitempty"`
}

// Message contains fields shared by dYdX Indexer WebSocket messages.
type Message struct {
	Type         string          `json:"type"`
	ConnectionID string          `json:"connection_id,omitempty"`
	Channel      string          `json:"channel,omitempty"`
	MessageID    int64           `json:"message_id,omitempty"`
	ID           string          `json:"id,omitempty"`
	Version      string          `json:"version,omitempty"`
	Contents     json.RawMessage `json:"contents,omitempty"`
}

// OrderBookMessage is a typed v4_orderbook message.
type OrderBookMessage struct {
	Type         string            `json:"type"`
	ConnectionID string            `json:"connection_id,omitempty"`
	Channel      string            `json:"channel,omitempty"`
	MessageID    int64             `json:"message_id,omitempty"`
	ID           string            `json:"id,omitempty"`
	Version      string            `json:"version,omitempty"`
	Contents     OrderBookContents `json:"contents"`
}

// OrderBookContents contains bid and ask changes. Snapshot levels are objects;
// incremental levels are tuples. PriceLevel accepts both wire representations.
type OrderBookContents struct {
	Bids []PriceLevel `json:"bids"`
	Asks []PriceLevel `json:"asks"`
}

type PriceLevel struct {
	Price  string `json:"price"`
	Size   string `json:"size"`
	Offset string `json:"offset,omitempty"`
}

func (p *PriceLevel) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("failed to decode price level: target=null")
	}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("failed to decode price level: value=empty")
	}

	if data[0] == '{' {
		type priceLevel PriceLevel
		var level priceLevel
		if err := json.Unmarshal(data, &level); err != nil {
			return fmt.Errorf("failed to decode price level object: %w", err)
		}
		*p = PriceLevel(level)
		return nil
	}

	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("failed to decode price level tuple: %w", err)
	}
	if len(values) != 2 && len(values) != 3 {
		return fmt.Errorf("failed to decode price level tuple: expected 2 or 3 values, got %d", len(values))
	}

	p.Price = values[0]
	p.Size = values[1]
	p.Offset = ""
	if len(values) == 3 {
		p.Offset = values[2]
	}
	return nil
}
