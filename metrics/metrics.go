package metrics

import (
	"context"
	"time"
)

type metrics struct {
}

func New(d Driver) Metrics {
	return &metrics{}
}

func (m *metrics) Counter(ctx context.Context, key string, value int) {

}

func (m *metrics) Increment(ctx context.Context, key string, value int) {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Gauge(ctx context.Context, key string, value any) {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Histogram(ctx context.Context, bucket string, value any) {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Timing(ctx context.Context, key string, ms int) {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Duration(ctx context.Context, key string, v time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Timer() func(ctx context.Context, key string) {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) WithTag(key string, value string) context.Context {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) WithTags(m map[string]string) context.Context {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Flush() {
	//TODO implement me
	panic("implement me")
}

func (m *metrics) Close() {
	//TODO implement me
	panic("implement me")
}
