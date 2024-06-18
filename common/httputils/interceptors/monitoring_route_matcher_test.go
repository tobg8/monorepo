package interceptors

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaticRouteMatcher(t *testing.T) {
	testRoute := "/foo/bar"
	testMethod := http.MethodGet
	expectedRoute := "get:" + testRoute

	newRequest := func(t *testing.T, method string, path string, queryParams []string, body io.Reader) *http.Request {
		url := fmt.Sprintf("http://test.example.com%s", path)
		if len(queryParams) > 0 {
			url += "?" + strings.Join(queryParams, "&")
		}

		req, err := http.NewRequest(method, url, body)
		require.NoError(t, err)

		return req
	}

	testStaticRoute := func(t *testing.T, req *http.Request) (string, bool) {
		rm := StaticRouteMatcher(testMethod, testRoute)
		return rm.MatchRequest(req)
	}

	t.Run("match", func(t *testing.T) {
		testValidStaticRoute := func(t *testing.T, req *http.Request) {
			route, ok := testStaticRoute(t, req)
			assert.True(t, ok)
			assert.Equal(t, expectedRoute, route)
		}

		t.Run("without-query-params", func(t *testing.T) {
			req := newRequest(t, testMethod, testRoute, nil, nil)
			testValidStaticRoute(t, req)
		})

		t.Run("with-query-params", func(t *testing.T) {
			req := newRequest(
				t, testMethod, testRoute, []string{
					"param1=value1",
					"param2=value2",
				}, nil,
			)
			testValidStaticRoute(t, req)
		})

		t.Run("with-body", func(t *testing.T) {
			req := newRequest(
				t, testMethod, testRoute, nil,
				strings.NewReader("test request body"),
			)
			testValidStaticRoute(t, req)
		})
	})

	t.Run("not-match", func(t *testing.T) {
		testNotValidStaticRoute := func(t *testing.T, req *http.Request) {
			route, ok := testStaticRoute(t, req)
			assert.False(t, ok)
			assert.Zero(t, route)
		}

		t.Run("with-different-method", func(t *testing.T) {
			for _, method := range []string{
				http.MethodConnect,
				http.MethodDelete,
				http.MethodGet,
				http.MethodHead,
				http.MethodOptions,
				http.MethodPatch,
				http.MethodPost,
				http.MethodPut,
				http.MethodTrace,
			} {
				if method == testMethod {
					continue
				}

				t.Run(method, func(t *testing.T) {
					req := newRequest(t, method, testRoute, nil, nil)
					testNotValidStaticRoute(t, req)
				})
			}
		})

		t.Run("with-different-path", func(t *testing.T) {
			for _, path := range []string{
				"/buzz/qux",
				testRoute + "/quuz",
				testRoute + "/quuz/garply",
				"/corge" + testRoute,
				"/corge" + testRoute + "/waldo",
			} {
				t.Run(path, func(t *testing.T) {
					t.Run("without-query-params", func(t *testing.T) {
						req := newRequest(t, testMethod, path, nil, nil)
						testNotValidStaticRoute(t, req)
					})

					t.Run("with-query-params", func(t *testing.T) {
						req := newRequest(t, testMethod, path, []string{
							"param1=value1",
							"param2=value2",
						}, nil)
						testNotValidStaticRoute(t, req)
					})

					t.Run("with-body", func(t *testing.T) {
						req := newRequest(
							t, testMethod, path, nil,
							strings.NewReader("test request body"),
						)
						testNotValidStaticRoute(t, req)
					})
				})
			}
		})
	})
}

func TestDynamicRouteMatcher(t *testing.T) {
	testRoutePattern, err := regexp.Compile(
		fmt.Sprintf(
			"/foo/\\d+/bar/%[1]s{8}-%[1]s{4}-%[1]s{4}-%[1]s{4}-%[1]s{12}",
			"[a-fA-F0-9]",
		),
	)
	require.NoError(t, err)

	testRoute := "/foo/<store-id>/bar/<user-id>"
	validTestRoutes := []string{
		"/foo/1/bar/d3076749-1b44-4520-8515-eb26ff84913c",
		"/foo/123456789/bar/0a373eca-df9b-4879-afd4-52a4b1c6d32c",
	}
	testMethod := http.MethodPost
	expectedRoute := "post:" + testRoute

	newRequest := func(t *testing.T, method string, path string, queryParams []string, body io.Reader) *http.Request {
		url := fmt.Sprintf("http://test.example.com%s", path)
		if len(queryParams) > 0 {
			url += "?" + strings.Join(queryParams, "&")
		}

		req, err := http.NewRequest(method, url, body)
		require.NoError(t, err)

		return req
	}

	testDynamicRoute := func(t *testing.T, req *http.Request) (string, bool) {
		rm := DynamicRouteMatcher(testMethod, testRoute, testRoutePattern)
		return rm.MatchRequest(req)
	}

	t.Run("match", func(t *testing.T) {
		testValidDynamicRoute := func(t *testing.T, req *http.Request) {
			route, ok := testDynamicRoute(t, req)
			assert.True(t, ok)
			assert.Equal(t, expectedRoute, route)
		}

		for idx, validTestRoute := range validTestRoutes {
			t.Run(fmt.Sprintf("valid-route-%d", idx), func(t *testing.T) {
				t.Logf("route: %s", validTestRoute)

				t.Run("without-query-params", func(t *testing.T) {
					req := newRequest(t, testMethod, validTestRoute, nil, nil)
					testValidDynamicRoute(t, req)
				})

				t.Run("with-query-params", func(t *testing.T) {
					req := newRequest(
						t, testMethod, validTestRoute, []string{
							"param1=value1",
							"param2=value2",
						}, nil,
					)
					testValidDynamicRoute(t, req)
				})

				t.Run("with-body", func(t *testing.T) {
					req := newRequest(
						t, testMethod, validTestRoute, nil,
						strings.NewReader("test request body"),
					)
					testValidDynamicRoute(t, req)
				})
			})
		}
	})

	t.Run("not-match", func(t *testing.T) {
		testNotValidDynamicRoute := func(t *testing.T, req *http.Request) {
			route, ok := testDynamicRoute(t, req)
			assert.False(t, ok)
			assert.Zero(t, route)
		}

		t.Run("with-different-method", func(t *testing.T) {
			for _, method := range []string{
				http.MethodConnect,
				http.MethodDelete,
				http.MethodGet,
				http.MethodHead,
				http.MethodOptions,
				http.MethodPatch,
				http.MethodPost,
				http.MethodPut,
				http.MethodTrace,
			} {
				if method == testMethod {
					continue
				}

				t.Run(method, func(t *testing.T) {
					req := newRequest(t, method, validTestRoutes[0], nil, nil)
					testNotValidDynamicRoute(t, req)
				})
			}
		})

		t.Run("with-different-path", func(t *testing.T) {
			for _, path := range []string{
				"/buzz/qux",
				testRoute + "/quuz",
				testRoute + "/quuz/garply",
				"/corge" + testRoute,
				"/corge" + testRoute + "/waldo",

				"/foo/bar/d3076749-1b44-4520-8515-eb26ff84913c",
				"/foo/1x2z/bar/0a373eca-df9b-4879-afd4-52a4b1c6d32c",
				"/foo/1/bar/123",
				"/foo/1/bar/buz",
				"/foo/1/bar/0a373eca",
			} {
				t.Run(path, func(t *testing.T) {
					t.Run("without-query-params", func(t *testing.T) {
						req := newRequest(t, testMethod, path, nil, nil)
						testNotValidDynamicRoute(t, req)
					})

					t.Run("with-query-params", func(t *testing.T) {
						req := newRequest(t, testMethod, path, []string{
							"param1=value1",
							"param2=value2",
						}, nil)
						testNotValidDynamicRoute(t, req)
					})

					t.Run("with-body", func(t *testing.T) {
						req := newRequest(
							t, testMethod, path, nil,
							strings.NewReader("test request body"),
						)
						testNotValidDynamicRoute(t, req)
					})
				})
			}
		})
	})
}
