package semconv

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// The software deployment.
const (
	// DeploymentEnvironmentKey is the attribute Key conforming to the
	// "deployment.environment" semantic conventions.
	//
	// It represents the name of the [deployment
	// environment](https://wikipedia.org/wiki/Deployment_environment) (aka
	// deployment tier).
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'staging', 'production'
	DeploymentEnvironmentKey = semconv.DeploymentEnvironmentKey
)

// DeploymentEnvironment returns an attribute KeyValue conforming to the
// "deployment.environment" semantic conventions.
//
// It represents the name of the
// [deployment environment](https://wikipedia.org/wiki/Deployment_environment)
// (aka deployment tier).
var DeploymentEnvironment = semconv.DeploymentEnvironment

// A service instance.
const (
	// ServiceNameKey is the attribute Key conforming to the "service.name"
	// semantic conventions.
	//
	// It represents the logical name of the service.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'shoppingcart'
	ServiceNameKey = semconv.ServiceNameKey

	// ServiceVersionKey is the attribute Key conforming to the
	// "service.version" semantic conventions.
	//
	// It represents the version string of the service API or implementation.
	// The format is not defined by these conventions.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: '2.0.0', 'a01dbef8a'
	ServiceVersionKey = semconv.ServiceVersionKey
)

// ServiceName returns an attribute KeyValue conforming to the
// "service.name" semantic conventions.
// It represents the logical name of the service.
var ServiceName = semconv.ServiceName

// ServiceVersion returns an attribute KeyValue conforming to the
// "service.version" semantic conventions.
// It represents the version string of the service API or implementation.
// The format is not defined by these conventions.
var ServiceVersion = semconv.ServiceVersion

// A service instance.
const (
	// ServiceNamespaceKey is the attribute Key conforming to the
	// "service.namespace" semantic conventions.
	//
	// It represents a namespace for `service.name`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 'Shop'
	ServiceNamespaceKey = semconv.ServiceNamespaceKey
)

// ServiceNamespace returns an attribute KeyValue conforming to the
// "service.namespace" semantic conventions.
//
// It represents a namespace for `service.name`.
var ServiceNamespace = semconv.ServiceNamespace

// The telemetry SDK used to capture data recorded by the instrumentation libraries.
const (
	// TelemetrySDKLanguageKey is the attribute Key conforming to the
	// "telemetry.sdk.language" semantic conventions.
	//
	// It represents the language of the telemetry SDK.
	//
	// Type: Enum
	// RequirementLevel: Required
	// Stability: experimental
	TelemetrySDKLanguageKey = semconv.TelemetrySDKLanguageKey

	// TelemetrySDKNameKey is the attribute Key conforming to the
	// "telemetry.sdk.name" semantic conventions.
	//
	// It represents the name of the telemetry SDK as defined above.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: 'opentelemetry'
	TelemetrySDKNameKey = semconv.TelemetrySDKNameKey

	// TelemetrySDKVersionKey is the attribute Key conforming to the
	// "telemetry.sdk.version" semantic conventions.
	//
	// It represents the version string of the telemetry SDK.
	//
	// Type: string
	// RequirementLevel: Required
	// Stability: experimental
	// Examples: '1.2.3'
	TelemetrySDKVersionKey = semconv.TelemetrySDKVersionKey
)

var (
	// TelemetrySDKLanguageGo is the Go Telemetry SDK language attribute.
	TelemetrySDKLanguageGo = semconv.TelemetrySDKLanguageGo
)

// TelemetrySDKName returns an attribute KeyValue conforming to the
// "telemetry.sdk.name" semantic conventions. It represents the name of the
// telemetry SDK as defined above.
func TelemetrySDKName(val string) attribute.KeyValue {
	return TelemetrySDKNameKey.String(val)
}

// TelemetrySDKVersion returns an attribute KeyValue conforming to the
// "telemetry.sdk.version" semantic conventions. It represents the version
// string of the telemetry SDK.
func TelemetrySDKVersion(val string) attribute.KeyValue {
	return TelemetrySDKVersionKey.String(val)
}
