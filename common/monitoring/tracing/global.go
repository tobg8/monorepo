package tracing

import (
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// SetGlobalTracerProvider registers mp as the global OTEL tracer provider.
func SetGlobalTracerProvider(tp TracerProvider) {
	otel.SetTracerProvider(tp)
}

// GetGlobalTracerProvider returns the registered global OTEL tracer provider.
func GetGlobalTracerProvider() TracerProvider {
	return otel.GetTracerProvider()
}

// Tracer creates a named tracer that implements `oteltrace.Tracer` interface.
// If the name is an empty string then provider uses default name.
//
// This is short for GetGlobalTracerProvider().Tracer(...).
func Tracer(name string, options ...oteltrace.TracerOption) oteltrace.Tracer {
	return GetGlobalTracerProvider().Tracer(name, options...)
}
