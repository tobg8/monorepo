package profiling

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/monitoring/internal/runtime"
)

// Service is the background service to manage the lifecycle of the Datadog profiler.
type Service struct {
	Config *Config
	Logger logging.Logger
}

// NewService creates a background service to start and stop the Datadog profiler.
func NewService(cfg *Config, logger logging.Logger) *Service {
	return &Service{
		Config: cfg,
		Logger: logger,
	}
}

// Run implements the `service.Service` interface.
//
// Run starts the global Datadog profiler.
func (svc *Service) Run() error {
	tags := svc.Config.Tags
	tags = append(tags, runtime.Tags()...)

	err := profiler.Start(
		profiler.WithService(svc.Config.ServiceName),
		profiler.WithEnv(svc.Config.Env),
		profiler.WithVersion(svc.Config.Version),
		profiler.WithTags(tags...),
		profiler.WithProfileTypes(
			profiler.HeapProfile,
			profiler.CPUProfile,
			profiler.BlockProfile,
			profiler.MutexProfile,
			profiler.GoroutineProfile,
			profiler.MetricsProfile,
		),
		profiler.WithLogStartup(false),
	)

	if err != nil {
		return err
	}

	svc.Logger.Info("Datadog profiler started")

	return nil
}

// Stop implements the `service.Service` interface.
//
// Stop stops the global Datadog profiler.
func (svc *Service) Stop(_ context.Context) error {
	profiler.Stop()

	svc.Logger.Info("Datadog profiler stopped")

	return nil
}
