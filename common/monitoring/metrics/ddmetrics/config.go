package ddmetrics

// Config is the configuration of a Datadog statsd client.
type Config struct {
	Namespace string
	Tags      []string
}
