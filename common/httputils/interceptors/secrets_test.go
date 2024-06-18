package interceptors

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_QueryParams_Intercept(t *testing.T) {
	t.Run("is ok", func(t *testing.T) {
		// Setup the client with the interceptor
		qp := secret.QueryParams{"bar", "qux"}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transportWithInterceptors := train.TransportWith(transport, &QueryObfuscator{qp})

		client := http.Client{Transport: transportWithInterceptors}

		params := url.Values{}
		params.Add("foo", "apple")
		params.Add("bar", "banana")
		params.Add("bar", "cherry")
		params.Add("baz", "lemon")
		params.Add("qux", "orange")
		req, err := http.NewRequest("GET", "invalid-request", nil)
		require.NoError(t, err)

		req.URL.RawQuery = params.Encode()

		res, err := client.Do(req)

		assert.Nil(t, res)
		assert.EqualError(t, err, "Get \"invalid-request?bar=xxxxx&bar=xxxxx&baz=lemon&foo=apple&qux=xxxxx\": unsupported protocol scheme \"\"")
	})
}
