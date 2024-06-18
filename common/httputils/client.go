package httputils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/monitoring/metrics"

	"github.com/monorepo/common/httputils/interceptors"
	"github.com/monorepo/common/httputils/svcauth"
	"github.com/monorepo/common/secret"
)

// Client is an http.Client wrapper.
type Client struct {
	*http.Client
	isTraced    bool
	isMonitored bool
	sqp         secret.QueryParams
}

// HTTPTransportWithInterceptors override http.Transport type to support the registration of
// train.Interceptors middleware.
// See Client.appendInterceptors() and Client.prependInterceptors() methods.
type HTTPTransportWithInterceptors struct {
	http.Transport
	interceptors []train.Interceptor
}

// NewHTTPTransportWithInterceptors returns a *HTTPTransportWithInterceptors.
// DefaultTransport properties values are reused except for MaxIdleConns and MaxIdleConnsPerHosts.
func NewHTTPTransportWithInterceptors(timeout, keepalive time.Duration) *HTTPTransportWithInterceptors {
	// Dialer with customized timeout
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: keepalive,
	}

	// Make a copy of http.DefaultTransport to use its values,
	// and be able to customize it without changing the stdlib's default
	dt := http.DefaultTransport.(*http.Transport)

	return &HTTPTransportWithInterceptors{
		Transport: http.Transport{
			Proxy: dt.Proxy,

			// As we often have a lot of traffic on a very few number of host
			// targets and the default value of http.DefaultMaxIdleConnsPerHost
			// is very low (2) we often close/open a great number TCP session
			// instead of keeping them idle and reuse them.
			MaxIdleConns:        300,
			MaxIdleConnsPerHost: 200,

			IdleConnTimeout:       dt.IdleConnTimeout,
			TLSHandshakeTimeout:   dt.TLSHandshakeTimeout,
			ExpectContinueTimeout: dt.ExpectContinueTimeout,
			DialContext:           dialer.DialContext,
		},
	}
}

// RoundTrip is the main routine called by the
// http module to execute the given request.
func (t *HTTPTransportWithInterceptors) RoundTrip(req *http.Request) (*http.Response, error) {
	return train.TransportWith(&t.Transport, t.interceptors...).RoundTrip(req)
}

// AppendInterceptors register new train.Interceptors middlewares.
// Add the given http middlewares at the end of the list.
func (t *HTTPTransportWithInterceptors) AppendInterceptors(is ...train.Interceptor) {
	t.interceptors = append(t.interceptors, is...)
}

// PrependInterceptors register new train.Interceptors middlewares.
// Add the given http middlewares at the begin of the list.
func (t *HTTPTransportWithInterceptors) PrependInterceptors(is ...train.Interceptor) {
	t.interceptors = append(is, t.interceptors...)
}

// NewClient returns a *Client with the specified timeout and keepalive.
// If keepalive is 0, it is disabled.
func NewClient(timeout, keepalive time.Duration) *Client {
	tr := NewHTTPTransportWithInterceptors(timeout, keepalive)
	tr.AppendInterceptors(
	//interceptors.NewUniqueID(),
	//interceptors.NewBrandForwarding(),
	)
	return &Client{
		Client: &http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
		isMonitored: false,
		isTraced:    false,
	}
}

// appendInterceptors allows to register new interceptors to the Client
func (client *Client) appendInterceptors(is ...train.Interceptor) {
	tr := client.Transport.(*HTTPTransportWithInterceptors)
	tr.AppendInterceptors(is...)
}

// prependInterceptors allows to register new interceptors to the Client
func (client *Client) prependInterceptors(is ...train.Interceptor) {
	tr := client.Transport.(*HTTPTransportWithInterceptors)
	tr.PrependInterceptors(is...)
}

// WithLimiter limits the number of concurrent requests.
func (client *Client) WithLimiter(count int) *Client {
	client.appendInterceptors(interceptors.NewLimiter(count))
	return client
}

// WithAuthorizationHeader add the introspection token to the request's Authorization header.
func (client *Client) WithAuthorizationHeader() *Client {
	client.appendInterceptors(interceptors.NewAuthorization())
	return client
}

// WithCircuitBreaker set a circuit Breaker with monitoring.
// backPressureThreshold is the number of possible consecutive failure
// duration of the open circuit
//func (client *Client) WithCircuitBreaker(name string, backPressureThreshold uint64, duration time.Duration) *Client {
//	client.appendInterceptors(interceptors.NewBreaker(name, backPressureThreshold, duration))
//	return client
//}

// WithConsentChecker propagates user consents through HTTP using the Didomi cookie
// if provided in the request context
//func (client *Client) WithConsentChecker(logger logging.Logger) *Client {
//	client.appendInterceptors(interceptors.NewConsentChecker(logger))
//	return client
//}

// Observe activates the monitoring and tracing of this client
func (client *Client) Observe(rp ...interceptors.RouteMatcher) *Client {
	panicIfAlreadySet(client, "both")
	sh := metrics.GetGlobalStatsdHandler()
	client.appendInterceptors(interceptors.NewMonitoring(sh, rp...))
	client.prependInterceptors(interceptors.
		NewTracing().
		WithSecretQueryParams(client.sqp...))
	return client
}

// WithMonitor activates the monitoring of the request associated with this client.
func (client *Client) WithMonitor(sh metrics.StatsdHandler, rp ...interceptors.RouteMatcher) *Client {
	panicIfAlreadySet(client, "monitor")
	client.appendInterceptors(interceptors.NewMonitoring(sh, rp...))
	return client
}

// WithTracer activates the tracer for the requests done with this client.
func (client *Client) WithTracer() *Client {
	panicIfAlreadySet(client, "tracer")
	client.prependInterceptors(interceptors.
		NewTracing().
		WithSecretQueryParams(client.sqp...))
	return client
}

// WithUserAgent define the user agent of the HTTP client.
func (client *Client) WithUserAgent(name, version string) *Client {
	client.appendInterceptors(interceptors.NewUserAgent(name, version))
	return client
}

// WithDialContext specifies the dial function for creating TCP connections.
func (client *Client) WithDialContext(dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) *Client {
	tr := client.Transport.(*HTTPTransportWithInterceptors)
	tr.DialContext = dialContext
	return client
}

// WithMaxIdleConnsPerHost defined the value of MaxIdleConnsPerHost of
// the HTTP transport.
func (client *Client) WithMaxIdleConnsPerHost(value int) *Client {
	tr := client.Transport.(*HTTPTransportWithInterceptors)
	tr.MaxIdleConnsPerHost = value
	return client
}

// WithServiceAuth defines the authorization interceptor.
func (client *Client) WithServiceAuth(conf svcauth.Conf) *Client {
	client.appendInterceptors(svcauth.NewServiceAuth(conf))
	return client
}

// WithSecretQueryParams set query params names whose values will be hidden in
// error messages and tracing metadata, e.g. api keys sent as query params.
//
// Note: http.NewRequest also displays full URL in errors and no secret query
// params hiding is provided for this in the present package. So
// secret.QueryParams.HideFromErr should be used on http.NewRequest errors in
// addition to this.
func (client *Client) WithSecretQueryParams(names ...string) *Client {
	if client.isTraced {
		panic("httputils Client.WithSecretQueryParams should be set before tracer to avoid leak through tracing")
	}
	client.sqp = names
	client.prependInterceptors(&interceptors.QueryObfuscator{QueryParams: client.sqp})

	return client
}

// Do executes a request and writes the response body in the provided io.Writer.
// That way, the user doesn't have to worry about closing the response's body anymore.
// If the body must be ignored, then pass io.Discard.
// The function returns the HTTP status code.
func (client *Client) Do(ctx context.Context, dst io.Writer, request *http.Request) (statusCode int, err error) {
	if dst == nil {
		panic(`destination cannot be nil, as this might panic when the response body is not empty; if you want to ignore the body, then pass io.Discard instead`)
	}

	response, err := client.Client.Do(request.WithContext(ctx))
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	_, err = io.Copy(dst, response.Body)
	if err != nil {
		return 0, fmt.Errorf("copy failed: %w", err)
	}

	return response.StatusCode, nil
}

// DoAndUnmarshalJSON executes a request and unmarshall the response body in the provided pointer
// The function returns the HTTP status code.
func (client *Client) DoAndUnmarshalJSON(ctx context.Context, v interface{}, request *http.Request) (statusCode int, err error) {
	response, err := client.Client.Do(request.WithContext(ctx))
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode >= http.StatusBadRequest {
		var buf bytes.Buffer
		r := io.TeeReader(response.Body, &buf)
		err = json.NewDecoder(r).Decode(v)
		if err != nil {
			return response.StatusCode, fmt.Errorf("decode status %d: body %q", response.StatusCode, buf.String())
		}
		return response.StatusCode, nil
	}

	err = json.NewDecoder(response.Body).Decode(v)
	if err != nil {
		return response.StatusCode, fmt.Errorf("decode status %d: %w", response.StatusCode, err)
	}
	return response.StatusCode, nil
}

func panicIfAlreadySet(c *Client, kind string) {
	switch kind {
	case "monitor":
		if c.isMonitored {
			panic(errors.New("trying to duplicate monitor on http client"))
		}
		c.isMonitored = true
	case "tracer":
		if c.isTraced {
			panic(errors.New("trying to duplicate tracer on http client"))
		}
		c.isTraced = true
	case "both":
		if c.isMonitored || c.isTraced {
			panic(errors.New("trying to duplicate observability on http client"))
		}
		c.isMonitored, c.isTraced = true, true
	default:
		panic(errors.New("unknown observability"))
	}
}
