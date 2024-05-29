package metrics

import (
	"time"

	otelnoop "go.opentelemetry.io/otel/metric/noop"
)

// NoopStatsdHandler is a noop datadog statsd handler.
var NoopStatsdHandler = &noopStatsdHandler{}

type noopStatsdHandler struct{}

func (*noopStatsdHandler) Gauge(name string, value float64, tags []string, rate float64)        {}
func (*noopStatsdHandler) Timing(name string, value time.Duration, tags []string, rate float64) {}
func (*noopStatsdHandler) Count(name string, value int64, tags []string, rate float64)          {}
func (*noopStatsdHandler) Histogram(name string, value float64, tags []string, rate float64)    {}
func (*noopStatsdHandler) Distribution(name string, value float64, tags []string, rate float64) {}
func (*noopStatsdHandler) ServiceCheck(sc *StatsdServiceCheck)                                  {}
func (*noopStatsdHandler) Event(evt *StatsdEvent)                                               {}
func (*noopStatsdHandler) GetNamespace() string                                                 { return "" }

// NoopMeterProvider is a noop OTEL meter provider.
var NoopMeterProvider = otelnoop.NewMeterProvider()
