package zap

import (
	"context"
	"time"

	"go.uber.org/zap"
	"moul.io/http2curl"

	"github.com/Pacman29/observability/internal/pool"
	"github.com/Pacman29/observability/logger"
)

type driver struct {
	l       *zap.SugaredLogger
	pool    *pool.Slice[any]
	options *options
}

func NewZapDriver(l *zap.SugaredLogger, opts ...Option) logger.Driver {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &driver{
		options: o,
		l:       l,
		pool:    pool.NewSlice[any](o.saveCap, o.createCap, nil),
	}
}

func (d *driver) writeLog(ctx context.Context, writter func(msg string, kvs ...interface{}), h logger.EventHandler) {
	args := d.pool.Get()
	defer func() {
		d.pool.Save(args)
	}()
	args = d.toZapArgs(ctx, args, h)
	writter(h.Msg(), args...)
}

func (d *driver) Trace(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, d.l.Debugw, h)
}

func (d *driver) Debug(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, d.l.Debugw, h)
}

func (d *driver) Warning(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, d.l.Warnw, h)
}

func (d *driver) Info(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, d.l.Infow, h)
}

func (d *driver) Error(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, d.l.Errorw, h)
}

func (d *driver) Fatal(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, d.l.Fatalw, h)
}

func (d *driver) Flush(timeout time.Duration) error {
	return d.l.Sync()
}

func (d *driver) toZapArgs(ctx context.Context, args []any, h logger.EventHandler) []any {
	for k, v := range h.Tags() {
		args = append(args, zap.String(k, v))
	}
	for k, v := range h.Fields() {
		args = append(args, zap.Any(k, v))
	}
	if err := h.Err(); err != nil {
		args = append(args, zap.Error(err))
	}
	for _, v := range h.Args() {
		args = append(args, v)
	}
	if req := h.Req(); req != nil {
		if reqString, err := http2curl.GetCurlCommand(req); err != nil {
			d.l.Warn("can't convert request to curl", err)
		} else {
			args = append(args, zap.String("request", reqString.String()))
		}
	}
	if d.options.ctxArgsResolver != nil {
		args = append(args, d.options.ctxArgsResolver(ctx)...)
	}

	return args
}
