package interceptors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/f2prateek/train"
)

// Limiter is a HTTP client middleware to limit the number of concurrent requests
type Limiter chan struct{}

// NewLimiter instantiates a new Limiter
func NewLimiter(size int) Limiter {
	if size == 0 {
		return nil
	}
	return make(chan struct{}, size)
}
func (l Limiter) acquire(ctx context.Context) error {
	if l == nil {
		return nil
	}

	select {
	case l <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (l Limiter) release() {
	if l == nil {
		return
	}
	<-l
}

// Intercept implements train.Interceptor interface
func (l Limiter) Intercept(chain train.Chain) (*http.Response, error) {
	err := l.acquire(chain.Request().Context())
	if err != nil {
		return nil, fmt.Errorf("request interrupted while waiting in limit queue: %w", err)
	}

	defer l.release()
	return chain.Proceed(chain.Request())
}
