package otelmetrics

import (
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	// Noop is a noop OTEL meter provider.
	Noop = noop.NewMeterProvider()
)
