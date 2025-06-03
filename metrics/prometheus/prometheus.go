package prometheus

import (
	"context"

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

func (d *driver) Counter(ctx context.Context, handler metrics.EventHandler) {
	counter, labelNames, ok := d.counters.Get(handler.GetKey())

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
			Name:      handler.GetKey(),
		}, labelNames)
		d.counters.Add(handler.GetKey(), counter, labelNames)
		d.registerer.MustRegister(counter)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	counter.WithLabelValues(labelValues...).Add(handler.GetValue())
}

func (d *driver) Increment(ctx context.Context, handler metrics.EventHandler) {
	counter, labelNames, ok := d.counters.Get(handler.GetKey())

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
			Name:      handler.GetKey(),
		}, labelNames)
		d.counters.Add(handler.GetKey(), counter, labelNames)
		d.registerer.MustRegister(counter)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	counter.WithLabelValues(labelValues...).Add(handler.GetValue())
}

func (d *driver) Gauge(ctx context.Context, handler metrics.EventHandler) {
	gauger, labelNames, ok := d.gauge.Get(handler.GetKey())

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
			Name:      handler.GetKey(),
		}, labelNames)
		d.gauge.Add(handler.GetKey(), gauger, labelNames)
		d.registerer.MustRegister(gauger)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	gauger.WithLabelValues(labelValues...).Add(handler.GetValue())
}

func (d *driver) Histogram(ctx context.Context, handler metrics.EventHandler) {
	histogrammer, labelNames, ok := d.histogram.Get(handler.GetKey())

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
			Name:      handler.GetKey(),
			Buckets:   handler.GetBuckets(),
		}, labelNames)
		d.histogram.Add(handler.GetKey(), histogrammer, labelNames)
		d.registerer.MustRegister(histogrammer)
		d.labelsPool.Save(labelNames)
	} else {
		m := handler.GetTags()
		for _, v := range labelNames {
			labelValues = append(labelValues, m[v])
		}
	}

	histogrammer.WithLabelValues(labelValues...).Observe(handler.GetValue())
}

func (d *driver) Timing(ctx context.Context, handler metrics.EventHandler) {
	d.Histogram(ctx, handler)
}

func (d *driver) Duration(ctx context.Context, handler metrics.EventHandler) {
	d.Histogram(ctx, handler)
}

func (d *driver) Flush() {
}

func (d *driver) Close() {
}
