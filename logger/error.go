package logger

import (
	"context"
	"errors"
)

type errorWithFields interface {
	LoggerFields() map[string]any
}

type errorWithTags interface {
	LoggerTags() map[string]string
}

type errWrapper struct {
	fields map[string]any
	tags   map[string]string

	err error
}

func (e *errWrapper) Error() string {
	return e.err.Error()
}

func (e *errWrapper) Unwrap() error {
	return e.err
}

// Cause is used for https://pkg.go.dev/github.com/pkg/errors#Cause
func (e *errWrapper) Cause() error {
	return e.err
}

func withError(ctx context.Context, err error) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if err == nil {
		return ctx
	}

	var (
		errFields = map[string]interface{}{}
		errTags   = map[string]string{}
	)

	var wrappedErr *errWrapper
	if errors.As(err, &wrappedErr) {
		// if we have errWrapper somewhere inside err, then we will extract its fields and tags
		errFields = wrappedErr.fields
		errTags = wrappedErr.tags
	}

	if wrappedErr, ok := err.(*errWrapper); ok {
		// if this err is already err wrapper, then we will unwrap it to reduce stack.
		// Fields and tags were extracted in the previous step
		err = wrappedErr.err
	}

	if errWithFields, ok := err.(errorWithFields); ok {
		for k, v := range errWithFields.LoggerFields() {
			errFields[k] = v
		}
	}

	if errWithTags, ok := err.(errorWithTags); ok {
		for k, v := range errWithTags.LoggerTags() {
			errTags[k] = v
		}
	}

	ctx = addFieldsToCtx(ctx, errFields)
	ctx = addTagsToCtx(ctx, errTags)

	return addErrorToCtx(ctx, err)
}

// wrapError оборачивает переданную ошибку err тегами и полями из ctx и возвращает новую ошибку,
// которую затем можно использовать в методах withField и подобных для логирования ее вместе с данными из контекста
func wrapError(ctx context.Context, err error) error {
	if err == nil {
		return err // maintain error type
	}

	if ctx == nil {
		return err
	}

	var (
		ctxFields = getFields(ctx)
		ctxTags   = getTags(ctx)
	)

	var wrappedErr *errWrapper
	if errors.As(err, &wrappedErr) {
		// if we have errWrapper somewhere inside err, then we will extract its fields and tags
		for name, value := range wrappedErr.fields {
			if _, ok := ctxFields[name]; !ok {
				ctxFields[name] = value
			}
		}

		for name, value := range wrappedErr.tags {
			if _, ok := ctxTags[name]; !ok {
				ctxTags[name] = value
			}
		}
	}

	if wrappedErr, ok := err.(*errWrapper); ok {
		// if this err is already err wrapper, then we will unwrap it to reduce stack.
		// Fields and tags were extracted in the previous step
		err = wrappedErr.err
	}

	return &errWrapper{
		fields: ctxFields,
		tags:   ctxTags,
		err:    err,
	}
}
