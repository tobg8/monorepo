package interceptors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f2prateek/train"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monorepo/common/contextkeys"
)

func TestAuthorization_Intercept(t *testing.T) {
	var (
		client http.Client
		token  = "my-super-token"
	)
	client.Transport = train.Transport(NewAuthorization())

	t.Run("add token in header", func(t *testing.T) {
		var srv = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, "Bearer "+token, req.Header.Get("Authorization"))
			}),
		)
		defer srv.Close()
		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(),
			contextkeys.AuthToken, token,
		))
		_, err = client.Do(req)
		require.NoError(t, err)
	})

	t.Run("already existing token in header", func(t *testing.T) {
		var srv = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, "Bearer my-other-token", req.Header.Get("Authorization"))
			}),
		)
		defer srv.Close()
		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(),
			contextkeys.AuthToken, token,
		))
		req.Header.Set("Authorization", "Bearer my-other-token")
		_, err = client.Do(req)
		require.NoError(t, err)
	})
}
