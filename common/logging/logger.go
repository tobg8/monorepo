package logging

import (
	"log"
	"log/slog"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Logger is the interface we use at LBC to interact
// with the log implementation.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})

	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})

	WithField(fn string, fv interface{}) Logger
	WithFields(newFields Fields) Logger

	WithError(err error) Logger
}

// LoggerLevel is a logger with a configurable level.
type LoggerLevel interface {
	Logger

	SetLevel(lvl Level)
	GetLevel() Level
}

// OverrideDefaultStandardLogger uses the given logger for all log printed by
// the standard library logger.
func OverrideDefaultStandardLogger(logger Logger) {
	stdLogger := logger.WithField("std", "unhandled call to standard log package")

	sl, ok := stdLogger.(interface {
		SlogLogger() *slog.Logger
	})

	if ok {
		l := sl.SlogLogger()

		// slog.SetDefault() also defines the global log logger.
		slog.SetDefault(l)
	} else {
		ll, ok := stdLogger.(interface {
			LogLogger() *log.Logger
		})

		if ok {
			l := ll.LogLogger()
			log.SetPrefix(l.Prefix())
			log.SetFlags(l.Flags())
			log.SetOutput(l.Writer())
		}
	}
}
