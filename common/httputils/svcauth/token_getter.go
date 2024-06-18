package svcauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/f2prateek/train"
	"golang.org/x/net/context/ctxhttp"

	"github.com/monorepo/common/configloader"
	"github.com/monorepo/common/httputils/interceptors"
	"github.com/monorepo/common/secret"
)

// Conf is the configuration object for the service authentication interceptor
type Conf struct {
	AuthorizerURL  string        `mapstructure:"authorizer_url"`
	ClientID       string        `mapstructure:"client_id"`
	ClientSecret   secret.String `mapstructure:"client_secret"`
	RequiredScopes []string      `mapstructure:"scopes"`
	Enabled        bool          `mapstructure:"enabled"`
}

// String implements the Stringer interface and mask the password to avoid leaking
func (c Conf) String() string {
	jsonConfig, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(jsonConfig)
}

// Envs bind environment keys to env variables
func (*Conf) Envs(l *configloader.Loader) {
	l.BindEnv("authorizer_url")
	l.BindEnv("client_id")
	l.BindEnv("client_secret")
	l.BindEnv("scopes")
}

// Defaults sets the default values for configuration keys
func (*Conf) Defaults(l *configloader.Loader) {
	l.SetDefault("authorizer_url", "http://authorizer.svc.disco/api/authorizer/v2/token")
	l.SetDefault("scopes", "*")
	l.SetDefault("enabled", true)
}

// TokenGetter fetches a token from authorizer and feed a request authorization header
type TokenGetter struct {
	authorizerURL         string
	httpClient            *http.Client
	authorizerRequestData url.Values
	tokenMutex            sync.RWMutex
	expiresAt             time.Time
	token                 TokenResp
}

// NewTokenGetter returns a new token getter.
func NewTokenGetter(conf Conf) *TokenGetter {
	return &TokenGetter{
		authorizerURL: conf.AuthorizerURL,
		httpClient: &http.Client{
			Transport: train.Transport(interceptors.NewTracing()),
		},
		authorizerRequestData: url.Values{
			"client_id": []string{conf.ClientID},
			"client_secret": []string{
				string(conf.ClientSecret),
			},
			"scope": []string{
				strings.Join(conf.RequiredScopes, " "),
			},
			"grant_type": []string{"client_credentials"},
		},
	}
}

func (l *TokenGetter) retrieveTokenLocked(ctx context.Context) (TokenResp, error) {
	token, err := l.retrieveNewToken(ctx)
	if err != nil {
		return TokenResp{}, fmt.Errorf("can't renew token: %w", err)
	}

	l.token = token
	l.expiresAt = time.Now().Add((time.Duration(token.ExpiresIn) * time.Second) - time.Minute)
	return l.token, nil
}

// ForceTokenRenew force a token cache renew without checking if it has expired
func (l *TokenGetter) ForceTokenRenew(ctx context.Context) (string, error) {
	l.tokenMutex.Lock()
	defer l.tokenMutex.Unlock()
	resp, err := l.retrieveTokenLocked(ctx)
	if err != nil {
		return "", err
	}
	return resp.AccessToken, err
}

// GetToken retrieve token from cache and renew it when it has expired
func (l *TokenGetter) GetToken(ctx context.Context) (string, error) {
	l.tokenMutex.RLock()
	if token, ok := l.access(); ok {
		l.tokenMutex.RUnlock()
		return token, nil
	}
	l.tokenMutex.RUnlock()
	l.tokenMutex.Lock()
	if token, ok := l.access(); ok {
		l.tokenMutex.Unlock()
		return token, nil
	}

	token, err := l.retrieveTokenLocked(ctx)
	l.tokenMutex.Unlock()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func (l *TokenGetter) access() (string, bool) {
	return l.token.AccessToken, !l.expiresAt.Before(time.Now())
}

func (l *TokenGetter) retrieveNewToken(ctx context.Context) (TokenResp, error) {
	res, err := ctxhttp.PostForm(ctx, l.httpClient, l.authorizerURL, l.authorizerRequestData)
	if err != nil {
		return TokenResp{}, fmt.Errorf("can't retrieve service token: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TokenResp{}, fmt.Errorf("can't read response from authorizer: %w", err)
	}
	if res.StatusCode >= 400 {
		return TokenResp{}, fmt.Errorf("authorizer return an error when retrieving token %q", string(body))
	}
	var resp TokenResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return TokenResp{}, fmt.Errorf("can't unmarshal authorizer response: %w", err)
	}
	return resp, nil
}

// TokenResp defines the structure of authorizer token response
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
