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
	ChannelTrades    = "v4_trades"
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

// OrderBookContents contains bid and ask changes in receive order. Batches has
// one entry for ordinary contents and preserves every entry for batched contents.
type OrderBookContents struct {
	Bids    []PriceLevel     `json:"bids"`
	Asks    []PriceLevel     `json:"asks"`
	Batches []OrderBookBatch `json:"-"`
}

// OrderBookBatch contains one atomic group of bid and ask changes.
type OrderBookBatch struct {
	Bids []PriceLevel `json:"bids"`
	Asks []PriceLevel `json:"asks"`
}

// UnmarshalJSON decodes single or batched order book contents while preserving batch order.
//
// Version:
//   - 2026-08-22: Added support for batched order book contents.
func (c *OrderBookContents) UnmarshalJSON(data []byte) error {
	if c == nil {
		return fmt.Errorf("failed to decode order book contents: target=null")
	}
	*c = OrderBookContents{}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("failed to decode order book contents: value=empty")
	}
	if bytes.Equal(data, []byte("null")) {
		return fmt.Errorf("failed to decode order book contents: value=null")
	}

	switch data[0] {
	case '{':
		var batch OrderBookBatch
		if err := json.Unmarshal(data, &batch); err != nil {
			return fmt.Errorf("failed to decode order book contents object: %w", err)
		}
		c.Bids = batch.Bids
		c.Asks = batch.Asks
		c.Batches = []OrderBookBatch{batch}
		return nil
	case '[':
		var batches []OrderBookBatch
		if err := json.Unmarshal(data, &batches); err != nil {
			return fmt.Errorf("failed to decode batched order book contents: %w", err)
		}
		c.Batches = batches
		for _, batch := range batches {
			c.Bids = append(c.Bids, batch.Bids...)
			c.Asks = append(c.Asks, batch.Asks...)
		}
		return nil
	default:
		return fmt.Errorf("failed to decode order book contents: value=invalid")
	}
}

// TradesMessage is a typed v4_trades message.
type TradesMessage struct {
	Type         string         `json:"type"`
	ConnectionID string         `json:"connection_id,omitempty"`
	Channel      string         `json:"channel,omitempty"`
	MessageID    int64          `json:"message_id,omitempty"`
	ID           string         `json:"id,omitempty"`
	Version      string         `json:"version,omitempty"`
	Contents     TradesContents `json:"contents"`
}

// TradesContents contains public trades for a market.
type TradesContents struct {
	Trades []Trade `json:"trades"`
}

// UnmarshalJSON decodes single or batched trade contents and flattens trades in receive order.
//
// Version:
//   - 2026-08-17: Added support for batched trade contents.
func (c *TradesContents) UnmarshalJSON(data []byte) error {
	if c == nil {
		return fmt.Errorf("failed to decode trades contents: target=null")
	}
	*c = TradesContents{}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("failed to decode trades contents: value=empty")
	}
	if bytes.Equal(data, []byte("null")) {
		return fmt.Errorf("failed to decode trades contents: value=null")
	}

	type tradesContents TradesContents
	switch data[0] {
	case '{':
		var contents tradesContents
		if err := json.Unmarshal(data, &contents); err != nil {
			return fmt.Errorf("failed to decode trades contents object: %w", err)
		}
		*c = TradesContents(contents)
		return nil
	case '[':
		var batches []tradesContents
		if err := json.Unmarshal(data, &batches); err != nil {
			return fmt.Errorf("failed to decode batched trades contents: %w", err)
		}
		trades := make([]Trade, 0)
		for _, batch := range batches {
			trades = append(trades, batch.Trades...)
		}
		c.Trades = trades
		return nil
	default:
		return fmt.Errorf("failed to decode trades contents: value=invalid")
	}
}

// Trade contains one public market trade.
type Trade struct {
	ID              string `json:"id"`
	CreatedAtHeight string `json:"createdAtHeight,omitempty"`
	CreatedAt       string `json:"createdAt"`
	Side            string `json:"side"`
	Price           string `json:"price"`
	Size            string `json:"size"`
	Type            string `json:"type"`
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
