package perpetual_markets

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type recordingExecutor struct {
	request transport.Request
	body    []byte
}

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return e.body, nil
}

func TestSendBuildsRequestAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{body: []byte(`{"markets":{"BTC-USD":{"clobPairId":"0","ticker":"BTC-USD","trades24H":12}}}`)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Send(context.Background(), Params{Market: "BTC-USD", Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if executor.request.Path != "/v4/perpetualMarkets" || executor.request.Query.Get("ticker") != "BTC-USD" || executor.request.Query.Get("limit") != "1" {
		t.Fatalf("unexpected request: %+v", executor.request)
	}
	market := result.Markets["BTC-USD"]
	if market.ClobPairID != "0" || market.Trades24H != 12 {
		t.Fatalf("unexpected market: %+v", market)
	}
}

func TestSendRejectsNegativeLimit(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	if _, err := client.Send(context.Background(), Params{Limit: -1}); err == nil {
		t.Fatal("expected an error")
	}
}
