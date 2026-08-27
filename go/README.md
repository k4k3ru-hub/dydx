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
	"github.com/k4k3ru-hub/dydx/go/rest/markets/historical_funding"
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

	funding, err := client.Markets().HistoricalFunding().Send(
		context.Background(),
		historical_funding.Params{Market: "BTC-USD", Limit: 100},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(funding.HistoricalFunding[0].Rate)
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

err = client.Markets().Subscribe(ctx, subscriptions.MarketsParams{
	Batched: true,
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

The `v4_markets` channel does not use a market ID because it streams all dYdX
markets. Decode its initial snapshot into `protocol.MarketsInitialMessage` and
incremental updates into `protocol.MarketsUpdateMessage`. The initial snapshot
contains `nextFundingRate`, `openInterest`, and `oraclePrice`. Incremental
`trading` updates contain the current funding estimate and open interest, while
`oraclePrices` updates include the oracle price, effective time, and block
height.

For startup state, use REST `PerpetualMarkets`; use `v4_markets` for subsequent
updates and REST `HistoricalFunding` for settled funding records. Values remain
decimal strings so consumers can normalize them without losing precision.
`oraclePrice` is an oracle price and must not be labeled as an exchange mark
price. Dataset conversion and timestamp policy belong to the consuming service.
