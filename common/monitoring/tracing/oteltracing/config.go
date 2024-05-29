package oteltracing

// Config is the configuration of the OTEL tracer.
type Config struct {
	ServiceName string
	Version     string
	Env         string
}
