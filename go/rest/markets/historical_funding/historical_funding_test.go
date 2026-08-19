package historical_funding

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
	executor := &recordingExecutor{body: []byte(`{"historicalFunding":[{"ticker":"BTC-USD","effectiveAt":"2026-08-19T00:00:00.000Z","effectiveAtHeight":"123","price":"60000","rate":"0.0001"}]}`)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Send(context.Background(), Params{
		Market:                    " BTC-USD ",
		EffectiveBeforeOrAt:       "2026-08-19T01:00:00.000Z",
		EffectiveBeforeOrAtHeight: 456,
		Limit:                     1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if executor.request.Path != "/v4/historicalFunding/BTC-USD" ||
		executor.request.Query.Get("effectiveBeforeOrAt") != "2026-08-19T01:00:00.000Z" ||
		executor.request.Query.Get("effectiveBeforeOrAtHeight") != "456" ||
		executor.request.Query.Get("limit") != "1" {
		t.Fatalf("unexpected request: %+v", executor.request)
	}
	if len(result.HistoricalFunding) != 1 {
		t.Fatalf("unexpected historical funding: %+v", result)
	}
	funding := result.HistoricalFunding[0]
	if funding.Ticker != "BTC-USD" || funding.EffectiveAtHeight != "123" || funding.Price != "60000" || funding.Rate != "0.0001" {
		t.Fatalf("unexpected historical funding: %+v", funding)
	}
}

func TestSendRejectsInvalidParams(t *testing.T) {
	client, err := NewClient(&recordingExecutor{})
	if err != nil {
		t.Fatal(err)
	}
	tests := []Params{
		{Market: " "},
		{Market: "BTC-USD", EffectiveBeforeOrAtHeight: -1},
		{Market: "BTC-USD", Limit: -1},
	}
	for _, params := range tests {
		if _, err := client.Send(context.Background(), params); err == nil {
			t.Fatalf("expected an error: %+v", params)
		}
	}
}
