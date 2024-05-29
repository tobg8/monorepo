package ddtracing

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// AddTag is a helper to add a tag to a tracing span attached to the given context.
func AddTag(ctx context.Context, key string, value interface{}) {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		span.SetTag(key, value)
	}
}
