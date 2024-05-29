package metrics

import (
	"time"

	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
)

var (
	globalStatsdHandler StatsdHandler = NoopStatsdHandler
)

// SetGlobalStatsdHandler defines the global statsd handler.
func SetGlobalStatsdHandler(sh StatsdHandler) {
	globalStatsdHandler = sh
}

// GetGlobalStatsdHandler returns the global statsd handler.
//
// If no global statsd handler has been registered, a noop StatsdHandler
// implementation is returned.
func GetGlobalStatsdHandler() StatsdHandler {
	return globalStatsdHandler
}

// Gauge measures the value of a metric at a particular time.
//
// This is short for GetGlobalStatsdHandler().Gauge(..).
func Gauge(name string, value float64, tags []string, rate float64) {
	GetGlobalStatsdHandler().Gauge(name, value, tags, rate)
}

// Timing sends timing information, it is an alias for TimeInMilliseconds.
//
// This is short for GetGlobalStatsdHandler().Timing(..).
func Timing(name string, value time.Duration, tags []string, rate float64) {
	GetGlobalStatsdHandler().Timing(name, value, tags, rate)
}

// Count tracks how many times something happened per second.
//
// This is short for GetGlobalStatsdHandler().Count(..).
func Count(name string, value int64, tags []string, rate float64) {
	GetGlobalStatsdHandler().Count(name, value, tags, rate)
}

// Histogram tracks the statistical distribution of a set of values on each host.
//
// This is short for GetGlobalStatsdHandler().Histogram(..).
func Histogram(name string, value float64, tags []string, rate float64) {
	GetGlobalStatsdHandler().Histogram(name, value, tags, rate)
}

// Distribution tracks the statistical distribution of a set of values across
// your infrastructure.
//
// This is short for GetGlobalStatsdHandler().Distributtion(..).
func Distribution(name string, value float64, tags []string, rate float64) {
	GetGlobalStatsdHandler().Distribution(name, value, tags, rate)
}

// ServiceCheck sends the provided ServiceCheck.
//
// This is short for GetGlobalStatsdHandler().ServiceCheck(..).
func ServiceCheck(sc *StatsdServiceCheck) {
	GetGlobalStatsdHandler().ServiceCheck(sc)
}

// Event sends the provided Event.
//
// This is short for GetGlobalStatsdHandler().Event(..).
func Event(e *StatsdEvent) {
	GetGlobalStatsdHandler().Event(e)
}

// SetGlobalMeterProvider registers mp as the global OTEL meter provider.
func SetGlobalMeterProvider(mp MeterProvider) {
	otel.SetMeterProvider(mp)
}

// GetGlobalMeterProvider returns the registered global OTEL meter provider.
func GetGlobalMeterProvider() MeterProvider {
	return otel.GetMeterProvider()
}

// Meter returns a `metric.Meter` from the global OTEL meter provider.
//
// If this is called before a global MeterProvider is registered the returned
// Meter will be a No-op implementation of a Meter.
//
// This is short for GetGlobalMeterProvider().Meter(...).
func Meter(name string, opts ...otelmetric.MeterOption) otelmetric.Meter {
	return GetGlobalMeterProvider().Meter(name, opts...)
}
