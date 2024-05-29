package ddmetrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/monitoring/internal/runtime"
	"github.com/monorepo/common/monitoring/metrics"
)

// NewClient creates a new Datadog statsd client.
func NewClient(conf *Config, logger logging.Logger) (*Client, error) {
	if conf.Namespace == "" {
		return nil, fmt.Errorf("empty namespace in configuration")
	}

	tags := conf.Tags
	tags = append(tags, runtime.Tags()...)

	dd, err := statsd.New("",
		statsd.WithNamespace(conf.Namespace),
		statsd.WithTags(tags),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client:    dd,
		logger:    logger,
		namespace: conf.Namespace,
	}, nil
}

// Client is a Datadog statsd client.
//
// It implements the `monitoring.StatsdHandler` interface.
type Client struct {
	client    *statsd.Client
	logger    logging.Logger
	namespace string
}

// Gauge measures the value of a metric at a particular time.
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) {
	if err := c.client.Gauge(name, value, tags, rate); err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s event to datadog", name)
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds.
//
// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info
// (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing).
func (c *Client) Timing(name string, value time.Duration, tags []string, rate float64) {
	if err := c.client.Timing(name, value, tags, rate); err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s event to datadog", name)
	}
}

// Count tracks how many times something happened per second.
func (c *Client) Count(name string, value int64, tags []string, rate float64) {
	if err := c.client.Count(name, value, tags, rate); err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s event to datadog", name)
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) {
	if err := c.client.Histogram(name, value, tags, rate); err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s event to datadog", name)
	}
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func (c *Client) Distribution(name string, value float64, tags []string, rate float64) {
	if err := c.client.Distribution(name, value, tags, rate); err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s event to datadog", name)
	}
}

// ServiceCheck sends the provided ServiceCheck.
func (c *Client) ServiceCheck(sc *metrics.StatsdServiceCheck) {
	var status statsd.ServiceCheckStatus
	switch sc.Status {
	case metrics.ServiceCheckStatusOk:
		status = statsd.Ok
	case metrics.ServiceCheckStatusWarn:
		status = statsd.Warn
	case metrics.ServiceCheckStatusCritical:
		status = statsd.Critical
	default:
		status = statsd.Unknown
	}

	err := c.client.ServiceCheck(&statsd.ServiceCheck{
		Name:      sc.Name,
		Status:    status,
		Timestamp: sc.Timestamp,
		Hostname:  sc.Hostname,
		Message:   sc.Message,
		Tags:      sc.Tags,
	})

	if err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s check to datadog", sc.Name)
	}
}

// Event sends the provided Event.
func (c *Client) Event(e *metrics.StatsdEvent) {
	var alertType statsd.EventAlertType
	switch e.AlertType {
	case metrics.AlertTypeInfo:
		alertType = statsd.Info
	case metrics.AlertTypeError:
		alertType = statsd.Error
	case metrics.AlertTypeWarning:
		alertType = statsd.Warning
	case metrics.AlertTypeSuccess:
		alertType = statsd.Success
	}

	err := c.client.Event(&statsd.Event{
		Title:     e.Title,
		Text:      e.Text,
		Timestamp: e.Timestamp,
		Hostname:  e.Hostname,
		Priority:  statsd.Normal,
		AlertType: alertType,
		Tags:      e.Tags,
	})

	if err != nil {
		c.logger.WithError(err).Errorf("Failed to send %s event to datadog", e.Title)
	}
}

// GetNamespace returns the configured namespace of this client.
func (c *Client) GetNamespace() string {
	return c.namespace
}

// Close the client connection.
func (c *Client) Close() {
	if err := c.client.Close(); err != nil {
		c.logger.WithError(err).Errorf("Failed to close datadog client")
	}
}
