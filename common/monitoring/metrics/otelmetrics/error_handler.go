package otelmetrics

import (
	"go.opentelemetry.io/otel"

	"github.com/monorepo/common/logging"
)

// SetGlobalErrorHandler registers a global `otel.ErrorHandler` to log all errors.
func SetGlobalErrorHandler(logger logging.Logger) {
	otel.SetErrorHandler(&otelErrorHandler{logger})
}

type otelErrorHandler struct {
	logger logging.Logger
}

func (h *otelErrorHandler) Handle(err error) {
	h.logger.WithError(err).Warning("otel error")
}
