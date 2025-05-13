package logger

import (
	"context"
	"maps"
	"net/http"
)

type key int

const (
	fieldKey key = iota
	tagKey
	errorKey
	requestKey
)

func getFields(ctx context.Context) map[string]any {
	if m, ok := ctx.Value(fieldKey).(map[string]any); ok && m != nil {
		return m
	}
	return nil
}

func getTags(ctx context.Context) map[string]string {
	if m, ok := ctx.Value(tagKey).(map[string]string); ok && m != nil {
		return m
	}
	return nil
}

func getError(ctx context.Context) error {
	if e, ok := ctx.Value(errorKey).(error); ok && e != nil {
		return e
	}
	return nil
}

func getRequest(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(requestKey).(*http.Request); ok && r != nil {
		return r
	}
	return nil
}

func addRequestToCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, r)
}

func addFieldToCtx(ctx context.Context, k string, v any) context.Context {
	m := getFields(ctx)
	if len(m) == 0 {
		return context.WithValue(ctx, fieldKey, map[string]any{
			k: v,
		})
	}

	nm := make(map[string]any, len(m)+1)
	maps.Copy(nm, m)
	nm[k] = v
	return context.WithValue(ctx, fieldKey, nm)
}

func addFieldsToCtx(ctx context.Context, fs map[string]any) context.Context {
	m := getFields(ctx)
	if len(m) == 0 {
		nm := make(map[string]any, len(fs))
		maps.Copy(nm, fs)
		return context.WithValue(ctx, fieldKey, nm)
	}

	nm := make(map[string]any, len(m)+len(fs))
	maps.Copy(nm, m)
	maps.Copy(nm, fs)
	return context.WithValue(ctx, fieldKey, nm)
}

func addTagToCtx(ctx context.Context, k string, v string) context.Context {
	m := getTags(ctx)
	if len(m) == 0 {
		return context.WithValue(ctx, tagKey, map[string]string{
			k: v,
		})
	}

	nm := make(map[string]string, len(m)+1)
	maps.Copy(nm, m)
	nm[k] = v
	return context.WithValue(ctx, tagKey, nm)
}

func addTagsToCtx(ctx context.Context, tgs map[string]string) context.Context {
	m := getTags(ctx)
	if len(m) == 0 {
		nm := make(map[string]string, len(tgs))
		maps.Copy(nm, tgs)
		return context.WithValue(ctx, tagKey, nm)
	}

	nm := make(map[string]string, len(m)+len(tgs))
	maps.Copy(nm, m)
	maps.Copy(nm, tgs)
	return context.WithValue(ctx, tagKey, nm)
}

func addErrorToCtx(ctx context.Context, e error) context.Context {
	return context.WithValue(ctx, errorKey, e)
}

func getCtxFields(ctx context.Context) map[string]any {
	m := getFields(ctx)
	if m == nil {
		return nil
	}
	nm := make(map[string]any, len(m))
	maps.Copy(nm, m)

	return nm
}

func getCtxTags(ctx context.Context) map[string]string {
	m := getTags(ctx)
	if m == nil {
		return nil
	}
	nm := make(map[string]string, len(m))
	maps.Copy(nm, m)

	return nm
}

func copyCtx(dst context.Context, src context.Context) context.Context {
	if srcm := getFields(src); srcm != nil {
		dstm := getFields(dst)
		if dstm == nil {
			dstm = srcm
		} else {
			maps.Insert(dstm, maps.All(srcm))
		}
		dst = context.WithValue(dst, fieldKey, dstm)
	}

	if srcm := getTags(src); srcm != nil {
		dstm := getTags(dst)
		if dstm == nil {
			dstm = srcm
		} else {
			maps.Insert(dstm, maps.All(srcm))
		}
		dst = context.WithValue(dst, tagKey, dstm)
	}

	err := getError(src)
	if err != nil {
		dst = addErrorToCtx(src, err)
	}
	return dst
}
