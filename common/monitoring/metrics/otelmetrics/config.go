package otelmetrics

// Config is the configuration of a OTEL meter provider.
type Config struct {
	ServiceName string
	Version     string
	Env         string
}
