package interceptors

import (
	"fmt"
	"net/http"

	"github.com/f2prateek/train"
)

// UserAgent propagates the User-Agent header with the application name in it.
type UserAgent string

// NewUserAgent instantiates an UserAgent interceptor
func NewUserAgent(appName string, appVersion string) *UserAgent {
	res := UserAgent(fmt.Sprintf("%s/%s", appName, appVersion))
	return &res
}

// Intercept implements the train.Interceptor interface
func (u *UserAgent) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()

	req.Header.Set("User-Agent", string(*u))

	return chain.Proceed(req)
}
