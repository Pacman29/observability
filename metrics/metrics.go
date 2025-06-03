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
	tags    map[string]string
	buckets []float64
	value   float64
	key     string
}

func (e *metricEventHandler) Tags() iter.Seq2[string, string] {
	return maps.All(e.tags)
}

func (e *metricEventHandler) GetTags() map[string]string {
	return maps.Clone(e.tags)
}

func (e *metricEventHandler) GetBuckets() []float64 {
	return e.buckets
}

func (e *metricEventHandler) GetKey() string {
	return e.key
}

func (e *metricEventHandler) GetValue() float64 {
	return e.value
}

func (m *metrics) withArgs(ctx context.Context, handler *metricEventHandler, opts ...MetricOption) {
	o := newMetricOption()
	for _, opt := range opts {
		opt(o)
	}

	// добавляем данные из ридеров
	for _, reader := range m.o.ctxReaders {
		maps.Copy(handler.tags, reader(ctx))
	}

	// потом все из контекста
	for k, v := range getTags(ctx) {
		handler.tags[k] = v
	}

	// потом все из опций
	for k, v := range o.tags {
		handler.tags[k] = v
	}

	if len(o.buckets) != 0 {
		handler.buckets = o.buckets
	}
}

func (m *metrics) Counter(ctx context.Context, key string, value float64, opts ...MetricOption) {
	ctx = defaultCtx(ctx)
	handler := &metricEventHandler{
		tags:    m.tagsPool.Get(),
		buckets: nil,
		value:   value,
		key:     key,
	}
	defer func() {
		m.tagsPool.Save(handler.tags)
	}()
	m.withArgs(ctx, handler, opts...)

	m.driver.Counter(ctx, handler)
}

func (m *metrics) Increment(ctx context.Context, key string, value float64, opts ...MetricOption) {
	ctx = defaultCtx(ctx)
	handler := &metricEventHandler{
		tags:    m.tagsPool.Get(),
		buckets: nil,
		value:   value,
		key:     key,
	}
	defer func() {
		m.tagsPool.Save(handler.tags)
	}()
	m.withArgs(ctx, handler, opts...)

	m.driver.Increment(ctx, handler)
}

func (m *metrics) Gauge(ctx context.Context, key string, value float64, opts ...MetricOption) {
	ctx = defaultCtx(ctx)
	handler := &metricEventHandler{
		tags:    m.tagsPool.Get(),
		buckets: nil,
		value:   value,
		key:     key,
	}
	defer func() {
		m.tagsPool.Save(handler.tags)
	}()
	m.withArgs(ctx, handler, opts...)

	m.driver.Gauge(ctx, handler)
}

func (m *metrics) Histogram(ctx context.Context, key string, value float64, opts ...MetricOption) {
	ctx = defaultCtx(ctx)
	handler := &metricEventHandler{
		tags:    m.tagsPool.Get(),
		buckets: nil,
		value:   value,
		key:     key,
	}
	defer func() {
		m.tagsPool.Save(handler.tags)
	}()
	m.withArgs(ctx, handler, opts...)

	m.driver.Histogram(ctx, handler)
}

func (m *metrics) Timing(ctx context.Context, key string, ms int, opts ...MetricOption) {
	ctx = defaultCtx(ctx)
	handler := &metricEventHandler{
		tags:    m.tagsPool.Get(),
		buckets: nil,
		value:   float64(ms),
		key:     key,
	}
	defer func() {
		m.tagsPool.Save(handler.tags)
	}()
	m.withArgs(ctx, handler, opts...)

	m.driver.Timing(ctx, handler)
}

func (m *metrics) Duration(ctx context.Context, key string, v time.Duration, opts ...MetricOption) {
	ctx = defaultCtx(ctx)
	handler := &metricEventHandler{
		tags:    m.tagsPool.Get(),
		buckets: nil,
		value:   float64(v.Milliseconds()),
		key:     key,
	}
	defer func() {
		m.tagsPool.Save(handler.tags)
	}()
	m.withArgs(ctx, handler, opts...)

	m.driver.Duration(ctx, handler)
}

func (m *metrics) Timer() func(ctx context.Context, key string, opts ...MetricOption) {
	now := time.Now()
	return func(ctx context.Context, key string, opts ...MetricOption) {
		m.Duration(ctx, key, time.Since(now), opts...)
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
