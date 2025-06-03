package metrics

import (
	"context"
	"iter"
	"time"
)

type Driver interface {
	Counter(ctx context.Context, handler EventHandler)
	Increment(ctx context.Context, handler EventHandler)
	Gauge(ctx context.Context, handler EventHandler)
	Histogram(ctx context.Context, handler EventHandler)
	Timing(ctx context.Context, handler EventHandler)
	Duration(ctx context.Context, handler EventHandler)
	Flush()
	Close()
}

type CtxReader func(ctx context.Context) map[string]string

type Metrics interface {
	Counter(ctx context.Context, key string, value float64, opts ...MetricOption)
	Increment(ctx context.Context, key string, value float64, opts ...MetricOption)
	Gauge(ctx context.Context, key string, value float64, opts ...MetricOption)
	Histogram(ctx context.Context, key string, value float64, opts ...MetricOption)
	Timing(ctx context.Context, key string, ms int, opts ...MetricOption)
	Duration(ctx context.Context, key string, v time.Duration, opts ...MetricOption)
	Timer() func(ctx context.Context, key string, opts ...MetricOption)
	WithTag(ctx context.Context, key string, value string) context.Context
	WithTags(ctx context.Context, m map[string]string) context.Context
	Flush()
	Close()
}

type EventHandler interface {
	Tags() iter.Seq2[string, string]
	GetBuckets() []float64
	GetTags() map[string]string
	GetValue() float64
	GetKey() string
}
