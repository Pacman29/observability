package metrics

import (
	"context"
	"iter"
	"maps"
	"time"

	"github.com/Pacman29/observability/internal/pool"
)

type metrics struct {
	driver   Driver
	tagsPool *pool.Map[string, string]
	o        *options
}

func New(d Driver, opts ...Option) Metrics {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &metrics{
		o:        o,
		driver:   d,
		tagsPool: pool.NewMap[string, string](o.tagsMapPoolSaveCapacity, o.tagsMapPoolCreateCapacity, o.defaultTags),
	}
}

type metricEventHandler struct {
	tags map[string]string
}

func (e *metricEventHandler) Tags() iter.Seq2[string, string] {
	return maps.All(e.tags)
}

func (e *metricEventHandler) GetTags() map[string]string {
	return maps.Clone(e.tags)
}

func (m *metrics) newHandler() (*metricEventHandler, func()) {
	h := &metricEventHandler{
		tags: m.tagsPool.Get(),
	}
	return h, func() {
		m.tagsPool.Save(h.tags)
	}
}

func (m *metrics) withArgs(ctx context.Context, handler *metricEventHandler) {
	// добавляем данные из ридеров
	for _, reader := range m.o.ctxReaders {
		maps.Copy(reader(ctx), handler.tags)
	}

	// потом все из контекста
	for k, v := range getTags(ctx) {
		handler.tags[k] = v
	}
}

func (m *metrics) Counter(ctx context.Context, key string, value float64) {
	ctx = defaultCtx(ctx)
	handler, releaser := m.newHandler()
	defer releaser()
	m.withArgs(ctx, handler)

	m.driver.Counter(ctx, handler, key, value)
}

func (m *metrics) Increment(ctx context.Context, key string, value float64) {
	ctx = defaultCtx(ctx)
	handler, releaser := m.newHandler()
	defer releaser()
	m.withArgs(ctx, handler)

	m.driver.Increment(ctx, handler, key, value)
}

func (m *metrics) Gauge(ctx context.Context, key string, value float64) {
	ctx = defaultCtx(ctx)
	handler, releaser := m.newHandler()
	defer releaser()
	m.withArgs(ctx, handler)

	m.driver.Gauge(ctx, handler, key, value)
}

func (m *metrics) Histogram(ctx context.Context, key string, buckets []float64, value float64) {
	ctx = defaultCtx(ctx)
	handler, releaser := m.newHandler()
	defer releaser()
	m.withArgs(ctx, handler)

	m.driver.Histogram(ctx, handler, key, buckets, value)
}

func (m *metrics) Timing(ctx context.Context, key string, ms int) {
	ctx = defaultCtx(ctx)
	handler, releaser := m.newHandler()
	defer releaser()
	m.withArgs(ctx, handler)

	m.driver.Timing(ctx, handler, key, ms)
}

func (m *metrics) Duration(ctx context.Context, key string, v time.Duration) {
	ctx = defaultCtx(ctx)
	handler, releaser := m.newHandler()
	defer releaser()
	m.withArgs(ctx, handler)

	m.driver.Duration(ctx, handler, key, v)
}

func (m *metrics) Timer() func(ctx context.Context, key string) {
	now := time.Now()
	return func(ctx context.Context, key string) {
		m.Duration(ctx, key, time.Since(now))
	}
}

func (m *metrics) WithTag(ctx context.Context, key string, value string) context.Context {
	ctx = defaultCtx(ctx)
	return addTagToCtx(ctx, key, value)
}

func (m *metrics) WithTags(ctx context.Context, tags map[string]string) context.Context {
	ctx = defaultCtx(ctx)
	return addTagsToCtx(ctx, tags)
}

func (m *metrics) Flush() {
	m.driver.Flush()
}

func (m *metrics) Close() {
	m.driver.Close()
}

func defaultCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
