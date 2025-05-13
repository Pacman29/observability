package slog

import (
	"context"
	"log/slog"
	"os"
	"time"

	"moul.io/http2curl"

	"observability/internal/pool"
	"observability/logger"
)

type driver struct {
	pool    *pool.Slice[any]
	l       *slog.Logger
	options *options
}

func NewSlogDriver(logger *slog.Logger, opts ...Option) logger.Driver {
	o := newOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &driver{
		options: o,
		l:       logger,
		pool:    pool.NewSlice[any](o.saveCap, o.createCap, nil),
	}
}

func (d *driver) writeLog(ctx context.Context, level slog.Level, h logger.EventHandler) {
	args := d.pool.Get()
	defer func() {
		d.pool.Save(args)
	}()
	var msg string
	msg, args = d.toSlogArgs(ctx, args, h)
	d.l.Log(ctx, level, msg, args...)
}

func (d *driver) Trace(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, slog.LevelDebug, h)
}

func (d *driver) Debug(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, slog.LevelDebug, h)
}

func (d *driver) Warning(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, slog.LevelWarn, h)
}

func (d *driver) Info(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, slog.LevelInfo, h)
}

func (d *driver) Error(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, slog.LevelError, h)
}

func (d *driver) Fatal(ctx context.Context, h logger.EventHandler) {
	d.writeLog(ctx, slog.LevelError, h)
	os.Exit(1)
}

func (d *driver) Flush(timeout time.Duration) error {
	return nil
}

func (d *driver) toSlogArgs(ctx context.Context, args []any, h logger.EventHandler) (string, []any) {
	for k, v := range h.Tags() {
		args = append(args, slog.String(k, v))
	}
	for k, v := range h.Fields() {
		args = append(args, slog.Any(k, v))
	}
	if err := h.Err(); err != nil {
		args = append(args, slog.Any("error", err))
	}
	for _, v := range h.Args() {
		args = append(args, v)
	}
	if req := h.Req(); req != nil {
		if reqString, err := http2curl.GetCurlCommand(req); err != nil {
			d.l.Warn("can't convert request to curl", err)
		} else {
			args = append(args, slog.String("request", reqString.String()))
		}
	}
	if d.options.ctxArgsResolver != nil {
		args = append(args, d.options.ctxArgsResolver(ctx)...)
	}
	return h.Msg(), args
}
