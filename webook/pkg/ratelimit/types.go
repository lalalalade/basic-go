package ratelimit

import "context"

type Limiter interface {
	// Limit 有没有触发限流, key就是限流对象
	Limit(ctx context.Context, key string) (bool, error)
}
