package svcauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/f2prateek/train"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAuth_Intercept(t *testing.T) {
	token := "test-token"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseForm())
		require.Equal(t, conf.ClientID, req.PostForm.Get("client_id"))
		require.Equal(t, string(conf.ClientSecret), req.PostForm.Get("client_secret"))
		require.Equal(t, conf.RequiredScopes, []string{
			req.PostForm.Get("scope"),
		})
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   10,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		require.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", token), req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusOK)
	}))
	conf.AuthorizerURL = authorizer.URL
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}
	resp, err := client.Get(svc.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServiceAuth_Intercept_doesnt_override_existing_token(t *testing.T) {
	token := "test-token-from-authorizer"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseForm())
		require.Equal(t, conf.ClientID, req.PostForm.Get("client_id"))
		require.Equal(t, string(conf.ClientSecret), req.PostForm.Get("client_secret"))
		require.Equal(t, conf.RequiredScopes, []string{
			req.PostForm.Get("scope"),
		})
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   10,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		require.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", "test-token-not-from-authorizer"), req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusOK)
	}))
	conf.AuthorizerURL = authorizer.URL
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}
	req, err := http.NewRequest("GET", svc.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer test-token-not-from-authorizer")

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServiceAuth_Intercept_renew_token_on_unauthorized(t *testing.T) {
	token := "test-token"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseForm())
		require.Equal(t, conf.ClientID, req.PostForm.Get("client_id"))
		require.Equal(t, string(conf.ClientSecret), req.PostForm.Get("client_secret"))
		require.Equal(t, conf.RequiredScopes, []string{
			req.PostForm.Get("scope"),
		})
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   10,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		require.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	conf.AuthorizerURL = authorizer.URL

	ncalls := 0
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", token), req.Header.Get("Authorization"))
		ncalls++
		switch ncalls {
		case 1:
			resp.WriteHeader(http.StatusUnauthorized)
		default:
			resp.WriteHeader(http.StatusOK)
		}
	}))
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}
	resp, err := client.Get(svc.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServiceAuth_InterceptDisabled(t *testing.T) {
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		t.Error("authorizer should not be called")
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		assert.Empty(t, req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusOK)
	}))
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(Conf{
			Enabled:       false,
			AuthorizerURL: authorizer.URL,
		})),
	}
	resp, err := client.Get(svc.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServiceAuth_InterceptDisabled_no_renewal_on_401(t *testing.T) {
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		t.Error("authorizer should not be called")
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		assert.Empty(t, req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusUnauthorized)
	}))
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(Conf{
			Enabled:       false,
			AuthorizerURL: authorizer.URL,
		})),
	}
	resp, err := client.Get(svc.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestServiceAuth_Intercept_dont_call_if_token_valid(t *testing.T) {
	token := "test-token"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	count := 0
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseForm())
		require.Equal(t, conf.ClientID, req.PostForm.Get("client_id"))
		require.Equal(t, string(conf.ClientSecret), req.PostForm.Get("client_secret"))
		require.Equal(t, conf.RequiredScopes, []string{
			req.PostForm.Get("scope"),
		})
		count++
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   1000,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		require.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", token), req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusOK)
	}))
	conf.AuthorizerURL = authorizer.URL
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}
	resp, err := client.Get(svc.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.Get(svc.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, 1, count)
}

func TestServiceAuth_Intercept_call_if_token_expire(t *testing.T) {
	token := "test-token"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	count := 0
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseForm())
		require.Equal(t, conf.ClientID, req.PostForm.Get("client_id"))
		require.Equal(t, string(conf.ClientSecret), req.PostForm.Get("client_secret"))
		require.Equal(t, conf.RequiredScopes, []string{
			req.PostForm.Get("scope"),
		})
		count++
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   -1,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		require.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", token), req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusOK)
	}))
	conf.AuthorizerURL = authorizer.URL
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}
	resp, err := client.Get(svc.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.Get(svc.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, 2, count)
}

func TestServiceAuth_Intercept_dont_call_twice_in_parrallel(t *testing.T) {
	token := "test-token"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	var count int64
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		atomic.AddInt64(&count, 1)

		require.NoError(t, req.ParseForm())
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   100000,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		require.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", token), req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusOK)
	}))
	conf.AuthorizerURL = authorizer.URL
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get(svc.URL)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(1), count)
}

func TestServiceAuth_Intercept_is_thread_safe_even_forced_renewal(t *testing.T) {
	token := "test-token"
	conf := Conf{
		ClientID:     "testID",
		ClientSecret: "testSecret",
		RequiredScopes: []string{
			"test",
		},
		Enabled: true,
	}
	var count int64
	authorizer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		atomic.AddInt64(&count, 1)

		assert.NoError(t, req.ParseForm())
		data, err := json.Marshal(TokenResp{
			AccessToken: token,
			ExpiresIn:   100000,
			Scope:       req.Form.Get("scope"),
			TokenType:   "bearer",
		})
		assert.NoError(t, err)
		_, _ = resp.Write(data)
	}))
	svc := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		require.Equal(t, fmt.Sprintf("Bearer %s", token), req.Header.Get("Authorization"))
		resp.WriteHeader(http.StatusUnauthorized)
	}))
	conf.AuthorizerURL = authorizer.URL
	client := &http.Client{
		Transport: train.Transport(NewServiceAuth(conf)),
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get(svc.URL)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(11), count, "should be called once at start and once per goroutine due to 401")
}
