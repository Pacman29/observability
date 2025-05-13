package sentry

import (
	"context"
	"errors"
	"maps"
	"os"
	"time"

	"github.com/getsentry/sentry-go"

	"observability/internal/pool"
	"observability/logger"
)

type driver struct {
	c          *sentry.Client
	options    *options
	tagsPool   *pool.Map[string, string]
	fieldsPool *pool.Map[string, any]
	argsPool   *pool.Slice[any]
}

func NewSentryDriver(client *sentry.Client, opts ...Option) logger.Driver {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &driver{
		options:    o,
		c:          client,
		tagsPool:   pool.NewMap[string, string](o.tagsPoolCapSave, o.tagsPoolCapCreate, nil),
		fieldsPool: pool.NewMap[string, any](o.fieldsPoolCapSave, o.fieldsPoolCapCreate, nil),
		argsPool:   pool.NewSlice[any](o.argsPoolCapSave, o.argsPoolCapCreate, nil),
	}
}

func (d *driver) captureException(ctx context.Context, h logger.EventHandler) {
	scope := sentry.NewScope()

	if d.options.sentryUserResolver != nil {
		scope.SetUser(d.options.sentryUserResolver(ctx))
	}

	if req := h.Req(); req != nil {
		scope.SetRequest(req)
	}

	tagsMap := d.tagsPool.Get()
	defer func() {
		d.tagsPool.Save(tagsMap)
	}()

	for k, v := range h.Tags() {
		tagsMap[k] = v
	}
	if d.options.additionalTagsResolver != nil {
		maps.Copy(tagsMap, d.options.additionalTagsResolver(ctx))
	}

	scope.SetTags(tagsMap)

	fieldsMap := d.fieldsPool.Get()
	defer func() {
		d.fieldsPool.Save(fieldsMap)
	}()

	for k, v := range h.Fields() {
		fieldsMap[k] = v
	}
	if d.options.additionalFieldsResolver != nil {
		maps.Copy(fieldsMap, d.options.additionalFieldsResolver(ctx))
	}

	args := d.argsPool.Get()
	defer func() {
		d.argsPool.Save(args)
	}()

	for _, v := range h.Args() {
		args = append(args, v)
	}
	fieldsMap["__additional_args"] = args

	scope.SetExtras(fieldsMap)
	d.c.CaptureException(h.Err(), nil, scope)
}

func (d *driver) Trace(ctx context.Context, h logger.EventHandler) {}

func (d *driver) Debug(ctx context.Context, h logger.EventHandler) {}

func (d *driver) Warning(ctx context.Context, h logger.EventHandler) {}

func (d *driver) Info(ctx context.Context, h logger.EventHandler) {}

func (d *driver) Error(ctx context.Context, h logger.EventHandler) {
	d.captureException(ctx, h)
}

func (d *driver) Fatal(ctx context.Context, h logger.EventHandler) {
	d.captureException(ctx, h)
	os.Exit(1)
}

func (d *driver) Flush(timeout time.Duration) error {
	if !d.c.Flush(timeout) {
		return errors.New("can't flush data")
	}
	return nil
}
