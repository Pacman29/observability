package metrics

import (
	"context"
	"maps"
)

type key int

const (
	tagKey key = iota
)

func getTags(ctx context.Context) map[string]string {
	if m, ok := ctx.Value(tagKey).(map[string]string); ok && m != nil {
		return m
	}
	return nil
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

func getCtxTags(ctx context.Context) map[string]string {
	m := getTags(ctx)
	if m == nil {
		return nil
	}
	nm := make(map[string]string, len(m))
	maps.Copy(nm, m)

	return nm
}
