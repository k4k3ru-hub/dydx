# dYdX Go SDK

## Install

```bash
go get github.com/k4k3ru-hub/dydx/go/rest
go get github.com/k4k3ru-hub/dydx/go/websocket
```

Import only the packages required by your application.

## REST

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/k4k3ru-hub/dydx/go/rest"
	"github.com/k4k3ru-hub/dydx/go/rest/markets/order_book"
	"github.com/k4k3ru-hub/dydx/go/rest/markets/perpetual_markets"
)

func main() {
	client, err := rest.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	markets, err := client.Markets().PerpetualMarkets().Send(
		context.Background(),
		perpetual_markets.Params{Market: "BTC-USD"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markets.Markets["BTC-USD"].OraclePrice)

	book, err := client.Markets().OrderBook().Send(
		context.Background(),
		order_book.Params{Market: "BTC-USD"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(book.Bids[0], book.Asks[0])
}
```

The default REST endpoint is `https://indexer.dydx.trade`.

## WebSocket market data

Implement `websocket.SessionHandler`, construct a client, then subscribe. The
example imports the required packages directly:

```go
import (
	"github.com/k4k3ru-hub/dydx/go/websocket"
	"github.com/k4k3ru-hub/dydx/go/websocket/subscriptions"
)

client, err := websocket.NewClient(ctx, websocket.DefaultEndpointURL, handler, nil)
if err != nil {
	return err
}
if err := client.Connect(ctx); err != nil {
	return err
}
defer client.Close()

err = client.OrderBook().Subscribe(ctx, subscriptions.OrderBookParams{
	Market:  "BTC-USD",
	Batched: false,
})

err = client.Trades().Subscribe(ctx, subscriptions.TradesParams{
	Market:  "BTC-USD",
	Batched: false,
})
```

Decode order-book events into `protocol.OrderBookMessage` from
`github.com/k4k3ru-hub/dydx/go/websocket/protocol`. Its price levels support both
snapshot objects (`{"price":"...","size":"..."}`) and incremental tuples
(`["price","size","offset"]`). The default WebSocket endpoint is
`wss://indexer.dydx.trade/v4/ws`.

Decode trade events into `protocol.TradesMessage`. Each entry in
`Contents.Trades` contains the trade ID, creation time and height, side, price,
size, and type.
