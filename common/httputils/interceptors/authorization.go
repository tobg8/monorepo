package interceptors

import (
	"fmt"
	"net/http"

	"github.com/f2prateek/train"

	"github.com/monorepo/common/contextkeys"
)

// Authorization implements train.Interceptor interface.
type Authorization struct{}

// NewAuthorization returns a new authorization interceptor.
func NewAuthorization() *Authorization {
	return &Authorization{}
}

// Intercept implements train.Interceptor interface.
func (l Authorization) Intercept(chain train.Chain) (*http.Response, error) {
	var (
		req   = chain.Request()
		token = req.Context().Value(contextkeys.AuthToken)
	)

	if req.Header.Get("Authorization") == "" {
		if t, ok := token.(string); ok && t != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
		}
	}

	return chain.Proceed(req)
}
