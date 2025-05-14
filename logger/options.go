package logger

import "maps"

type options struct {
	defaultFields               map[string]any
	defaultTags                 map[string]string
	ctxReaders                  []CtxReader
	fieldsMapPoolCreateCapacity int
	fieldsMapPoolSaveCapacity   int
	tagsMapPoolCreateCapacity   int
	tagsMapPoolSaveCapacity     int
	argsArrayPoolCreateCapacity int
	argsArrayPoolSaveCapacity   int
}

type Option func(*options)

func newOptions() *options {
	return &options{
		defaultFields: map[string]any{},
		defaultTags:   map[string]string{},
		ctxReaders:    nil,
	}
}

func WithDefaultFields(fields map[string]any) Option {
	return func(o *options) {
		maps.Copy(o.defaultFields, fields)
	}
}

func WithDefaultField(k string, v any) Option {
	return func(o *options) {
		o.defaultFields[k] = v
	}
}

func WithDefaultTags(tags map[string]string) Option {
	return func(o *options) {
		maps.Copy(o.defaultTags, tags)
	}
}

func WithDefaultTag(k string, v string) Option {
	return func(o *options) {
		o.defaultTags[k] = v
	}
}

func WithCtxReader(reader CtxReader) Option {
	return func(o *options) {
		o.ctxReaders = append(o.ctxReaders, reader)
	}
}

func WithFieldsMapPoolCreateCapacity(c int) Option {
	return func(o *options) {
		o.fieldsMapPoolCreateCapacity = c
	}
}

func WithFieldsMapPoolSaveCapacity(c int) Option {
	return func(o *options) {
		o.fieldsMapPoolSaveCapacity = c
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

func WithArgsArrayPoolCreateCapacity(c int) Option {
	return func(o *options) {
		o.argsArrayPoolCreateCapacity = c
	}
}

func WithArgsArrayPoolSaveCapacity(c int) Option {
	return func(o *options) {
		o.argsArrayPoolSaveCapacity = c
	}
}
