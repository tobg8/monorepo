package ddtracing

import (
	"context"
	"fmt"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/monorepo/common/logging"
)

// Service is the background service to manage the lifecycle of the Datadog tracer.
type Service struct {
	Config *Config
	Logger logging.Logger
}

// NewService creates a background service to start and stop the Datadog tracer.
func NewService(cfg *Config, logger logging.Logger) *Service {
	return &Service{
		Config: cfg,
		Logger: logger,
	}
}

// Run implements the `service.Service` interface.
//
// Run starts the global Datadog tracer.
func (svc *Service) Run() error {
	ddtrace.UseLogger(&ddLogger{
		logger: svc.Logger.WithField("lib", "ddtrace"),
	})

	opts := []tracer.StartOption{
		tracer.WithService(svc.Config.ServiceName),
		tracer.WithServiceVersion(svc.Config.Version),
		tracer.WithGlobalServiceName(true),
		tracer.WithEnv(svc.Config.Env),
		tracer.WithRuntimeMetrics(),
	}

	if svc.Config.RateSampled {
		if svc.Config.Rate > 1.0 || svc.Config.Rate < 0.0 {
			return fmt.Errorf("rate for rate sampler is invalid : %f", svc.Config.Rate)
		}

		opts = append(opts, tracer.WithSampler(
			tracer.NewRateSampler(svc.Config.Rate),
		))
	}

	tracer.Start(opts...)

	svc.Logger.Info("Datadog tracer started")

	return nil
}

// Stop implements the `service.Service` interface.
//
// Stop stops the global Datadog tracer.
func (svc *Service) Stop(_ context.Context) error {
	tracer.Stop()

	svc.Logger.Info("Datadog tracer stopped")

	return nil
}
