package subscriptions

import "context"

type Executor interface {
	Subscribe(ctx context.Context, key string, payload []byte) error
	Unsubscribe(ctx context.Context, key string, payload []byte) error
}
