package ddtracing

// Config is the configuration of the Datadog tracer service.
type Config struct {
	ServiceName string
	Version     string
	Env         string
	RateSampled bool
	Rate        float64
}
