package tracing

import (
	oteltrace "go.opentelemetry.io/otel/trace"
)

// TracerProvider is an alias to the tracer provider of the OTEL API.
//
// This interface is used to centralize all OTEL libraries import, and
// maybe in the next future, to limit the scope of the interface.
type TracerProvider interface {
	oteltrace.TracerProvider
}
