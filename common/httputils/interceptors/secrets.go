package interceptors

import (
	"net/http"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/secret"
)

// QueryObfuscator will obfuscate the query from a request.
//
// This allows to hide sensitives informations from the logs
type QueryObfuscator struct {
	secret.QueryParams
}

// Intercept is an implementation of train.Intercept.
//
// It is used to catch and interact with the http request transparently.
func (o *QueryObfuscator) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()

	res, err := chain.Proceed(req)

	// The request is not more used so we can modify its queries.
	//
	// Those are the queries used for the url generation inside the response errors.
	req.URL.RawQuery = o.HideFromValues(req.URL.Query()).Encode()

	return res, err
}
