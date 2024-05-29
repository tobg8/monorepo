package semconv

import (
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// HTTPServerRequestDuration is the metric conforming to the
// "http.server.request.duration" semantic conventions.
//
// It represents the duration of HTTP server requests.
//
// Instrument: histogram
// Unit: s
// Stability: Stable
const (
	HTTPServerRequestDurationName        = semconv.HTTPServerRequestDurationName
	HTTPServerRequestDurationUnit        = semconv.HTTPServerRequestDurationUnit
	HTTPServerRequestDurationDescription = semconv.HTTPServerRequestDurationDescription
)

// HTTPServerActiveRequests is the metric conforming to the
// "http.server.active_requests" semantic conventions.
//
// It represents the number of active HTTP server requests.
//
// Instrument: updowncounter
// Unit: {request}
// Stability: Experimental
const (
	HTTPServerActiveRequestsName        = semconv.HTTPServerActiveRequestsName
	HTTPServerActiveRequestsUnit        = semconv.HTTPServerActiveRequestsUnit
	HTTPServerActiveRequestsDescription = semconv.HTTPServerActiveRequestsDescription
)

// HTTPServerRequestBodySize is the metric conforming to the
// "http.server.request.body.size" semantic conventions.
//
// It represents the size of HTTP server request bodies.
//
// Instrument: histogram
// Unit: By
// Stability: Experimental
const (
	HTTPServerRequestBodySizeName        = semconv.HTTPServerRequestBodySizeName
	HTTPServerRequestBodySizeUnit        = semconv.HTTPServerRequestBodySizeUnit
	HTTPServerRequestBodySizeDescription = semconv.HTTPServerRequestBodySizeDescription
)

// HTTPServerResponseBodySize is the metric conforming to the
// "http.server.response.body.size" semantic conventions.
//
// It represents the size of HTTP server response bodies.
//
// Instrument: histogram
// Unit: By
// Stability: Experimental
const (
	HTTPServerResponseBodySizeName        = semconv.HTTPServerResponseBodySizeName
	HTTPServerResponseBodySizeUnit        = semconv.HTTPServerResponseBodySizeUnit
	HTTPServerResponseBodySizeDescription = semconv.HTTPServerResponseBodySizeDescription
)
