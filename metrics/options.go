package metrics

import "maps"

type Option func(*options)

type options struct {
	defaultTags                    map[string]string
	tagsMapPoolSaveCapacity        int
	tagsMapPoolCreateCapacity      int
	ctxReaders                     []CtxReader
	bucketsSlicePoolSaveCapacity   int
	bucketsSlicePoolCreateCapacity int
	defaultBuckets                 []float64
}

func newOptions() *options {
	return &options{
		defaultTags: map[string]string{},
		ctxReaders:  nil,
	}
}

func WithDefaultTags(tags map[string]string) Option {
	return func(o *options) {
		o.defaultTags = tags
	}
}

func WithTagsMapPoolCreateCapacity(c int) Option {
	return func(o *options) {
		o.tagsMapPoolCreateCapacity = c
	}
}

func WithTagsMapPoolSaveCapacity(c int) Option {
	return func(o *options) {
		o.tagsMapPoolSaveCapacity = c
	}
}

func WithCtxReader(reader CtxReader) Option {
	return func(o *options) {
		o.ctxReaders = append(o.ctxReaders, reader)
	}
}

type metricOptions struct {
	tags    map[string]string
	buckets []float64
}
type MetricOption func(option *metricOptions)

func newMetricOption() *metricOptions {
	return &metricOptions{
		tags: map[string]string{},
	}
}

func WithTag(k string, v string) MetricOption {
	return func(option *metricOptions) {
		option.tags[k] = v
	}
}

func WithTags(m map[string]string) MetricOption {
	return func(option *metricOptions) {
		maps.Copy(option.tags, m)
	}
}

func WithBuckets(buckets ...float64) MetricOption {
	return func(option *metricOptions) {
		option.buckets = buckets
	}
}
