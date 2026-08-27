# dYdX SDK

dYdX Indexer の公開 Market Data API を扱う Go SDK です。

現在は次の API を提供します。

- REST: perpetual market metadata
- REST: perpetual market order-book snapshot
- REST: perpetual market historical funding rates
- WebSocket: `v4_orderbook` subscribe / unsubscribe
- WebSocket: `v4_trades` subscribe / unsubscribe
- WebSocket: `v4_markets` Funding Rate・Open Interest・Oracle Price updates

実装と利用例は [`go`](./go) を参照してください。

## License

MIT
