package httputils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/monitoring/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genInterceptor generates a custom HTTP interceptor.
// The latter will modify a header of the HTTP intercepted request to
// keep a trace of its execution.
// To do that, we update a specific header in the HTTP request with a
// list of all encountered interceptors identifier.
// This list will be used by the HTTP handler in the test function to
// verify the order execution of all interceptors.
func genInterceptor(id int, headerName string) train.Interceptor {
	return train.InterceptorFunc(
		func(chain train.Chain) (*http.Response, error) {
			httpRequest := chain.Request()

			value := httpRequest.Header.Get(headerName)
			if value != "" {
				value += ","
			}
			value += strconv.Itoa(id)
			httpRequest.Header.Set(headerName, value)

			return chain.Proceed(httpRequest)
		},
	)
}

// Test_Client_Interceptors check the order execution of all
// registered HTTP interceptors.
func Test_Client_Interceptors(t *testing.T) {
	const headerName = "X-Test-Interceptors"

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, r.Header.Get(headerName))
		}),
	)
	defer ts.Close()

	interceptors := make([]train.Interceptor, 6)
	for idx := range interceptors {
		interceptors[idx] = genInterceptor(idx+1, headerName)
	}

	client := NewClient(1*time.Second, 0)
	client.appendInterceptors(interceptors[3], interceptors[4])
	client.prependInterceptors(interceptors[2])
	client.appendInterceptors(interceptors[5])
	client.prependInterceptors(interceptors[0], interceptors[1])

	httpRequest, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	var responseBody bytes.Buffer
	statusCode, err := client.Do(
		context.Background(),
		&responseBody,
		httpRequest,
	)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "1,2,3,4,5,6", responseBody.String())
}

func Test_Client_Do_can_ignore_the_body_if_the_destination_parameter_is_nil(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello World"))
		}),
	)
	defer ts.Close()

	client := NewClient(1*time.Second, 0)
	httpRequest, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	statusCode, err := client.Do(
		context.Background(),
		io.Discard,
		httpRequest,
	)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
}

func Test_Client_DoAndUnmarshalJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"Test": "test"}`))
		}),
	)
	defer ts.Close()

	client := NewClient(1*time.Second, 0)
	httpRequest, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	to := make(map[string]string)
	statusCode, err := client.DoAndUnmarshalJSON(
		context.Background(),
		&to,
		httpRequest,
	)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, map[string]string{"Test": "test"}, to)
}

func Test_Client_DoAndUnmarshalJSON_nicer_error_on_504(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusGatewayTimeout)
			_, _ = w.Write([]byte(`<gateway timeout>`))
		}),
	)
	defer ts.Close()

	client := NewClient(1*time.Second, 0)
	httpRequest, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	to := make(map[string]string)
	statusCode, err := client.DoAndUnmarshalJSON(
		context.Background(),
		&to,
		httpRequest,
	)
	require.Equal(t, http.StatusGatewayTimeout, statusCode)
	require.Equal(t, errors.New(`decode status 504: body "<gateway timeout>"`), err)
	require.Empty(t, to)
}

func Test_Client_DoAndUnmarshalJSON_504_still_default_to_parsing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusGatewayTimeout)
			_, _ = w.Write([]byte(`{"Test": "test"}`))
		}),
	)
	defer ts.Close()

	client := NewClient(1*time.Second, 0)
	httpRequest, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	to := make(map[string]string)
	statusCode, err := client.DoAndUnmarshalJSON(
		context.Background(),
		&to,
		httpRequest,
	)

	require.NoError(t, err)
	require.Equal(t, http.StatusGatewayTimeout, statusCode)
	require.Equal(t, map[string]string{"Test": "test"}, to)
}

//func TestNewClient_adds_unique_id_interceptor(t *testing.T) {
//	client := NewClient(time.Second, 0)
//
//	var uniqueID string
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		uniqueID = req.Header.Get(polarisheaders.UniqueID)
//	}))
//	defer ts.Close()
//
//	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
//	require.NoError(t, err)
//
//	ctx := context.WithValue(context.Background(), contextkeys.UniqueID, "my_unique_id")
//	req = req.WithContext(ctx)
//
//	_, err = client.Client.Do(req)
//	assert.NoError(t, err)
//	assert.Equal(t, "my_unique_id", uniqueID)
//}

//func TestNewClient_adds_brand_forwarding_interceptor(t *testing.T) {
//	client := NewClient(time.Second, 0)
//
//	var forwardedHost string
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		forwardedHost = req.Header.Get(polarisheaders.ForwardedHost)
//	}))
//	defer ts.Close()
//
//	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
//	require.NoError(t, err)
//
//	ctx := context.WithValue(context.Background(), contextkeys.BrandOriginHost, "my_brand")
//	req = req.WithContext(ctx)
//
//	_, err = client.Client.Do(req)
//	assert.NoError(t, err)
//	assert.Equal(t, "my_brand", forwardedHost)
//}

func TestNewClient_timeouts(t *testing.T) {
	client := NewClient(100*time.Millisecond, 0)

	var counter int64
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		atomic.AddInt64(&counter, 1)
		<-req.Context().Done()
	}))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)

	_, err = client.Client.Do(req)
	assert.Error(t, err)

	ts.Close()
	assert.Equal(t, int64(1), atomic.LoadInt64(&counter)) // No retry
}

func Test_panicIfAlreadySet_should_panic_if_double_tracer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if !ok || e == nil {
				t.Fail()
			}
		}
	}()
	client := NewClient(1*time.Second, 1*time.Second)
	client.WithTracer()
	client.WithTracer()
	t.Fail()
}

func Test_panicIfAlreadySet_should_panic_if_double_monitor(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if !ok || e == nil {
				t.Fail()
			}
		}
	}()
	client := NewClient(1*time.Second, 1*time.Second)
	client.WithMonitor(metrics.NoopStatsdHandler)
	client.WithMonitor(metrics.NoopStatsdHandler)
	t.Fail()
}

func Test_panicIfAlreadySet_should_not_panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fail()
		}
	}()
	client := NewClient(1*time.Second, 1*time.Second)
	client.WithTracer()
	client.WithMonitor(metrics.NoopStatsdHandler)
}

func Test_Client_WithSecretQueryParams(t *testing.T) {
	t.Run("removes secrets from Do error", func(t *testing.T) {
		const testSecret = `WjAh_Bt.ZXf"HG@p!Nuo*RXAAx9db9rxE`

		// Use unsupported scheme to generate http Do error that contains full
		// url in message without dns resolution.
		u, err := url.Parse("example://domain.test/path")
		require.NoError(t, err)
		u.RawQuery = url.Values{"secret": []string{testSecret}}.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		require.NoError(t, err)

		_, err = NewClient(time.Second, time.Second).
			WithSecretQueryParams("secret").
			Do(context.Background(), io.Discard, req)
		assert.Error(t, err)
		assert.NotContains(t, err.Error(), testSecret)
		assert.NotContains(t, err.Error(), url.QueryEscape(testSecret))
	})

	t.Run("panics when tracer is already set", func(t *testing.T) {
		assert.Panics(t, func() {
			NewClient(time.Second, time.Second).
				WithTracer().
				WithSecretQueryParams("secret")
		})
	})
}
