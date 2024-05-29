package tracing

import (
	"context"

	"github.com/monorepo/common/monitoring/tracing/ddtracing"
)

// Deprecated: Use ddtracing.AddTag() instead.
func AddTag(ctx context.Context, key string, value interface{}) {
	ddtracing.AddTag(ctx, key, value)
}
