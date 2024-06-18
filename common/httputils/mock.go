package httputils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/mock"
)

// Doer is the interface of Client.
type Doer interface {
	Do(ctx context.Context, dst io.Writer, request *http.Request) (statusCode int, err error)
	DoAndUnmarshalJSON(ctx context.Context, v interface{}, request *http.Request) (statusCode int, err error)
}

// MockClient is a Client mock. Mocked response body is third parameter of
// Mock.On() method; it may be an io.Reader, a string or a []byte. See examples.
//
// Note that using a Reader does not allow the mock to be used multiple times:
// it is read at first call and subsequent calls would get it empty (already
// read, not rewound).
type MockClient struct {
	mock.Mock
}

// Do implements and mocks Doer.
func (m *MockClient) Do(ctx context.Context, dst io.Writer, request *http.Request) (statusCode int, err error) {
	if dst == nil {
		panic(`destination cannot be nil, as this might panic when the response body is not empty; if you want to ignore the body, then pass io.Discard instead`)
	}

	args := m.Called(dst, request)

	body := args.Get(2)
	if body != nil {
		var src io.Reader
		switch v := body.(type) {
		case io.Reader:
			src = v
		case string:
			src = strings.NewReader(v)
		case []byte:
			src = bytes.NewReader(v)
		default:
			panic(fmt.Sprintf("httputils: MockClient.Do body: expected io.Reader, string or []byte, got %T", v))
		}
		if _, err := io.Copy(dst, src); err != nil {
			panic("httputils: MockClient.Do failed to copy body into destination Writer")
		}
	}

	return args.Int(0), args.Error(1)
}

// DoAndUnmarshalJSON implements and mocks JsonDoer
func (m *MockClient) DoAndUnmarshalJSON(ctx context.Context, v interface{}, request *http.Request) (statusCode int, err error) {
	args := m.Called(v, request)

	body := args.Get(2)
	if body != nil {
		var src io.Reader
		switch v := body.(type) {
		case io.Reader:
			src = v
		case string:
			src = strings.NewReader(v)
		case []byte:
			src = bytes.NewReader(v)
		default:
			panic(fmt.Sprintf("httputils: MockClient.DoAndUnmarshalJSON body: expected io.Reader, string or []byte, got %T", v))
		}
		decoder := json.NewDecoder(src)
		if err := decoder.Decode(v); err != nil {
			return args.Int(0), err
		}
	}

	return args.Int(0), args.Error(1)
}

const ignoreBody = "<ignore body>"

// Mock returns a matcher for mock.Mock that supports matching criteria on a request.
// The object must be converted to an argument matcher with Match() so that it can be used with mock.Mock.
func Mock(call string) *RequestMatcher {
	split := strings.Split(call, " ")
	if len(split) != 2 {
		panic(fmt.Sprintf("parameter should be with format METHOD URL, e.g. GET http://www.site.com/, was: %s", call))
	}
	method := split[0]
	url := split[1]
	return &RequestMatcher{
		method:  method,
		url:     url,
		body:    ignoreBody,
		rawBody: ignoreBody,
		headers: make(map[string]string, 0),
	}
}

// RequestMatcher holds the criteria that the request should match with
type RequestMatcher struct {
	t       testing.TB
	method  string
	url     string
	body    string
	rawBody string
	headers map[string]string
}

// Test sets the test object of the matcher.
//
// This is optional but you should use it: When set, the diff of expected and
// actual requests is printed when the test fails, but when not set, raw output
// from mock is not helpful.
func (m *RequestMatcher) Test(t testing.TB) *RequestMatcher {
	m.t = t
	return m
}

// Body ensures that the body of the request is the same json object as the one in parameter
func (m *RequestMatcher) Body(json string) *RequestMatcher {
	m.body = json
	return m
}

// RawBody ensures that the body of the request is the same json object as the one in parameter
func (m *RequestMatcher) RawBody(raw string) *RequestMatcher {
	m.rawBody = raw
	return m
}

// Header checks that one of the headers has the specified value
func (m *RequestMatcher) Header(key, value string) *RequestMatcher {
	m.headers[key] = value
	return m
}

// Match uses the data in RequestMatcher to create an argument matcher that
// works with mock.Mock. Outputs a readable diff on failure.
//
// NOTE: Content-Length present in diff is not the actual value but the value
// of the indented JSON generated for diff readability.
func (m *RequestMatcher) Match() interface{} {
	// Create a http.Request with RequestMatcher data with its JSON body in
	// diff-friendly format (key-sorted and indented), and write it as HTTP/1.1
	// string which is then used for matching and diff output.

	var expectedBody io.Reader = http.NoBody
	if m.body != ignoreBody {
		nj, err := normaliseJSON([]byte(m.body))
		if err != nil {
			panic(fmt.Sprintf("httputils.Match(): parse expected json: %s", err))
		}
		expectedBody = bytes.NewReader(nj)
	}

	if m.rawBody != ignoreBody {
		expectedBody = strings.NewReader(m.rawBody)
	}

	expectedReq, err := http.NewRequest(m.method, m.url, expectedBody)
	if err != nil {
		panic(fmt.Sprintf("httputils.Match(): new expected request: %s", err))
	}

	expectedReq.Header.Set("User-Agent", "") // avoid default user-agent
	for k, v := range m.headers {
		expectedReq.Header.Set(k, v)
	}

	buf := &bytes.Buffer{}
	err = expectedReq.Write(buf)
	if err != nil {
		panic(fmt.Sprintf("httputils.Match(): write expected request: %s", err))
	}
	expected := buf.String()

	// diffLogged is used to deduplicate diff outputs (the argument matcher is
	// called muliple times by mock).
	diffLogged := map[string]bool{}

	return mock.MatchedBy(func(req *http.Request) bool {
		// Copy actual request with its JSON body formatted to diff-friendly
		// format and its headers filtered to expectation subset, write it as
		// HTTP/1.1 string and compare it to expected.

		var actualBody io.Reader = http.NoBody
		if m.body != ignoreBody {
			b, err := io.ReadAll(req.Body)
			if err != nil {
				panic(fmt.Sprintf("httputils.Match(): read actual request body: %s", err))
			}
			req.Body = io.NopCloser(bytes.NewReader(b))
			nj, err := normaliseJSON(b)
			if err != nil {
				panic(fmt.Sprintf("httputils.Match(): parse actual json: %s", err))
			}
			actualBody = bytes.NewReader(nj)
		}

		if m.rawBody != ignoreBody {
			b, err := io.ReadAll(req.Body)
			if err != nil {
				panic(fmt.Sprintf("httputils.Match(): read actual request body: %s", err))
			}
			req.Body = io.NopCloser(bytes.NewReader(b))
			actualBody = bytes.NewReader(b)
		}

		actualReq, err := http.NewRequest(req.Method, req.URL.String(), actualBody)
		if err != nil {
			panic(fmt.Sprintf("httputils.Match(): new actual request copy: %s", err))
		}

		for k := range expectedReq.Header {
			for _, v := range req.Header.Values(k) {
				actualReq.Header.Add(k, v)
			}
		}
		if _, ok := actualReq.Header["User-Agent"]; !ok {
			actualReq.Header.Set("User-Agent", "") // avoid default user-agent
		}

		buf := &bytes.Buffer{}
		err = actualReq.Write(buf)
		if err != nil {
			panic(fmt.Sprintf("httputils.Match(): write actual request: %s", err))
		}
		actual := buf.String()

		match := expected == actual

		if !match && m.t != nil {
			m.t.Cleanup(func() {
				// Output readable diff if both match and test failed.

				if !m.t.Failed() {
					return
				}

				diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
					A:        difflib.SplitLines(strings.ReplaceAll(expected, "\r\n", "\n")),
					B:        difflib.SplitLines(strings.ReplaceAll(actual, "\r\n", "\n")),
					FromFile: "expected",
					ToFile:   "actual",
					Context:  3,
				})

				if expectedReq.URL.RawQuery != actualReq.URL.RawQuery {
					a := strings.Split(expectedReq.URL.RawQuery, "&")
					for i := range a {
						a[i], _ = url.QueryUnescape(a[i])
						a[i] += "\n"
					}
					b := strings.Split(actualReq.URL.RawQuery, "&")
					for i := range b {
						b[i], _ = url.QueryUnescape(b[i])
						b[i] += "\n"
					}
					qDiff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
						A:        a,
						B:        b,
						FromFile: "expected",
						ToFile:   "actual",
						Context:  3,
					})
					diff = fmt.Sprintf("%s=== query parameters:\n%s", diff, qDiff)
				}

				if diffLogged[diff] {
					return
				}
				diffLogged[diff] = true

				m.t.Logf("httputils unmatched expectation:\n%s", diff)
			})
		}

		return match
	})
}

// normaliseJSON transforms given JSON in key-sorted, indented format suitable
// to equivalence comparison and diff.
func normaliseJSON(j []byte) ([]byte, error) {
	var m interface{}
	err := json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}
	b, _ := json.MarshalIndent(m, "", "  ") // unmarshal then marshal doesn't err
	return b, nil
}

// NewMockClient creates a new instance of MockClient. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockClient(t testing.TB) *MockClient {
	m := &MockClient{}
	m.Mock.Test(t)

	t.Cleanup(func() { m.AssertExpectations(t) })

	return m
}
