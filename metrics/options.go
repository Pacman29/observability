package metrics

type Option func(*options)

type options struct {
	defaultTags               map[string]string
	tagsMapPoolSaveCapacity   int
	tagsMapPoolCreateCapacity int
	ctxReaders                []CtxReader
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
