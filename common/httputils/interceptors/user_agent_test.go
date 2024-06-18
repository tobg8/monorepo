package interceptors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/f2prateek/train"
	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	client := http.Client{}
	client.Transport = train.Transport(NewUserAgent("my-service", "1.2.3"))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "my-service/1.2.3", req.Header.Get("User-Agent"))
	}))
	defer ts.Close()

	_, err := client.Get(ts.URL)
	assert.NoError(t, err)
}
