package httputils

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type tMock struct {
	mock.Mock
	testing.T
}

func (m *tMock) Cleanup(f func()) {
	m.Called(f)
}

func (m *tMock) Logf(format string, args ...interface{}) {
	m.Called(format, args)
}

func Test_RequestMatcher(t *testing.T) {
	t.Run("matches full request", func(t *testing.T) {
		mc := NewMockClient(t)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(t).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`{"foo":"1","bar":{"baz":["2","3"]}}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1")
		r.Header.Set("Bar", "2")

		_, err = mc.Do(context.Background(), io.Discard, r)
		require.NoError(t, err)

		mc.AssertExpectations(t)
	})

	t.Run("matches json in any key order", func(t *testing.T) {
		mc := NewMockClient(t)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(t).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`{"bar":{"baz":["2","3"]},"foo":"1"}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1")
		r.Header.Set("Bar", "2")

		_, err = mc.Do(context.Background(), io.Discard, r)
		require.NoError(t, err)

		mc.AssertExpectations(t)
	})

	t.Run("matches raw content even if not parseable", func(t *testing.T) {
		t.Skip()
		mc := NewMockClient(t)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(t).
				RawBody(`some raw content`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`some raw conent`))
		require.NoError(t, err)

		_, err = mc.Do(context.Background(), io.Discard, r)
		require.NoError(t, err)

		mc.AssertExpectations(t)
	})

	t.Run("reject raw content even if the json would match", func(t *testing.T) {
		expected := `--- expected
+++ actual
@@ -2,4 +2,4 @@
 Host: domain.test
 Content-Length: 41
 
-{"a": "first value", "b": "second value"}
+{"b": "second value", "a": "first value"}
`

		cleanupFuncs := []func(){}

		tm := &tMock{}
		tm.Test(t)

		tm.On("Cleanup", mock.Anything).Run(func(args mock.Arguments) {
			cleanupFuncs = append(cleanupFuncs, args.Get(0).(func()))
		}).Return()

		tm.On("Logf", "httputils unmatched expectation:\n%s", []interface{}{expected}).
			Return().Once()

		tm.On("Logf", mock.Anything, mock.Anything).Return()

		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(tm).
				RawBody(`{"a": "first value", "b": "second value"}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`{"b": "second value", "a": "first value"}`))
		require.NoError(t, err)

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}

		mc.AssertExpectations(tm)
		tm.AssertExpectations(t)
		assert.True(t, tm.Failed())
	})

	t.Run("matches ignoring extra headers", func(t *testing.T) {
		mc := NewMockClient(t)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(t).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`{"foo":"1","bar":{"baz":["2","3"]}}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1")
		r.Header.Set("Bar", "2")
		r.Header.Set("Baz", "3")
		r.Header.Set("Qux", "4")

		_, err = mc.Do(context.Background(), io.Discard, r)
		require.NoError(t, err)

		mc.AssertExpectations(t)
	})

	t.Run("matches user-agent header", func(t *testing.T) {
		mc := NewMockClient(t)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(t).
				Header("User-Agent", "some-agent/1.0").
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path", nil)
		require.NoError(t, err)
		r.Header.Set("User-Agent", "some-agent/1.0")

		_, err = mc.Do(context.Background(), io.Discard, r)
		require.NoError(t, err)

		mc.AssertExpectations(t)
	})

	t.Run("doesn't match wrong method", func(t *testing.T) {
		tm := new(testing.T)
		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(tm).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("POST", "http://domain.test/path",
			strings.NewReader(`{"foo":"1","bar":{"baz":["2","3"]}}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1")
		r.Header.Set("Bar", "2")

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		mc.AssertExpectations(tm)

		assert.True(t, tm.Failed())
	})

	t.Run("doesn't match wrong url", func(t *testing.T) {
		tm := new(testing.T)
		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(tm).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/other-path",
			strings.NewReader(`{"foo":"1","bar":{"baz":["2","3"]}}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1")
		r.Header.Set("Bar", "2")

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		mc.AssertExpectations(tm)

		assert.True(t, tm.Failed())
	})

	t.Run("doesn't match wrong header", func(t *testing.T) {
		tm := new(testing.T)
		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(tm).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`{"foo":"1","bar":{"baz":["2","3"]}}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1000")
		r.Header.Set("Bar", "2")

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		mc.AssertExpectations(tm)

		assert.True(t, tm.Failed())
	})

	t.Run("doesn't match wrong body", func(t *testing.T) {
		tm := new(testing.T)
		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(tm).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/path",
			strings.NewReader(`{"foo":"1","bar":{"qux":["2","3"]}}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "1")
		r.Header.Set("Bar", "2")

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		mc.AssertExpectations(tm)

		assert.True(t, tm.Failed())
	})

	t.Run("doesn't match malformed json", func(t *testing.T) {
		t.Run("expected", func(t *testing.T) {
			tm := new(testing.T)
			mc := NewMockClient(tm)

			assert.PanicsWithValue(t,
				"httputils.Match(): parse expected json: invalid character 'm' looking for beginning of value",
				func() {
					mc.On("Do", mock.Anything,
						Mock("PUT http://domain.test/path").
							Test(tm).
							Header("Foo", "1").
							Header("Bar", "2").
							Body(`malformed json`).
							Match(),
					).Return(200, nil, nil).Once()
				},
			)
		})

		t.Run("actual", func(t *testing.T) {
			tm := new(testing.T)
			mc := NewMockClient(tm)
			mc.On("Do", mock.Anything,
				Mock("PUT http://domain.test/path").
					Test(tm).
					Header("Foo", "1").
					Header("Bar", "2").
					Body(`{
						"foo": "1",
						"bar": {
							"baz": ["2", "3"]
						}
					}`).
					Match(),
			).Return(200, nil, nil).Once()

			r, err := http.NewRequest("PUT", "http://domain.test/path",
				strings.NewReader(`malformed json`))
			require.NoError(t, err)
			r.Header.Set("Foo", "1")
			r.Header.Set("Bar", "2")

			// Use a goroutine because mock calls Goexit when no match.
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = mc.Do(context.Background(), io.Discard, r)
			}()
			wg.Wait()

			mc.AssertExpectations(tm)

			assert.True(t, tm.Failed())
		})
	})

	t.Run("logs diff when no match", func(t *testing.T) {
		expected := `--- expected
+++ actual
@@ -1,15 +1,16 @@
-PUT /path HTTP/1.1
+PUT /other-path HTTP/1.1
 Host: domain.test
-Content-Length: 72
+Content-Length: 86
 Bar: 2
-Foo: 1
+Foo: 0
 
 {
   "bar": {
     "baz": [
       "2",
-      "3"
+      "4"
     ]
   },
-  "foo": "1"
+  "foo": "1",
+  "qux": "5"
 }
`

		cleanupFuncs := []func(){}

		tm := &tMock{}
		tm.Test(t)

		tm.On("Cleanup", mock.Anything).Run(func(args mock.Arguments) {
			cleanupFuncs = append(cleanupFuncs, args.Get(0).(func()))
		}).Return()

		tm.On("Logf", "httputils unmatched expectation:\n%s", []interface{}{expected}).
			Return().Once()

		tm.On("Logf", mock.Anything, mock.Anything).Return()

		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("PUT http://domain.test/path").
				Test(tm).
				Header("Foo", "1").
				Header("Bar", "2").
				Body(`{
					"foo": "1",
					"bar": {
						"baz": ["2", "3"]
					}
				}`).
				Match(),
		).Return(200, nil, nil).Once()

		r, err := http.NewRequest("PUT", "http://domain.test/other-path",
			strings.NewReader(`{"foo":"1","bar":{"baz":["2","4"]},"qux":"5"}`))
		require.NoError(t, err)
		r.Header.Set("Foo", "0")
		r.Header.Set("Bar", "2")

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}

		mc.AssertExpectations(tm)
		tm.AssertExpectations(t)
	})

	t.Run("logs diff of query params when no match", func(t *testing.T) {
		expected := `--- expected
+++ actual
@@ -1,4 +1,4 @@
-GET /path?bar=2&baz=3&foo=1&qux=4&qux=5 HTTP/1.1
+GET /path?bar=2%262&baz=3&foo=1&qux=5 HTTP/1.1
 Host: domain.test
 
 
=== query parameters:
--- expected
+++ actual
@@ -1,5 +1,4 @@
-bar=2
+bar=2&2
 baz=3
 foo=1
-qux=4
 qux=5
`

		cleanupFuncs := []func(){}

		tm := &tMock{}
		tm.Test(t)

		tm.On("Cleanup", mock.Anything).Run(func(args mock.Arguments) {
			cleanupFuncs = append(cleanupFuncs, args.Get(0).(func()))
		}).Return()

		tm.On("Logf", "httputils unmatched expectation:\n%s", []interface{}{expected}).
			Return().Once()

		tm.On("Logf", mock.Anything, mock.Anything).Return()

		expectedQ := url.Values{}
		expectedQ.Add("foo", "1")
		expectedQ.Add("bar", "2")
		expectedQ.Add("baz", "3")
		expectedQ.Add("qux", "4")
		expectedQ.Add("qux", "5")

		mc := NewMockClient(tm)
		mc.On("Do", mock.Anything,
			Mock("GET http://domain.test/path?"+expectedQ.Encode()).
				Test(tm).
				Match(),
		).Return(200, nil, nil).Once()

		actualQ := url.Values{}
		actualQ.Add("foo", "1")
		actualQ.Add("bar", "2&2")
		actualQ.Add("baz", "3")
		actualQ.Add("qux", "5")

		r, err := http.NewRequest("GET", "http://domain.test/path?"+actualQ.Encode(), nil)
		require.NoError(t, err)

		// Use a goroutine because mock calls Goexit when no match.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = mc.Do(context.Background(), io.Discard, r)
		}()
		wg.Wait()

		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}

		mc.AssertExpectations(tm)
		tm.AssertExpectations(t)
	})
}
