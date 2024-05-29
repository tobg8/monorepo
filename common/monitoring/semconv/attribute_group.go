package semconv

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Semantic convention attributes in the HTTP namespace.
const (
	// HTTPRequestBodySizeKey is the attribute Key conforming to the
	// "http.request.body.size" semantic conventions.
	//
	// It represents the size of the request payload body in bytes.
	// This is the number of bytes transferred excluding headers and is often,
	// but not always, present as the
	// [Content-Length](https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length)
	// header. For requests using transport encoding, this should be the
	// compressed size.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3495
	HTTPRequestBodySizeKey = semconv.HTTPRequestBodySizeKey

	// HTTPRequestMethodKey is the attribute Key conforming to the
	// "http.request.method" semantic conventions.
	//
	// It represents the HTTP request method.
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'GET', 'POST', 'HEAD'
	HTTPRequestMethodKey = semconv.HTTPRequestMethodKey

	// HTTPRequestMethodOriginalKey is the attribute Key conforming to the
	// "http.request.method_original" semantic conventions.
	//
	// It represents the original HTTP method sent by the client in the request line.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'GeT', 'ACL', 'foo'
	HTTPRequestMethodOriginalKey = semconv.HTTPRequestMethodOriginalKey

	// HTTPRequestResendCountKey is the attribute Key conforming to the
	// "http.request.resend_count" semantic conventions.
	//
	// It represents the ordinal number of request resending attempt
	// (for any reason, including redirects).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 3
	HTTPRequestResendCountKey = semconv.HTTPRequestResendCountKey

	// HTTPResponseBodySizeKey is the attribute Key conforming to the
	// "http.response.body.size" semantic conventions.
	//
	// It represents the size of the response payload body in bytes.
	// This is the number of bytes transferred excluding headers and is often,
	// but not always, present as the
	// [Content-Length](https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length)
	// header. For requests using transport encoding, this should be the
	// compressed size.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: experimental
	// Examples: 3495
	HTTPResponseBodySizeKey = semconv.HTTPResponseBodySizeKey

	// HTTPResponseStatusCodeKey is the attribute Key conforming to the
	// "http.response.status_code" semantic conventions.
	//
	// It represents the [HTTP response status code](https://tools.ietf.org/html/rfc7231#section-6).
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 200
	HTTPResponseStatusCodeKey = semconv.HTTPResponseStatusCodeKey

	// HTTPRouteKey is the attribute Key conforming to the "http.route"
	// semantic conventions.
	//
	// It represents the matched route, that is, the path template in the format used
	// by the respective server framework.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '/users/:userID?', '{controller}/{action}/{id?}'
	HTTPRouteKey = semconv.HTTPRouteKey
)

// All available values for the `HTTPRequestMethodKey`.
var (
	HTTPRequestMethodConnect = semconv.HTTPRequestMethodConnect
	HTTPRequestMethodDelete  = semconv.HTTPRequestMethodDelete
	HTTPRequestMethodGet     = semconv.HTTPRequestMethodGet
	HTTPRequestMethodHead    = semconv.HTTPRequestMethodHead
	HTTPRequestMethodOptions = semconv.HTTPRequestMethodOptions
	HTTPRequestMethodPatch   = semconv.HTTPRequestMethodPatch
	HTTPRequestMethodPost    = semconv.HTTPRequestMethodPost
	HTTPRequestMethodPut     = semconv.HTTPRequestMethodPut
	HTTPRequestMethodTrace   = semconv.HTTPRequestMethodTrace
	HTTPRequestMethodOther   = semconv.HTTPRequestMethodOther
)

// HTTPRequestBodySize returns an attribute KeyValue conforming to the
// "http.request.body.size" semantic conventions.
//
// It represents the size of the request payload body in bytes.
// This is the number of bytes transferred excluding headers and is often,
// but not always, present as the
// [Content-Length](https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length)
// header. For requests using transport encoding, this should be the compressed
// size.
var HTTPRequestBodySize = semconv.HTTPRequestBodySize

// HTTPRequestMethodOriginal returns an attribute KeyValue conforming to the
// "http.request.method_original" semantic conventions.
//
// It represents the original HTTP method sent by the client in the request line.
var HTTPRequestMethodOriginal = semconv.HTTPRequestMethodOriginal

// HTTPRequestResendCount returns an attribute KeyValue conforming to the
// "http.request.resend_count" semantic conventions.
//
// It represents the ordinal number of request resending attempt (for any reason, including redirects).
var HTTPRequestResendCount = semconv.HTTPRequestResendCount

// HTTPResponseBodySize returns an attribute KeyValue conforming to the
// "http.response.body.size" semantic conventions.
//
// It represents the size of the response payload body in bytes.
// This is the number of bytes transferred excluding headers and is often,
// but not always, present as the
// [Content-Length](https://www.rfc-editor.org/rfc/rfc9110.html#field.content-length)
// header. For requests using transport encoding, this should be the compressed
// size.
var HTTPResponseBodySize = semconv.HTTPResponseBodySize

// HTTPResponseStatusCode returns an attribute KeyValue conforming to the
// "http.response.status_code" semantic conventions.
//
// It represents the [HTTP response status code](https://tools.ietf.org/html/rfc7231#section-6).
var HTTPResponseStatusCode = semconv.HTTPResponseStatusCode

// HTTPRoute returns an attribute KeyValue conforming to the "http.route"
// semantic conventions.
//
// It represents the matched route, that is, the path template in the format
// used by the respective server framework.
var HTTPRoute = semconv.HTTPRoute

// These attributes may be used for any network related operation.
const (
	// NetworkPeerAddressKey is the attribute Key conforming to the
	// "network.peer.address" semantic conventions.
	//
	// It represents the peer address of the network connection -
	// IP address or Unix domain socket name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '10.1.2.80', '/tmp/my.sock'
	NetworkPeerAddressKey = semconv.NetworkPeerAddressKey

	// NetworkPeerPortKey is the attribute Key conforming to the
	// "network.peer.port" semantic conventions.
	//
	// It represents the peer port number of the network connection.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 65123
	NetworkPeerPortKey = semconv.NetworkPeerPortKey

	// NetworkProtocolNameKey is the attribute Key conforming to the
	// "network.protocol.name" semantic conventions.
	//
	// It represents the [OSI application layer](https://osi-model.com/application-layer/)
	// or non-OSI equivalent.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'amqp', 'http', 'mqtt'
	NetworkProtocolNameKey = semconv.NetworkProtocolNameKey

	// NetworkProtocolVersionKey is the attribute Key conforming to the
	// "network.protocol.version" semantic conventions.
	//
	// It represents the version of the protocol specified in `network.protocol.name`.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: '3.1.1'
	NetworkProtocolVersionKey = semconv.NetworkProtocolVersionKey

	// NetworkTransportKey is the attribute Key conforming to the
	// "network.transport" semantic conventions.
	//
	// It represents the [OSI transport layer](https://osi-model.com/transport-layer/) or
	// [inter-process communication method](https://wikipedia.org/wiki/Inter-process_communication).
	//
	// Type: Enum
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'tcp', 'udp'
	NetworkTransportKey = semconv.NetworkTransportKey
)

// NetworkPeerAddress returns an attribute KeyValue conforming to the
// "network.peer.address" semantic conventions.
//
// It represents the peer address of the network connection -
// IP address or Unix domain socket name.
var NetworkPeerAddress = semconv.NetworkPeerAddress

// NetworkPeerPort returns an attribute KeyValue conforming to the
// "network.peer.port" semantic conventions.
//
// It represents the peer port number of the network connection.
var NetworkPeerPort = semconv.NetworkPeerPort

// NetworkProtocolName returns an attribute KeyValue conforming to the
// "network.protocol.name" semantic conventions.
//
// It represents the [OSI application layer](https://osi-model.com/application-layer/)
// or non-OSI equivalent.
var NetworkProtocolName = semconv.NetworkProtocolName

// NetworkProtocolVersion returns an attribute KeyValue conforming to the
// "network.protocol.version" semantic conventions.
//
// It represents the version of the protocol specified in `network.protocol.name`.
var NetworkProtocolVersion = semconv.NetworkProtocolVersion

// All available values for `NetworkTransportKey`.
var (
	NetworkTransportTCP  = NetworkTransportKey.String("tcp")
	NetworkTransportUDP  = NetworkTransportKey.String("udp")
	NetworkTransportPipe = NetworkTransportKey.String("pipe")
	NetworkTransportUnix = NetworkTransportKey.String("unix")
)

// These attributes may be used to describe the server in a connection-based
// network interaction where there is one side that initiates the connection
// (the client is the side that initiates the connection). This covers all TCP
// network interactions since TCP is connection-based and one side initiates
// the connection (an exception is made for peer-to-peer communication over TCP
// where the "user-facing" surface of the protocol / API doesn't expose a clear
// notion of client and server). This also covers UDP network interactions
// where one side initiates the interaction, e.g. QUIC (HTTP/3) and DNS.
const (
	// ServerAddressKey is the attribute Key conforming to the "server.address"
	// semantic conventions.
	//
	// It represents the server domain name if available without reverse DNS lookup;
	// otherwise, IP address or Unix domain socket name.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'example.com', '10.1.2.80', '/tmp/my.sock'
	ServerAddressKey = semconv.ServerAddressKey

	// ServerPortKey is the attribute Key conforming to the "server.port"
	// semantic conventions.
	//
	// It represents the server port number.
	//
	// Type: int
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 80, 8080, 443
	ServerPortKey = semconv.ServerPortKey
)

// ServerAddress returns an attribute KeyValue conforming to the "server.address" semantic conventions.
//
// It represents the server domain name if available without reverse DNS lookup;
// otherwise, IP address or Unix domain socket name.
var ServerAddress = semconv.ServerAddress

// ServerPort returns an attribute KeyValue conforming to the "server.port" semantic conventions.
//
// It represents the server port number.
var ServerPort = semconv.ServerPort

// Attributes describing URL.
const (
	// URLFullKey is the attribute Key conforming to the "url.full" semantic
	// conventions.
	//
	// It represents the absolute URL describing a network resource according to
	// [RFC3986](https://www.rfc-editor.org/rfc/rfc3986)
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'https://www.foo.bar/search?q=OpenTelemetry#SemConv',
	// '//localhost'
	URLFullKey = attribute.Key("url.full")

	// URLSchemeKey is the attribute Key conforming to the "url.scheme" semantic conventions.
	//
	// It represents the [URI scheme](https://www.rfc-editor.org/rfc/rfc3986#section-3.1) component
	// identifying the used protocol.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'https', 'ftp', 'telnet'
	URLSchemeKey = semconv.URLSchemeKey
)

// URLFull returns an attribute KeyValue conforming to the "url.full"
// semantic conventions.
//
// It represents the absolute URL describing a network resource according to
// [RFC3986](https://www.rfc-editor.org/rfc/rfc3986)
var URLFull = semconv.URLFull

// URLScheme returns an attribute KeyValue conforming to the "url.scheme" semantic conventions.
//
// It represents the [URI scheme](https://www.rfc-editor.org/rfc/rfc3986#section-3.1) component
// identifying the used protocol.
var URLScheme = semconv.URLScheme

// Describes user-agent attributes.
const (
	// UserAgentOriginalKey is the attribute Key conforming to the
	// "user_agent.original" semantic conventions.
	//
	// It represents the value of the
	// [HTTP User-Agent](https://www.rfc-editor.org/rfc/rfc9110.html#field.user-agent)
	// header sent by the client.
	//
	// Type: string
	// RequirementLevel: Optional
	// Stability: stable
	// Examples: 'CERN-LineMode/2.15 libwww/2.17b3', 'Mozilla/5.0 (iPhone; CPU
	// iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)
	// Version/14.1.2 Mobile/15E148 Safari/604.1'
	UserAgentOriginalKey = semconv.UserAgentOriginalKey
)

// UserAgentOriginal returns an attribute KeyValue conforming to the
// "user_agent.original" semantic conventions.
//
// It represents the value of the
// [HTTP User-Agent](https://www.rfc-editor.org/rfc/rfc9110.html#field.user-agent)
// header sent by the client.
var UserAgentOriginal = semconv.UserAgentOriginal
