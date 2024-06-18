package interceptors

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// RouteMatcher defines a way to match a HTTP
//
//	request and to return an unique route identifier.
type RouteMatcher interface {
	MatchRequest(req *http.Request) (string, bool)
}

type routeMatcherWithFunc struct {
	method    string
	route     string
	matchFunc func(req *http.Request) bool
}

func (rmwf *routeMatcherWithFunc) MatchRequest(req *http.Request) (string, bool) {
	if rmwf.method != req.Method || !rmwf.matchFunc(req) {
		return "", false
	}

	return fmt.Sprintf("%s:%s", strings.ToLower(rmwf.method), rmwf.route), true
}

// StaticRouteMatcher creates a new RouteMatcher to match a static route.
// A static route is a route whithout any dynamic parts
// (e.g. /foo/bar, /users/ads).
func StaticRouteMatcher(method string, route string) RouteMatcher {
	return &routeMatcherWithFunc{
		method: method,
		route:  route,
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == route
		},
	}
}

// DynamicRouteMatcher creates a new RouteMatcher to match a dynamic route.
// A dynamic route is a route with at least one part which may vary between two
// requests (e.g. /stores/100/users/1b305b3d-6793-43d5-b447-dc619c5618c4).
// The pattern given is a raw regular expression which will be compiled.
func DynamicRouteMatcher(method string, route string, pattern *regexp.Regexp) RouteMatcher {
	return &routeMatcherWithFunc{
		method: method,
		route:  route,
		matchFunc: func(req *http.Request) bool {
			return pattern.MatchString(req.URL.Path)
		},
	}
}

type loggingRouteMatcher struct {
	logRequest func(*http.Request)
}

func (lrm *loggingRouteMatcher) MatchRequest(req *http.Request) (string, bool) {
	lrm.logRequest(req)
	return "", false
}

// LoggingRouteMatcher creates a new RouteMatcher to call
// the given function for each request.
// This route matcher can be used to log all requests not catched by previous
// registered route matchers.
// It useful to log information about all unexpected sent HTTP requests.
func LoggingRouteMatcher(logRequest func(*http.Request)) RouteMatcher {
	return &loggingRouteMatcher{
		logRequest: logRequest,
	}
}
