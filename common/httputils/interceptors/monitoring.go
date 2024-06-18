package interceptors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/monitoring/metrics"

	"github.com/monorepo/common/monitoring"
)

// Monitoring implements a monitoring interceptor
type Monitoring struct {
	Monitor       metrics.StatsdHandler
	RouteMatchers []RouteMatcher
}

// NewMonitoring instantiates a new Monitoring interceptor
func NewMonitoring(sh metrics.StatsdHandler, lrm ...RouteMatcher) *Monitoring {
	return &Monitoring{
		Monitor:       sh,
		RouteMatchers: lrm,
	}
}

// Intercept implements the train.Interceptor interface
func (m *Monitoring) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	start := time.Now()

	resp, err := chain.Proceed(req)

	var tags []string
	if err != nil {
		reason := "error"
		if err := req.Context().Err(); err != nil {
			reason = "context_deadline_exceeded"
			if err == context.Canceled {
				reason = "context_canceled"
			}
		}
		tags = append(tags, "request_failed:"+reason)
	} else {
		tags = append(tags, fmt.Sprintf("status_code:%d", resp.StatusCode))
		tags = append(tags, fmt.Sprintf("status_class:%dxx", resp.StatusCode/100))
	}

	tags = append(tags, "target:"+req.URL.Host)
	if m.RouteMatchers != nil {
		for _, rm := range m.RouteMatchers {
			if route, ok := rm.MatchRequest(req); ok {
				tags = append(tags, "route:"+route)
				break
			}
		}
	}

	for k, v := range monitoring.GetTagsFromContext(req.Context()) {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}

	m.Monitor.Timing("http.request.latency", time.Since(start), tags, 1)
	m.Monitor.Count("http.request.count", 1, tags, 1)

	return resp, err
}
