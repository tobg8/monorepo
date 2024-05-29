package oteltracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otelsdk "go.opentelemetry.io/otel/sdk"
	otelsdkresource "go.opentelemetry.io/otel/sdk/resource"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/monitoring/semconv"
)

// NewTracerProvider instantiates a new OTEL tracer provider.
func NewTracerProvider(cfg *Config, logger logging.Logger) (trace.TracerProvider, func(), error) {
	const (
		otelExporterTracesEndpointEnvName = "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"
		otelExporterEndpointEnvName       = "OTEL_EXPORTER_OTLP_ENDPOINT"

		otelExporterTracesProtocolEnvName = "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"
		otelExporterProtocolEnvName       = "OTEL_EXPORTER_OTLP_PROTOCOL"
	)

	otelExporterOtlpEndpoint := os.Getenv(otelExporterTracesEndpointEnvName)
	if otelExporterOtlpEndpoint == "" {
		otelExporterOtlpEndpoint = os.Getenv(otelExporterEndpointEnvName)
	}

	otelExporterProtocol := os.Getenv(otelExporterTracesProtocolEnvName)
	if otelExporterProtocol == "" {
		otelExporterProtocol = os.Getenv(otelExporterProtocolEnvName)
	}

	if otelExporterProtocol == "" {
		if otelExporterOtlpEndpoint == "" {
			logger.Warningf(
				"Missing %s/%s environment variable, fallback to the console trace exporter",
				otelExporterTracesEndpointEnvName, otelExporterEndpointEnvName,
			)

			otelExporterProtocol = "console"
		} else {
			otelExporterProtocol = "http/protobuf"
		}
	}

	var (
		se  otelsdktrace.SpanExporter
		err error
	)

	switch otelExporterProtocol {
	case "http/protobuf":
		if otelExporterOtlpEndpoint == "" {
			return nil, nil, fmt.Errorf(
				"missing %s/%s environment variable",
				otelExporterTracesEndpointEnvName, otelExporterEndpointEnvName,
			)
		}

		se, err = otlptracehttp.New(context.Background(),
			otlptracehttp.WithEndpointURL(otelExporterOtlpEndpoint),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			otlptracehttp.WithTimeout(10*time.Second),
		)

	case "grpc":
		if otelExporterOtlpEndpoint == "" {
			return nil, nil, fmt.Errorf(
				"missing %s/%s environment variable",
				otelExporterTracesEndpointEnvName, otelExporterEndpointEnvName,
			)
		}

		se, err = otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithEndpointURL(otelExporterOtlpEndpoint),
			otlptracegrpc.WithTimeout(10*time.Second),
		)

	case "console":
		se, err = stdouttrace.New(
			// TODO(lbcjbb): Use a wrapper around the given logger.
			stdouttrace.WithWriter(os.Stderr),
			stdouttrace.WithPrettyPrint(),
		)

	default:
		return nil, nil, fmt.Errorf("unsupported otel exporter protocol %q", otelExporterProtocol)
	}

	if err != nil {
		return nil, nil, err
	}

	r := otelsdkresource.NewWithAttributes(semconv.SchemaURL,
		semconv.TelemetrySDKName("opentelemetry"),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersion(otelsdk.Version()),

		// semconv.ServiceNamespace("polaris.backend"),
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.Version),
		semconv.DeploymentEnvironment(cfg.Env),
	)

	tp := otelsdktrace.NewTracerProvider(
		otelsdktrace.WithResource(r),
		otelsdktrace.WithBatcher(se,
			otelsdktrace.WithBatchTimeout(5*time.Second),
			// otelsdktrace.WithBlocking(),
			otelsdktrace.WithExportTimeout(30*time.Second),
			otelsdktrace.WithMaxExportBatchSize(512),
			otelsdktrace.WithMaxQueueSize(2048),
		),
		otelsdktrace.WithSampler(otelsdktrace.AlwaysSample()),
	)

	return tp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			logger.
				WithError(err).
				Warning("can't shutdown the OTEL tracer provider")
		}
	}, nil
}
