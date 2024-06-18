package interceptors

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"

	"github.com/f2prateek/train"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/monorepo/common/secret"
)

const tracingRequestBodyMaxContentLength = 500

// Tracing is the structure used to store interceptor metadata and context
type Tracing struct {
	sqp             secret.QueryParams
	withRequestBody bool
	tryIndentJSON   bool
}

// NewTracing instantiates a new tracing interceptor
func NewTracing() *Tracing {
	return &Tracing{
		sqp: []string{},
	}
}

// WithSecretQueryParams sets the names of the query params whose value won't be
// sent in trace metadata. E.g. for api keys sent as query params.
func (t *Tracing) WithSecretQueryParams(names ...string) *Tracing {
	t.sqp = names
	return t
}

// WithRequestBodyLog activate the logging of request body in the trace.
// If tryIndentJSON is true, logged body will be json indented.
// The full request body content will be sent to datadog.
// Use with caution and try to avoid sending GDPR or sensible data to datadog.
func (t *Tracing) WithRequestBodyLog(tryIndentJSON bool) *Tracing {
	t.withRequestBody = true
	t.tryIndentJSON = tryIndentJSON
	return t
}

// Intercept implements the train.Interceptor interface for the tracing
func (t *Tracing) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	span, ctx := tracer.StartSpanFromContext(req.Context(), "http.request",
		tracer.SpanType(ext.SpanTypeHTTP),
		tracer.ResourceName(req.URL.Host),
	)
	defer span.Finish()

	span.SetTag("http.method", req.Method)
	span.SetTag("http.scheme", req.URL.Scheme)
	span.SetTag("http.host", req.Host)
	span.SetTag("http.path", req.URL.Path)
	span.SetTag("http.queryparams", t.sqp.HideFromValues(req.URL.Query()).Encode())

	if !strings.Contains(req.Host, "svc.cluster.local") && !strings.Contains(req.Host, "svc.disco") {
		span.SetTag(ext.ServiceName, req.Host)
	}

	requestBodyTagValue := t.getRequestBodyTagValue(req)
	if len(requestBodyTagValue) > 0 {
		span.SetTag("http.request_body", requestBodyTagValue)
	}

	var dnsSpan tracer.Span
	var connectSpan tracer.Span
	var tlsSpan tracer.Span
	trace := &httptrace.ClientTrace{
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			dnsSpan, _ = tracer.StartSpanFromContext(ctx, "dns.resolution", tracer.SpanType("dns"))
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			dnsSpan.SetTag("dns.addrs", dnsInfo.Addrs)
			dnsSpan.Finish(tracer.WithError(dnsInfo.Err))
		},

		ConnectStart: func(network, addr string) {
			connectSpan, _ = tracer.StartSpanFromContext(ctx, "http.connect", tracer.SpanType("connect"))
		},
		ConnectDone: func(_, _ string, err error) {
			connectSpan.Finish(tracer.WithError(err))
		},

		TLSHandshakeStart: func() {
			tlsSpan, _ = tracer.StartSpanFromContext(ctx, "http.tls", tracer.SpanType("tls_handshake"))
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
			tlsSpan.Finish(tracer.WithError(err))
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(ctx, trace))

	err := tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
	if err != nil {
		span.SetTag("tracer_inject_error", err)
	}

	resp, err := chain.Proceed(req)
	if err != nil {
		if errContext := ctx.Err(); errContext != nil {
			span.SetTag(ext.Error, errContext)
		} else {
			span.SetTag(ext.Error, err)
		}
	} else {
		span.SetTag(ext.HTTPCode, resp.StatusCode)
		if resp.StatusCode >= 500 {
			span.SetTag(ext.Error, errors.New("status code 5XX"))
		}
	}

	return resp, err
}

func (t *Tracing) getRequestBodyTagValue(req *http.Request) string {
	if !t.withRequestBody || req.GetBody == nil || req.ContentLength <= 0 {
		return ""
	}
	if req.ContentLength > tracingRequestBodyMaxContentLength {
		return fmt.Sprintf("tracing info: contentlength is to big for tracing")
	}

	body, err := req.GetBody()
	if err != nil {
		return fmt.Sprintf("tracing error: req.GetBody: %s", err)
	}

	rawBody, err := io.ReadAll(body)
	if err != nil {
		return fmt.Sprintf("tracing error: io.ReadAll(body): %s", err)
	}

	bodyContent := string(rawBody)
	if !t.tryIndentJSON {
		return bodyContent
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(bodyContent), "", "  "); err != nil {
		return bodyContent
	}
	return prettyJSON.String()
}
