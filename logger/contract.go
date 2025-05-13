package logger

import (
	"context"
	"iter"
	"net/http"
	"time"
)

type Driver interface {
	Trace(ctx context.Context, h EventHandler)
	Debug(ctx context.Context, h EventHandler)
	Warning(ctx context.Context, h EventHandler)
	Info(ctx context.Context, h EventHandler)
	Error(ctx context.Context, h EventHandler)
	Fatal(ctx context.Context, h EventHandler)
	Flush(timeout time.Duration) error
}

type CtxReader func(ctx context.Context) []any

type Logger interface {
	Trace(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warning(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Fatal(ctx context.Context, msg string, args ...any)
	WithField(ctx context.Context, k string, v any) context.Context
	WithFields(ctx context.Context, fields map[string]any) context.Context
	WithError(ctx context.Context, err error) context.Context
	Fields(ctx context.Context) map[string]any
	WithTag(ctx context.Context, k, v string) context.Context
	WithTags(ctx context.Context, tags map[string]string) context.Context
	Tags(ctx context.Context) map[string]string
	WithContext(ctx context.Context, src context.Context) context.Context
	WithRequest(ctx context.Context, request *http.Request) context.Context
	WrapError(ctx context.Context, err error) error
	Field(k string, v any) any
	Err(err error) any
	Tag(k string, v string) any
	Flush(timeout time.Duration) bool
}

type EventHandler interface {
	Msg() string
	Fields() iter.Seq2[string, any]
	Tags() iter.Seq2[string, string]
	Args() iter.Seq2[int, any]
	Err() error
	Req() *http.Request
}
