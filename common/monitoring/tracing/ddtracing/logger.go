package ddtracing

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/monorepo/common/logging"
)

var (
	// The log pattern is retrieved from the internal/log package in dd-trace-go.
	// https://github.com/DataDog/dd-trace-go/blob/v1.58.0/internal/log/log.go#L194
	// https://github.com/DataDog/dd-trace-go/blob/v1.58.0/internal/log/log.go#L31
	ddLogRegexp = regexp.MustCompile(`^Datadog Tracer v[0-9.]+ (DEBUG|INFO|WARN|ERROR): (.+)\n?$`)
)

const (
	// https://github.com/DataDog/dd-trace-go/blob/v1.58.0/ddtrace/tracer/log.go#L135
	ddTracerConfigPrefix = "DATADOG TRACER CONFIGURATION"
)

type ddLogger struct {
	logger logging.Logger
}

func (l ddLogger) Log(msg string) {
	m := ddLogRegexp.FindStringSubmatch(msg)
	if len(m) == 0 {
		l.logger.Debug(msg)
		return
	}

	lvl, ddMsg := m[1], m[2]

	var logFunc func(a ...any)
	switch lvl {
	case "DEBUG":
		logFunc = l.logger.Debug
	case "INFO":
		// Don't show the tracer configuration in info level to avoid to reduce noice in logs.
		// And pretty print the json configuration of the tracer if possible.
		if strings.HasPrefix(ddMsg, ddTracerConfigPrefix+" ") {
			ddMsg = strings.TrimPrefix(ddMsg, ddTracerConfigPrefix+" ")

			var ddConf any
			if err := json.Unmarshal([]byte(ddMsg), &ddConf); err == nil {
				logFunc = l.logger.WithField("dd_tracer_config", ddConf).Debug
				ddMsg = ddTracerConfigPrefix
			} else {
				logFunc = l.logger.Debug
			}
		} else {
			logFunc = l.logger.Info
		}
	case "WARN":
		logFunc = l.logger.Warning
	case "ERROR":
		logFunc = l.logger.Error
	default:
		logFunc = l.logger.Debug
	}

	logFunc(ddMsg)
}
