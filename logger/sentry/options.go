package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
)

type options struct {
	sentryUserResolver       func(ctx context.Context) sentry.User
	additionalTagsResolver   func(ctx context.Context) map[string]string
	additionalFieldsResolver func(ctx context.Context) map[string]any
	tagsPoolCapSave          int
	tagsPoolCapCreate        int
	fieldsPoolCapSave        int
	fieldsPoolCapCreate      int
	argsPoolCapSave          int
	argsPoolCapCreate        int
}

type Option func(o *options)

func newOptions() *options {
	return &options{
		sentryUserResolver:       nil,
		additionalTagsResolver:   nil,
		additionalFieldsResolver: nil,
		tagsPoolCapSave:          20,
		tagsPoolCapCreate:        10,
		fieldsPoolCapSave:        20,
		fieldsPoolCapCreate:      10,
		argsPoolCapSave:          20,
		argsPoolCapCreate:        10,
	}
}

func WithSentryUserResolver(f func(ctx context.Context) sentry.User) Option {
	return func(o *options) {
		o.sentryUserResolver = f
	}
}

func WithAdditionalTagsResolver(f func(ctx context.Context) map[string]string) Option {
	return func(o *options) {
		o.additionalTagsResolver = f
	}
}

func WithAdditionalFieldsResolver(f func(ctx context.Context) map[string]any) Option {
	return func(o *options) {
		o.additionalFieldsResolver = f
	}
}

func WithTagsPoolCapSave(n int) Option {
	return func(o *options) {
		o.tagsPoolCapSave = n
	}
}

func WithTagsPoolCapCreate(n int) Option {
	return func(o *options) {
		o.tagsPoolCapCreate = n
	}
}

func WithFieldsPoolCapSave(n int) Option {
	return func(o *options) {
		o.fieldsPoolCapSave = n
	}
}

func WithFieldsPoolCapCreate(n int) Option {
	return func(o *options) {
		o.fieldsPoolCapCreate = n
	}
}

func WithArgsPoolCapCreate(n int) Option {
	return func(o *options) {
		o.argsPoolCapCreate = n
	}
}

func WithArgsPoolCapSave(n int) Option {
	return func(o *options) {
		o.argsPoolCapSave = n
	}
}
