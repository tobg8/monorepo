# Profiling

Profiling cannot be programmatically enabled. It is enabled through environment variables, as described in
[this document](https://backstage.mpi-internal.com/docs/polaris/system/guild-backend/how-to/observability/profiling/).

The `ServiceProfiler` (defined in `profiler.go`) is configured with the following profiles:

- `HeapProfile`: reports memory usage
- `CPUProfile`: reports CPU usage
- `BlockProfile`: reports blockage on mutex and channel operations
- `MutexProfile`: reports lock contentions
- `GoroutineProfile`: reports stack traces of goroutines
- `MetricsProfile`: reports user metrics

WIP: additional information?
