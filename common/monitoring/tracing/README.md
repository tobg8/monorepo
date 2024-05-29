# Tracing

Tracing is automatically configured when using the `common/application` library.

By default, traces are sampled using the Datadog's priority sampler, which
[applies some rules on selection of traces to ingest](https://docs.datadoghq.com/tracing/trace_pipeline/ingestion_mechanisms/).
We can however set a percentage of traces to ingest, as described in the
[Configuration](../README.md#Configuration) section.

The main function we will use with tracing is:

- `AddTag(ctx context.Context, key string, value interface{})`

This function adds a tag to a tracing span for a context, which is then visible on Datadog.

For example, the following snippet adds a `message` tag with the value of the error to the context of a client request:

```go
tracing.AddTag(c.Request.Context(), "message", err.Error())
```
