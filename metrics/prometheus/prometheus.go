package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Pacman29/observability/internal/pool"
	"github.com/Pacman29/observability/metrics"
)

type driver struct {
	o          *options
	registerer prometheus.Registerer
	counters   *metricWrapper[prometheus.CounterVec]
	labelsPool *pool.Slice[string]
	gauge      *metricWrapper[prometheus.GaugeVec]
	histogram  *metricWrapper[prometheus.HistogramVec]
}

func NewPrometheusDriver(registerer prometheus.Registerer, opts ...Option) metrics.Driver {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &driver{
		o:          o,
		registerer: registerer,
		labelsPool: pool.NewSlice[string](o.labelsPoolCapSave, o.labelsPoolCapCreate, nil),

		counters:  newWrapper[prometheus.CounterVec](o.metricsSize),
		gauge:     newWrapper[prometheus.GaugeVec](o.metricsSize),
		histogram: newWrapper[prometheus.HistogramVec](o.metricsSize),
	}
}

func (d *driver) Counter(ctx context.Context, handler metrics.EventHandler, key string, value float64) {
	counter, labelNames, ok := d.counters.Get(key)

	labelValues := d.labelsPool.Get()
	defer d.labelsPool.Save(labelValues)

	if !ok {
		labelNames = d.labelsPool.Get()

		for k, v := range handler.Tags() {
			labelValues = append(labelValues, v)
			labelNames = append(labelNames, k)
		}

		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: d.o.namespace,
			Subsystem: d.o.subsystem,
			Name:      key,
		}, labelNames)
		d.counters.Add(key, counter, labelNames)
		d.registerer.MustRegister(counter)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	counter.WithLabelValues(labelValues...).Add(float64(value))
}

func (d *driver) Increment(ctx context.Context, handler metrics.EventHandler, key string, value float64) {
	counter, labelNames, ok := d.counters.Get(key)

	labelValues := d.labelsPool.Get()
	defer d.labelsPool.Save(labelValues)

	if !ok {
		labelNames = d.labelsPool.Get()

		for k, v := range handler.Tags() {
			labelValues = append(labelValues, v)
			labelNames = append(labelNames, k)
		}

		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: d.o.namespace,
			Subsystem: d.o.subsystem,
			Name:      key,
		}, labelNames)
		d.counters.Add(key, counter, labelNames)
		d.registerer.MustRegister(counter)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	counter.WithLabelValues(labelValues...).Add(float64(value))
}

func (d *driver) Gauge(ctx context.Context, handler metrics.EventHandler, key string, value float64) {
	gauger, labelNames, ok := d.gauge.Get(key)

	labelValues := d.labelsPool.Get()
	defer d.labelsPool.Save(labelValues)

	if !ok {
		labelNames = d.labelsPool.Get()

		for k, v := range handler.Tags() {
			labelValues = append(labelValues, v)
			labelNames = append(labelNames, k)
		}

		gauger = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: d.o.namespace,
			Subsystem: d.o.subsystem,
			Name:      key,
		}, labelNames)
		d.gauge.Add(key, gauger, labelNames)
		d.registerer.MustRegister(gauger)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	gauger.WithLabelValues(labelValues...).Add(value)
}

func (d *driver) Histogram(ctx context.Context, handler metrics.EventHandler, key string, buckets []float64, value float64) {
	histogrammer, labelNames, ok := d.histogram.Get(key)

	labelValues := d.labelsPool.Get()
	defer d.labelsPool.Save(labelValues)

	if !ok {
		labelNames = d.labelsPool.Get()

		for k, v := range handler.Tags() {
			labelValues = append(labelValues, v)
			labelNames = append(labelNames, k)
		}

		histogrammer = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: d.o.namespace,
			Subsystem: d.o.subsystem,
			Name:      key,
			Buckets:   buckets,
		}, labelNames)
		d.histogram.Add(key, histogrammer, labelNames)
		d.registerer.MustRegister(histogrammer)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	histogrammer.WithLabelValues(labelValues...).Observe(value)
}

func (d *driver) Timing(ctx context.Context, handler metrics.EventHandler, key string, ms int) {
	d.Histogram(ctx, handler, key, nil, float64(ms))
}

func (d *driver) Duration(ctx context.Context, handler metrics.EventHandler, key string, v time.Duration) {
	d.Histogram(ctx, handler, key, nil, float64(v.Milliseconds()))
}

func (d *driver) Flush() {
}

func (d *driver) Close() {
}
