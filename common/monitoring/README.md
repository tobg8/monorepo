# monitoring

`monitoring` is a convenience library which provides monitoring and tracing through Datadog.
It is built on top of Datadog libraries, such as [datadog-go](https://github.com/DataDog/datadog-go) and [dd-trace-go](https://github.com/DataDog/dd-trace-go).

## Global description

This library provides observability-related functions to use with Datadog:

- Monitoring (`statsd_handler*.go`) enables the creation and the tracking of custom metrics, which can be of multiple types (count, gauge, timing, ...).
- Tracing (`tracing.go`) is a more specific type of monitoring, which provides tracking of a request through multiple services.
- Profiling (`profiler.go`) lets the developer get low-level performance related metrics, such as CPU usage, memory consumption, ... However, profiling generates a certain amount of additional load, and is thus used sparingly.

More information about Observability with Datadog can be found at [the Backend documentation about Datadog](https://backstage.mpi-internal.com/docs/polaris/system/guild-backend/how-to/observability/datadog/).

This library also provides default handlers for Health checks (/health and /ping). These default handlers are useful for checking that the service is up and running.

## Configuration

Some parts (monitoring and tracing) of the library can be configured by using `datadog` flags in the configuration files of the service, in the same manner that we can enable or disable datadog depending on the environment.
The configuration flags (`struct DDConf`) are as follows:

| field name   | description                                          | possible values                                                                            |
|--------------|------------------------------------------------------|--------------------------------------------------------------------------------------------|
| enabled      | enables or disables the datadog agent                | true / false                                                                               |
| namespace    | sets the namespace for custom metrics                | string, by default set to the service name                                                 |
| version      | version of the service to monitor                    | string, automatically set                                                                  |
| address      | address of the DD daemon                             | string, automatically set                                                                  |
| buffer-size  | buffer pool sizes for sending metrics                | int, by default set to 100                                                                 |
| rate-sampled | enables or disables rate-sampling for tracing        | true / false, defines if the `rate` value is used                                          |
| rate         | custom sample rate used for tracing, if rate-sampled | between [0.0; 1.0], by default set to NaN => the priority sampler of DD takes the decision |
| tags         | monitoring tags to attach to monitoring events       | array of strings                                                                           |

We can provide these flags in the `conf/*.yaml` files.
For example, the following datadog configuration enables rate-sampling on tracing, sets the rate to 80% and sets the custom metrics namespace to `my-service`:

```yaml
datadog:
  enabled: true
  namespace: my-service
  rate-sampled: true
  rate: 0.8
```

The datadog client is automatically setup when using the `common/application` library, meaning that the developer can directly use exported functions of the `monitoring` package without worrying of the underlying structs.
This also means that Datadog's APMs (Application Performance Monitoring, not to be confused with custom metrics) are automatically available as long as the Datadog client is enabled in the configuration.

More details for usage are provided in the following sections.

## How to: monitoring w/ custom metrics

Sending custom metrics is done by using an `StatsdHandler` interface.
The monitoring lib provides 3 concrete `StatsdHandler` implementations:

- `statsd` from `datadog-go`, which is the standard Datadog client
- `LoggingStatsdHandler` which is `statsd` with additional debug logs
- `NoopStatsdHandler` which is a dummy client, mainly for usage in tests

By default, we don't need to manipulate these implementations directly.
Indeed, when using the `common/application` library for creating a service, a global `StatsdHandler` is initialized (`statsd_handler_global.go`) and is usable through exported functions of the `monitoring` package. This global `StatsdHandler` relies on the standard `statsd`.

Thus, we can directly use the following exported functions to send custom metrics:

- `Count(name string, value int64, tags []string, rate float64)`
- `Gauge(name string, value float64, tags []string, rate float64)`
- `Timing(name string, value time.Duration, tags []string, rate float64)`
- `Distribution(name string, value float64, tags []string, rate float64)`
- `Event(e *metrics.StatsdEvent)`

A description of each metric is available [here](https://backstage.mpi-internal.com/docs/polaris/system/guild-backend/backend-tour/09-datadog/#metrics).

The `Histogram()` metric should not be used as it is inaccurate: it is computed by each Datadog agent, resulting in an incorrect computation of the distribution. `Distribution()` should be used instead. 

For example, the following line sends a `Count` metric, with the name `argus.errors`, a value of `1`, a tag `reason:timeout` and a rate of `1`:

```go
monitoring.Count("argus.errors", 1, []string{"reason:timeout"}, 1)
```

WIP: add usage for all types of metrics, if possible

## References

- [Backend Guild's Observability documentation](https://backstage.mpi-internal.com/docs/polaris/system/guild-backend/how-to/observability/)
- [Observability Community of Practice's Slack channel](https://adevinta.slack.com/archives/C018JRKG40N)
