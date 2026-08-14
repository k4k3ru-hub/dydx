package order_book

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/dydx/go/rest/transport"
)

type recordingExecutor struct {
	request transport.Request
}

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{"bids":[{"price":"100","size":"2"}],"asks":[{"price":"101","size":"3"}]}`), nil
}

func TestSendBuildsRequestAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewClient(executor)
	result, err := client.Send(context.Background(), Params{Market: " BTC-USD "})
	if err != nil {
		t.Fatal(err)
	}
	if executor.request.Path != "/v4/orderbooks/perpetualMarket/BTC-USD" {
		t.Fatalf("unexpected path: %s", executor.request.Path)
	}
	if len(result.Bids) != 1 || result.Bids[0].Price != "100" || len(result.Asks) != 1 {
		t.Fatalf("unexpected order book: %+v", result)
	}
}

func TestSendRequiresMarket(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	if _, err := client.Send(context.Background(), Params{Market: "  "}); err == nil {
		t.Fatal("expected an error")
	}
}
