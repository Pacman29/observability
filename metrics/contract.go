package metrics

import (
	"context"
	"iter"
	"time"
)

type Driver interface {
	Counter(ctx context.Context, handler EventHandler, key string, value float64)
	Increment(ctx context.Context, handler EventHandler, key string, value float64)
	Gauge(ctx context.Context, handler EventHandler, key string, value float64)
	Histogram(ctx context.Context, handler EventHandler, key string, buckets []float64, value float64)
	Timing(ctx context.Context, handler EventHandler, key string, ms int)
	Duration(ctx context.Context, handler EventHandler, key string, v time.Duration)
	Flush()
	Close()
}

type CtxReader func(ctx context.Context) map[string]string

type Metrics interface {
	Counter(ctx context.Context, key string, value float64)
	Increment(ctx context.Context, key string, value float64)
	Gauge(ctx context.Context, key string, value float64)
	Histogram(ctx context.Context, key string, buckets []float64, value float64)
	Timing(ctx context.Context, key string, ms int)
	Duration(ctx context.Context, key string, v time.Duration)
	Timer() func(ctx context.Context, key string)
	WithTag(ctx context.Context, key string, value string) context.Context
	WithTags(ctx context.Context, m map[string]string) context.Context
	Flush()
	Close()
}

type EventHandler interface {
	Tags() iter.Seq2[string, string]
	GetTags() map[string]string
}
