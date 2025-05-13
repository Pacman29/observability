package multiple

import (
	"context"
	"os"
	"time"

	"go.uber.org/multierr"

	"github.com/Pacman29/observability/logger"
)

type drivers []logger.Driver

func NewMultiple(loggers ...logger.Driver) logger.Driver {
	var ds drivers
	ds = append(ds, loggers...)
	return ds
}

func (ds drivers) Trace(ctx context.Context, h logger.EventHandler) {
	for _, d := range ds {
		d.Trace(ctx, h)
	}
}

func (ds drivers) Debug(ctx context.Context, h logger.EventHandler) {
	for _, d := range ds {
		d.Debug(ctx, h)
	}
}

func (ds drivers) Warning(ctx context.Context, h logger.EventHandler) {
	for _, d := range ds {
		d.Warning(ctx, h)
	}
}

func (ds drivers) Info(ctx context.Context, h logger.EventHandler) {
	for _, d := range ds {
		d.Info(ctx, h)
	}
}

func (ds drivers) Error(ctx context.Context, h logger.EventHandler) {
	for _, d := range ds {
		d.Error(ctx, h)
	}
}

func (ds drivers) Fatal(ctx context.Context, h logger.EventHandler) {
	for _, d := range ds {
		d.Error(ctx, h)
	}
	os.Exit(1)
}

func (ds drivers) Flush(timeout time.Duration) error {
	var errs []error
	for _, d := range ds {
		if err := d.Flush(timeout); err != nil {
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}
