package otelmetrics

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	otelsdk "go.opentelemetry.io/otel/sdk"
	otelsdkmetric "go.opentelemetry.io/otel/sdk/metric"
	otelsdkresource "go.opentelemetry.io/otel/sdk/resource"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/monitoring/semconv"
)

// NewMeterProvider instantiates a new OTEL meter provider.
func NewMeterProvider(cfg *Config, logger logging.Logger) (metric.MeterProvider, func(), error) {
	// TODO(lbcjbb): Move in the common/application library.
	SetGlobalErrorHandler(logger)

	const (
		otelExporterMetricsEndpointEnvName = "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"
		otelExporterEndpointEnvName        = "OTEL_EXPORTER_OTLP_ENDPOINT"

		otelExporterMetricsProtocolEnvName = "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"
		otelExporterProtocolEnvName        = "OTEL_EXPORTER_OTLP_PROTOCOL"
	)

	otelExporterOtlpEndpoint := os.Getenv(otelExporterMetricsEndpointEnvName)
	if otelExporterOtlpEndpoint == "" {
		otelExporterOtlpEndpoint = os.Getenv(otelExporterEndpointEnvName)
	}

	otelExporterProtocol := os.Getenv(otelExporterMetricsProtocolEnvName)
	if otelExporterProtocol == "" {
		otelExporterProtocol = os.Getenv(otelExporterProtocolEnvName)
	}

	if otelExporterProtocol == "" {
		if otelExporterOtlpEndpoint == "" {
			logger.Warningf(
				"Missing %s/%s environment variable, fallback to the console metrics exporter",
				otelExporterMetricsEndpointEnvName, otelExporterEndpointEnvName,
			)

			otelExporterProtocol = "console"
		} else {
			otelExporterProtocol = "http/protobuf"
		}
	}

	var (
		me  otelsdkmetric.Exporter
		err error
	)

	switch otelExporterProtocol {
	case "http/protobuf":
		if otelExporterOtlpEndpoint == "" {
			return nil, nil, fmt.Errorf(
				"missing %s/%s environment variable",
				otelExporterMetricsEndpointEnvName, otelExporterEndpointEnvName,
			)
		}

		me, err = otlpmetrichttp.New(context.Background(),
			otlpmetrichttp.WithEndpointURL(otelExporterOtlpEndpoint),
			otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
			otlpmetrichttp.WithTimeout(10*time.Second),
		)

	case "grpc":
		if otelExporterOtlpEndpoint == "" {
			return nil, nil, fmt.Errorf(
				"missing %s/%s environment variable",
				otelExporterMetricsEndpointEnvName, otelExporterEndpointEnvName,
			)
		}

		me, err = otlpmetricgrpc.New(context.Background(),
			otlpmetricgrpc.WithEndpointURL(otelExporterOtlpEndpoint),
			otlpmetricgrpc.WithTimeout(10*time.Second),
		)

	case "console":
		me, err = stdoutmetric.New(
			// TODO(lbcjbb): Use a wrapper around the given logger.
			stdoutmetric.WithWriter(os.Stderr),
			stdoutmetric.WithPrettyPrint(),
		)

	default:
		return nil, nil, fmt.Errorf("unsupported otel exporter protocol %q", otelExporterProtocol)
	}

	if err != nil {
		return nil, nil, err
	}

	pr := otelsdkmetric.NewPeriodicReader(me,
		otelsdkmetric.WithTimeout(5*time.Second),
		otelsdkmetric.WithInterval(10*time.Second),
	)

	r := otelsdkresource.NewWithAttributes(semconv.SchemaURL,
		semconv.TelemetrySDKName("opentelemetry"),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersion(otelsdk.Version()),

		// semconv.ServiceNamespace("polaris.backend"),
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.Version),
		semconv.DeploymentEnvironment(cfg.Env),
	)

	mp := otelsdkmetric.NewMeterProvider(
		otelsdkmetric.WithReader(pr),
		otelsdkmetric.WithResource(r),
	)

	return mp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mp.Shutdown(ctx); err != nil {
			logger.
				WithError(err).
				Warning("can't shutdown the OTEL meter provider")
		}
	}, nil
}

// SetGlobalMeterProvider registers mp as the global `metric.MeterProvider`.
func SetGlobalMeterProvider(mp metric.MeterProvider) {
	otel.SetMeterProvider(mp)
}
