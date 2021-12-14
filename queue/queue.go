package queue

import "context"

// IQueue 简单队列约定
type IQueue interface {
	Push(ctx context.Context, transNo string) error
	Pop(ctx context.Context) (string, error)
}
