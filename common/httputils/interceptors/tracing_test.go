package interceptors

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/f2prateek/train"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/monorepo/common/pointer"
)

func TestTracing_Intercept(t *testing.T) {
	client := new(http.Client)

	// Start the mock tracer.
	mt := mocktracer.Start()
	defer mt.Stop()

	client.Transport = train.Transport(NewTracing())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.NotEmpty(t, req.Header.Get(tracer.DefaultParentIDHeader))
		require.NotEmpty(t, req.Header.Get(tracer.DefaultTraceIDHeader))
	}))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)

	_, ctx := tracer.StartSpanFromContext(context.Background(), "test")
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.NoError(t, err)

	// Query the mock tracer for finished spans.
	spans := mt.FinishedSpans()
	require.Len(t, spans, 2)
}

func TestTracing_Intercept_with_a_custom_resourceNameFunc(t *testing.T) {
	client := new(http.Client)

	// Start the mock tracer.
	mt := mocktracer.Start()
	defer mt.Stop()

	client.Transport = train.Transport(NewTracing())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.NotEmpty(t, req.Header.Get(tracer.DefaultParentIDHeader))
		require.NotEmpty(t, req.Header.Get(tracer.DefaultTraceIDHeader))
	}))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)

	_, ctx := tracer.StartSpanFromContext(context.Background(), "test")
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.NoError(t, err)

	// Query the mock tracer for finished spans.
	spans := mt.FinishedSpans()
	require.Len(t, spans, 2)
}

func TestTracing_Intercept_with_a_call_to_an_internal_service(t *testing.T) {
	client := new(http.Client)

	// Start the mock tracer.
	mt := mocktracer.Start()
	defer mt.Stop()

	client.Transport = train.Transport(NewTracing())

	scheme := "https"
	host := "some-service.some-env.k8s.leboncoin.lan"
	path := "/messages/me"
	queryparams := "since=2024-06-01&withUserID=1234"
	url := fmt.Sprintf("%s://%s%s?%s", scheme, host, path, queryparams)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	_, ctx := tracer.StartSpanFromContext(context.Background(), "test")
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.Error(t, err) // The service doesn't exist so an error should occur.

	// Query the mock tracer for finished spans.
	spans := mt.FinishedSpans()

	// The target url finish by "k8s.leboncoin.lan" so it's handled as an micro-service and no
	// custom ServiceName is set.
	require.GreaterOrEqual(t, len(spans), 1)
	span := spans[1]

	assert.Equal(t, "some-service.some-env.k8s.leboncoin.lan", span.Tag(ext.ServiceName))
	assert.Equal(t, http.MethodGet, span.Tag("http.method"))
	assert.Equal(t, scheme, span.Tag("http.scheme"))
	assert.Equal(t, host, span.Tag("http.host"))
	assert.Equal(t, path, span.Tag("http.path"))
	assert.Equal(t, queryparams, span.Tag("http.queryparams"))
}

func TestTracing_Intercept_with_a_call_to_an_external_service(t *testing.T) {
	client := new(http.Client)

	// Start the mock tracer.
	mt := mocktracer.Start()
	defer mt.Stop()

	client.Transport = train.Transport(NewTracing())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)

	_, ctx := tracer.StartSpanFromContext(context.Background(), "test")
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.NoError(t, err)

	// Query the mock tracer for finished spans.
	spans := mt.FinishedSpans()

	// The server is in localhost so it's handled as an external service and a
	// custom ServiceName is set.
	require.NotNil(t, spans[0].Tag(ext.ServiceName))
}

func TestTracing_Intercept_with_secrets_in_query_params(t *testing.T) {
	client := new(http.Client)

	mt := mocktracer.Start()
	defer mt.Stop()

	client.Transport = train.Transport(
		NewTracing().WithSecretQueryParams("bar", "qux", "quux", "quuux"),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	u.Path += "/path"
	params := url.Values{}
	params.Add("foo", "apple")
	params.Add("bar", "")
	params.Add("bar", "banana")
	params.Add("bar", "cherry")
	params.Add("baz", "lemon")
	params.Add("qux", "orange")
	params.Add("quuux", "")
	params.Add("quuux", "")
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	require.NoError(t, err)

	_, ctx := tracer.StartSpanFromContext(context.Background(), "test")
	req = req.WithContext(ctx)

	_, err = client.Do(req)
	require.NoError(t, err)

	spans := mt.FinishedSpans()
	tag := spans[len(spans)-1].Tag("http.queryparams")
	assert.NotContains(t, tag, "banana")
	assert.NotContains(t, tag, "cherry")
	assert.NotContains(t, tag, "orange")
	assert.Contains(t, tag, "apple")
	assert.Contains(t, tag, "lemon")

	assert.NotContains(t, tag, "quux", "missing secret param should not be added")

	qp, err := url.ParseQuery(tag.(string))
	require.NoError(t, err)
	assert.Empty(t, strings.Join(qp["quuux"], ""), "empty secret param should be left emtpy")
}

func TestTracingWithRequestBody_Intercept_with_a_call_to_an_internal_service(t *testing.T) {
	var tests = []struct {
		tryIndentJSON bool
		method        string
		body          *string
		expectedBody  any
	}{
		{
			tryIndentJSON: true,
			method:        http.MethodGet,
			body:          nil,
			expectedBody:  nil,
		},
		{
			tryIndentJSON: false,
			method:        http.MethodGet,
			body:          pointer.String(`{"withUserID":"1234","adID":"789","msg":"Hello, I wand to buy your jean for 10 euros"}`),
			expectedBody:  `{"withUserID":"1234","adID":"789","msg":"Hello, I wand to buy your jean for 10 euros"}`,
		},
		{
			tryIndentJSON: true,
			method:        http.MethodGet,
			body:          pointer.String(`{"withUserID":"21234","adID":"789", "msg":"Hello, I wand to buy your jean for 10 euros"}`),
			expectedBody: `{
  "withUserID": "21234",
  "adID": "789",
  "msg": "Hello, I wand to buy your jean for 10 euros"
}`,
		},
		{
			tryIndentJSON: false,
			method:        http.MethodGet,
			body:          randStringRunes(tracingRequestBodyMaxContentLength + 1),
			expectedBody:  `tracing info: contentlength is to big for tracing`,
		},
	}

	for _, test := range tests {
		tracing := NewTracing().WithRequestBodyLog(test.tryIndentJSON)
		client := new(http.Client)
		client.Transport = train.Transport(tracing)

		// Start the mock tracer.
		mt := mocktracer.Start()
		defer mt.Stop()

		url := "https://some-service.some-env.k8s.leboncoin.lan/messages"
		var body io.Reader
		if test.body != nil {
			body = strings.NewReader(*test.body)
		}

		req, err := http.NewRequest(test.method, url, body)
		require.NoError(t, err)

		_, ctx := tracer.StartSpanFromContext(context.Background(), "test")
		req = req.WithContext(ctx)

		_, err = client.Do(req)
		require.Error(t, err) // The service doesn't exist so an error should occur.

		// Query the mock tracer for finished spans.
		spans := mt.FinishedSpans()
		require.GreaterOrEqual(t, len(spans), 1)

		span := spans[1]
		assert.Equal(t, test.method, span.Tag("http.method"))
		assert.Equal(t, test.expectedBody, span.Tag("http.request_body"))
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randStringRunes returns a pointer to a random string with n characters.
func randStringRunes(n int) *string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	value := string(b)
	return &value
}
