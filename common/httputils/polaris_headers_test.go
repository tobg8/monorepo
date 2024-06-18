package httputils

import (
	"testing"
)

func Test_get_all_headers(t *testing.T) {
	//req, _ := http.NewRequest(http.MethodGet, "http://test", nil)
	//req.Header.Set("X-Unique-Id", "unique id")
	//req.Header.Set("Authorization", "auth")
	//req.Header.Set("not-copied", "dummy")
	//
	//headers := GetPolarisHeaders(req)
	//
	//assert.Equal(t, []string{"unique id"}, headers[polarisheaders.UniqueID])
	//assert.Equal(t, []string{"auth"}, headers["Authorization"])
	//assert.Nil(t, headers["not-copied"])
}

func Test_set_all_headers(t *testing.T) {

	//req, _ := http.NewRequest(http.MethodGet, "http://test", nil)
	//
	//headers := PolarisHeaders{}
	//
	//headers[polarisheaders.UniqueID] = []string{"unique id"}
	//headers[polarisheaders.Authorization] = []string{"auth"}
	//
	//SetPolarisHeaders(req, headers)
	//
	//assert.Equal(t, "unique id", req.Header.Get(polarisheaders.UniqueID))
	//assert.Equal(t, "auth", req.Header.Get(polarisheaders.Authorization))
}
