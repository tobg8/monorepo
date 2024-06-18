package httputils

import (
	"time"
)

// HTTPClient contains configuration for HTTP clients
type HTTPClient struct {
	Timeout   time.Duration `mapstructure:"timeout"`
	KeepAlive time.Duration `mapstructure:"keepalive"`
}
