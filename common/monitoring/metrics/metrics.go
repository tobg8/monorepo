package metrics

import (
	"time"

	"go.opentelemetry.io/otel/metric"
)

// StatsdHandler allows monitoring in a way slightly simplified compared to datadog.
type StatsdHandler interface {
	Gauge(name string, value float64, tags []string, rate float64)
	Timing(name string, value time.Duration, tags []string, rate float64)
	Count(name string, value int64, tags []string, rate float64)
	Histogram(name string, value float64, tags []string, rate float64)
	Distribution(name string, value float64, tags []string, rate float64)
	ServiceCheck(sc *StatsdServiceCheck)
	Event(e *StatsdEvent)

	// Deprecated: This method should not be called.
	//
	// The namespace is global to the application, and should be defined once
	// by the common/application library. No library should be dependent
	// on this namespace.
	GetNamespace() string
}

// An StatsdEvent is an object that can be posted to your DataDog event stream.
type StatsdEvent struct {
	Title     string
	Text      string
	Timestamp time.Time
	Hostname  string
	AlertType StatsdEventAlertType
	Tags      []string
}

// StatsdEventAlertType is the alert type for events
type StatsdEventAlertType uint8

// All available values for StatsdEventAlertType.
const (
	AlertTypeInfo StatsdEventAlertType = iota + 1
	AlertTypeError
	AlertTypeWarning
	AlertTypeSuccess
)

// A StatsdServiceCheck is an object that contains status of DataDog service check.
type StatsdServiceCheck struct {
	Name      string
	Status    StatsdServiceCheckStatus
	Timestamp time.Time
	Hostname  string
	Message   string
	Tags      []string
}

// StatsdServiceCheckStatus is the status of the ServiceCheck.
type StatsdServiceCheckStatus uint8

// All available values for ServiceCheckStatus.
const (
	ServiceCheckStatusOk StatsdServiceCheckStatus = iota + 1
	ServiceCheckStatusWarn
	ServiceCheckStatusCritical
)

// MeterProvider is an alias to the meter provider of the OTEL API.
//
// This interface is used to centralize all OTEL libraries import, and
// maybe in the next future, to limit the scope of the interface.
type MeterProvider interface {
	metric.MeterProvider
}
