package svcauth

import (
	"fmt"
	"net/http"

	"github.com/f2prateek/train"
)

// ServiceAuth implements train.Interceptor interface.
type ServiceAuth struct {
	*TokenGetter
	intercept func(chain train.Chain) (*http.Response, error)
}

// NewServiceAuth returns a new authorization interceptor.
func NewServiceAuth(conf Conf) *ServiceAuth {
	svcAuth := ServiceAuth{
		TokenGetter: NewTokenGetter(conf),
	}

	if conf.Enabled {
		svcAuth.intercept = svcAuth.interceptEnabled
	} else {
		svcAuth.intercept = svcAuth.interceptDisabled
	}

	return &svcAuth
}

// Intercept implements train.Interceptor interface.
func (l *ServiceAuth) Intercept(chain train.Chain) (*http.Response, error) {
	return l.intercept(chain)
}

func (l *ServiceAuth) interceptDisabled(chain train.Chain) (*http.Response, error) {
	return chain.Proceed(chain.Request())
}

func (l *ServiceAuth) interceptEnabled(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	token, err := l.GetToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("can't execute request: %w", err)
	}
	if err := l.addAuthorizationToken(req, token); err != nil {
		return nil, err
	}

	resp, err := chain.Proceed(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	// unauthorized, try with a new token
	token, err = l.ForceTokenRenew(req.Context())
	if err != nil {
		return nil, fmt.Errorf("can't execute request: %w", err)
	}
	req = req.Clone(req.Context())
	if err := l.addAuthorizationToken(req, token); err != nil {
		return nil, err
	}
	return chain.Proceed(req)
}

func (l *ServiceAuth) addAuthorizationToken(req *http.Request, token string) error {
	const authorizationHeaderKey = "Authorization"

	if req.Header.Get(authorizationHeaderKey) == "" {
		req.Header.Set(authorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
	}

	return nil
}
