package httputils

import (
	"net/http"
	"net/http/httputil"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/httputils/interceptors"
	"github.com/monorepo/common/monitoring/metrics"
)

// ReverseProxy is an httputil.ReverseProxy wrapper.
type ReverseProxy struct {
	httputil.ReverseProxy
	isTraced    bool
	isMonitored bool
}

type transportWithInterceptors interface {
	RoundTrip(*http.Request) (*http.Response, error)
	AppendInterceptors(is ...train.Interceptor)
	PrependInterceptors(is ...train.Interceptor)
}

// NewReverseProxy returns a *ReverseProxy with the specified timeout and keepalive.
// If keepalive is 0, it is disabled.
func NewReverseProxy(director func(*http.Request), transport transportWithInterceptors) *ReverseProxy {
	return &ReverseProxy{
		ReverseProxy: httputil.ReverseProxy{
			Director:  director,
			Transport: transport,
		},
	}
}

// appendInterceptors allows to register new interceptors to the Client
func (rp *ReverseProxy) appendInterceptors(is ...train.Interceptor) {
	tr := rp.Transport.(transportWithInterceptors)
	tr.AppendInterceptors(is...)
}

// prependInterceptors allows to register new interceptors to the Client
func (rp *ReverseProxy) prependInterceptors(is ...train.Interceptor) {
	tr := rp.Transport.(transportWithInterceptors)
	tr.PrependInterceptors(is...)
}

// Observe activates the monitoring and tracing of this client
// The metrics will be identified by its name prefixed by the defaultnamespace
func (rp *ReverseProxy) Observe(sh metrics.StatsdHandler, rm ...interceptors.RouteMatcher) *ReverseProxy {
	return rp.WithMonitor(sh, rm...).WithTracer()
}

// WithMonitor activates the monitoring of the request associated with this client.
// The metrics will be identified by its name given as parameter.
func (rp *ReverseProxy) WithMonitor(sh metrics.StatsdHandler, rm ...interceptors.RouteMatcher) *ReverseProxy {
	if rp.isMonitored {
		return rp
	}
	rp.isMonitored = true
	rp.appendInterceptors(interceptors.NewMonitoring(sh, rm...))
	return rp
}

// WithTracer activates the tracer for the requests done with this client.
func (rp *ReverseProxy) WithTracer() *ReverseProxy {
	if rp.isTraced {
		return rp
	}
	rp.isTraced = true
	rp.prependInterceptors(interceptors.NewTracing())
	return rp
}
