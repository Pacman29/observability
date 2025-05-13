package zap

import "context"

type options struct {
	createCap       int
	saveCap         int
	ctxArgsResolver func(ctx context.Context) []any
}

type Option func(o *options)

func newOptions() *options {
	return &options{
		createCap:       10,
		saveCap:         20,
		ctxArgsResolver: nil,
	}
}

func WithPoolSaveCapacity(c int) Option {
	return func(o *options) {
		o.saveCap = c
	}
}

func WithPoolCreateCapacity(c int) Option {
	return func(o *options) {
		o.createCap = c
	}
}

func WithCtxArgsResolver(f func(ctx context.Context) []any) Option {
	return func(o *options) {
		o.ctxArgsResolver = f
	}
}
