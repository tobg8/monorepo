package tracing

import (
	oteltracenoop "go.opentelemetry.io/otel/trace/noop"
)

// NoopTracerProvider is a noop OTEL tracer provider.
var NoopTracerProvider TracerProvider = oteltracenoop.NewTracerProvider()
