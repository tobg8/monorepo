package monitoring

import (
	"context"
)

type ctxKey uint8

const (
	monitoringTagCtxKey ctxKey = iota
)

// AddTagsInContext returns a copy of the given context, with the
// given list of monitoring tags inside.
func AddTagsInContext(ctx context.Context, tags map[string]string) context.Context {
	m := GetTagsFromContext(ctx)
	for key, val := range tags {
		m[key] = val
	}

	return context.WithValue(ctx, monitoringTagCtxKey, m)
}

// GetTagsFromContext returns the list of monitoring tags of the given
// context.
func GetTagsFromContext(ctx context.Context) map[string]string {
	if tags, ok := ctx.Value(monitoringTagCtxKey).(map[string]string); ok {
		return tags
	}

	return map[string]string{}
}
