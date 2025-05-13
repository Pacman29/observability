package metrics

import (
	"context"
	"time"
)

type Driver interface {
}

type Metrics interface {
	Counter(ctx context.Context, key string, value int)
	Increment(ctx context.Context, key string, value int)
	Gauge(ctx context.Context, key string, value any)
	Histogram(ctx context.Context, bucket string, value any)
	Timing(ctx context.Context, key string, ms int)
	Duration(ctx context.Context, key string, v time.Duration)
	Timer() func(ctx context.Context, key string)
	WithTag(key string, value string) context.Context
	WithTags(m map[string]string) context.Context
	Flush()
	Close()
}
