package httptester

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Head calls the specified url with the HEAD method. Body and headers are optional.
//
//	response := Head(app, "/url")
//
//	response := Head(app, "/url", strings.NewReader("request body"))
//
//	response := Head(app, "/url", map[string]string{"header name": "header value"})
//
//	response := Head(app, "/url", strings.NewReader("request body"), map[string]string{"header name": "header value"})
func Head(app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return HeadWithContext(context.Background(), app, url, bodyAndHeaders...)
}

// HeadWithContext calls the specified url with the HEAD method and a context. Body and headers are optional.
func HeadWithContext(ctx context.Context, app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodHead, bodyAndHeaders...)
}

// Get calls the specified url with the GET method. Body and headers are optional.
//
//	response := Get(app, "/url")
//
//	response := Get(app, "/url", strings.NewReader("request body"))
//
//	response := Get(app, "/url", map[string]string{"header name": "header value"})
//
//	response := Get(app, "/url", strings.NewReader("request body"), map[string]string{"header name": "header value"})
func Get(app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return GetWithContext(context.Background(), app, url, bodyAndHeaders...)
}

// GetWithContext calls the specified url with the GET method and a context. Body and headers are optional.
func GetWithContext(ctx context.Context, app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodGet, bodyAndHeaders...)
}

// Options calls the specified url with the OPTIONS method. Body and headers are optional.
//
//	response := Options(app, "/url")
//
//	response := Options(app, "/url", strings.NewReader("request body"))
//
//	response := Options(app, "/url", map[string]string{"header name": "header value"})
//
//	response := Options(app, "/url", strings.NewReader("request body"), map[string]string{"header name": "header value"})
func Options(app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return OptionsWithContext(context.Background(), app, url, bodyAndHeaders...)
}

// OptionsWithContext calls the specified url with the OPTIONS method and a context. Body and headers are optional.
func OptionsWithContext(ctx context.Context, app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodOptions, bodyAndHeaders...)
}

// Post calls the specified url with the POST method. Headers (as a single map passed in parameter) are optional.
// The content type header is set to application/json, if not already specified in the parameters.
// This is because we want to support json in almost all cases.
//
//	response := Post(app, "/url", "request body")
//
//	response := Post(app, "/url", "request body", map[string]string{"header name": "header value"})
func Post(app http.Handler, url string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	return PostWithContext(context.Background(), app, url, body, headers...)
}

// PostWithContext calls the specified url with the POST method and context. Headers (as a single map passed in parameter) are optional.
func PostWithContext(ctx context.Context, app http.Handler, url string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodPost, body, headersWithContentType(headers...))
}

// Put calls the specified url with the PUT method. Headers (as a single map passed in parameter) are optional.
// The content type header is set to application/json, if not already specified in the parameters.
// This is because we want to support json in almost all cases.
//
//	response := Put(app, "/url", "request body")
//
//	response := Put(app, "/url", "request body", map[string]string{"header name": "header value"})
func Put(app http.Handler, url string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	return PutWithContext(context.Background(), app, url, body, headers...)
}

// PutWithContext calls the specified url with the PUT method and context. Headers (as a single map passed in parameter) are optional.
func PutWithContext(ctx context.Context, app http.Handler, url string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodPut, body, headersWithContentType(headers...))
}

// Patch calls the specified url with the PATCH method. Headers (as a single map passed in parameter) are optional.
// The content type header is set to application/json, if not already specified in the parameters.
// This is because we want to support json in almost all cases.
//
//	response := Patch(app, "/url", "request body")
//
//	response := Patch(app, "/url", "request body", map[string]string{"header name": "header value"})
func Patch(app http.Handler, url string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	return PatchWithContext(context.Background(), app, url, body, headers...)
}

// PatchWithContext calls the specified url with the PATCH method and context. Headers (as a single map passed in parameter) are optional.
func PatchWithContext(ctx context.Context, app http.Handler, url string, body string, headers ...map[string]string) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodPatch, body, headersWithContentType(headers...))
}

// Delete calls the specified url with the DELETE method. Body and headers are optional.
//
//	response := Delete(app, "/url")
//
//	response := Delete(app, "/url", strings.NewReader("request body"))
//
//	response := Delete(app, "/url", map[string]string{"header name": "header value"})
//
//	response := Delete(app, "/url", strings.NewReader("request body"), map[string]string{"header name": "header value"})
func Delete(app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return DeleteWithContext(context.Background(), app, url, bodyAndHeaders...)
}

// DeleteWithContext calls the specified url with the DELETE method and context. Body and headers are optional.
func DeleteWithContext(ctx context.Context, app http.Handler, url string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	return do(ctx, app, url, http.MethodDelete, bodyAndHeaders...)
}

func do(ctx context.Context, app http.Handler, url string, method string, bodyAndHeaders ...any) *httptest.ResponseRecorder {
	var body string
	headers := make(map[string]string)
	for _, v := range bodyAndHeaders {
		switch t := v.(type) {
		default:
			fmt.Printf("unexpected type %T\n", t)
		case string:
			if body != "" {
				panic("unexpected second body")
			}
			body = t
		case map[string]string:
			if len(headers) > 0 {
				panic("unexpected second headers")
			}
			headers = t
		}
	}
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = strings.NewReader(body)
	}
	request := httptest.NewRequest(method, url, bodyReader).WithContext(ctx)
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)
	return recorder
}

func headersWithContentType(headers ...map[string]string) map[string]string {
	if len(headers) > 1 {
		panic("more than one map for the headers is not supported")
	}

	h := map[string]string{"Content-Type": "application/json"}
	if len(headers) == 1 {
		for k, v := range headers[0] {
			h[k] = v
		}
	}
	return h
}
