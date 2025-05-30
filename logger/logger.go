package logger

import (
	"context"
	"errors"
	"iter"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/Pacman29/observability/internal/pool"
)

type logger struct {
	d          Driver
	fieldsPool *pool.Map[string, any]
	tagsPool   *pool.Map[string, string]
	argsPool   *pool.Slice[any]
	o          *options
}

func New(d Driver, opts ...Option) Logger {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &logger{
		d:          d,
		fieldsPool: pool.NewMap[string, any](o.fieldsMapPoolSaveCapacity, o.fieldsMapPoolCreateCapacity, o.defaultFields),
		tagsPool:   pool.NewMap[string, string](o.tagsMapPoolSaveCapacity, o.tagsMapPoolCreateCapacity, o.defaultTags),
		argsPool:   pool.NewSlice[any](o.argsArrayPoolSaveCapacity, o.argsArrayPoolCreateCapacity, nil),
		o:          o,
	}
}

type logEventHandler struct {
	msg    string
	fields map[string]any
	tags   map[string]string
	args   []any
	err    error
	req    *http.Request
}

func (l *logger) newHandler(msg string) (*logEventHandler, func()) {
	h := &logEventHandler{
		msg:    msg,
		fields: l.fieldsPool.Get(),
		tags:   l.tagsPool.Get(),
		args:   l.argsPool.Get(),
		err:    nil,
	}
	return h, func() {
		l.fieldsPool.Save(h.fields)
		l.tagsPool.Save(h.tags)
		l.argsPool.Save(h.args)
	}
}

func (h *logEventHandler) resolveArgs(args ...any) {
	for _, arg := range args {
		switch argument := arg.(type) {
		case *field:
			h.fields[argument.k] = argument.v
		case *tag:
			h.tags[argument.k] = argument.v
		case error:
			h.err = argument
		default:
			h.args = append(h.args, arg)
		}
	}
}

func (h *logEventHandler) Fields() iter.Seq2[string, any] {
	return maps.All(h.fields)
}

func (h *logEventHandler) Args() iter.Seq2[int, any] {
	return slices.All(h.args)
}

func (h *logEventHandler) Tags() iter.Seq2[string, string] {
	return maps.All(h.tags)
}

func (h *logEventHandler) Err() error {
	return h.err
}

func (h *logEventHandler) Msg() string {
	return h.msg
}

func (h *logEventHandler) Req() *http.Request {
	return h.req
}

func (l *logger) withArgs(ctx context.Context, handler *logEventHandler, args ...any) {
	// добавляем данные из ридеров
	for _, reader := range l.o.ctxReaders {
		handler.resolveArgs(reader(ctx)...)
	}

	// потом все из контекста
	for k, v := range getFields(ctx) {
		handler.fields[k] = v
	}
	for k, v := range getTags(ctx) {
		handler.tags[k] = v
	}
	handler.err = getError(ctx)
	// потом все из ошибки
	if handler.err != nil {
		var wErr *errWrapper
		if errors.As(handler.err, &wErr) {
			for k, v := range wErr.fields {
				handler.fields[k] = v
			}
			for k, v := range wErr.tags {
				handler.tags[k] = v
			}
		}
	}
	handler.req = getRequest(ctx)

	// вытаскиваем все из аргументов текущих
	handler.resolveArgs(args...)
}

func (l *logger) Trace(ctx context.Context, msg string, args ...any) {
	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler(msg)
	defer handlerClose()
	l.withArgs(ctx, h, args...)

	l.d.Trace(ctx, h)
}

func (l *logger) Debug(ctx context.Context, msg string, args ...any) {
	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler(msg)
	defer handlerClose()
	l.withArgs(ctx, h, args...)

	l.d.Debug(ctx, h)
}

func (l *logger) Info(ctx context.Context, msg string, args ...any) {
	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler(msg)
	defer handlerClose()
	l.withArgs(ctx, h, args...)

	l.d.Info(ctx, h)
}

func (l *logger) Warning(ctx context.Context, msg string, args ...any) {
	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler(msg)
	defer handlerClose()
	l.withArgs(ctx, h, args...)

	l.d.Warning(ctx, h)
}

func (l *logger) Error(ctx context.Context, msg string, args ...any) {
	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler(msg)
	defer handlerClose()
	l.withArgs(ctx, h, args...)

	l.d.Error(ctx, h)
}

func (l *logger) Fatal(ctx context.Context, msg string, args ...any) {
	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler(msg)
	defer handlerClose()
	l.withArgs(ctx, h, args...)

	l.d.Fatal(ctx, h)
}

func (l *logger) Recover(ctx context.Context) {
	err := recover()
	if err == nil {
		return
	}

	ctx = defaultCtx(ctx)
	h, handlerClose := l.newHandler("")
	defer handlerClose()
	l.withArgs(ctx, h)

	l.d.Recover(err, ctx, h)
}

func (l *logger) WithField(ctx context.Context, k string, v any) context.Context {
	ctx = defaultCtx(ctx)
	return addFieldToCtx(ctx, k, v)
}

func (l *logger) WithFields(ctx context.Context, fields map[string]any) context.Context {
	ctx = defaultCtx(ctx)
	return addFieldsToCtx(ctx, fields)
}

func (l *logger) WithError(ctx context.Context, err error) context.Context {
	ctx = defaultCtx(ctx)
	return withError(ctx, err)
}

func (l *logger) Fields(ctx context.Context) map[string]any {
	ctx = defaultCtx(ctx)
	return getCtxFields(ctx)
}

func (l *logger) WithTag(ctx context.Context, k, v string) context.Context {
	ctx = defaultCtx(ctx)
	return addTagToCtx(ctx, k, v)
}

func (l *logger) WithTags(ctx context.Context, tags map[string]string) context.Context {
	ctx = defaultCtx(ctx)
	return addTagsToCtx(ctx, tags)
}

func (l *logger) Tags(ctx context.Context) map[string]string {
	ctx = defaultCtx(ctx)
	return getCtxTags(ctx)
}

func (l *logger) WithContext(ctx context.Context, src context.Context) context.Context {
	ctx = defaultCtx(ctx)
	return copyCtx(ctx, src)
}

func (l *logger) WithRequest(ctx context.Context, request *http.Request) context.Context {
	ctx = defaultCtx(ctx)
	return addRequestToCtx(ctx, request)
}

func (l *logger) WrapError(ctx context.Context, err error) error {
	ctx = defaultCtx(ctx)
	return wrapError(ctx, err)
}

func (l *logger) Field(k string, v any) any {
	return &field{k: k, v: v}
}

func (l *logger) Err(err error) any {
	return err
}

func (l *logger) Tag(k string, v string) any {
	return &tag{k: k, v: v}
}

func (l *logger) Flush(timeout time.Duration) bool {
	if err := l.d.Flush(timeout); err != nil {
		return false
	}

	return true
}

func defaultCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
