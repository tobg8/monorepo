package profiling

// Config is the configuration of the Datadog profiler service.
type Config struct {
	ServiceName string
	Version     string
	Env         string
	Tags        []string
}
