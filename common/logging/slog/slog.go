package slog

import (
	"fmt"
	"io"
	"log/slog"
	"maps"
	"strings"
	"time"

	"github.com/monorepo/common/logging"
)

var (
	duplicateAttr = slog.String("__duplicated_field", "This log entry contains duplicated fields")
)

type loggerSlog struct {
	levelVar   *slog.LevelVar
	rootLogger *slog.Logger
	fieldNames map[string]struct{}
}

// New returns a new logger slog
func New(out io.Writer, config *logging.Config, appName string) (logging.LoggerLevel, error) {
	levelVar := new(slog.LevelVar)

	lvl, err := logging.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("can't parse the log level %q: %v", config.Level, err)
	}

	levelVar.Set(getSlogLevel(lvl))

	l := slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level: levelVar,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Key = "@timestamp"
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			case slog.LevelKey:
				v := strings.ToLower(a.Value.String())
				if v == "warn" {
					v = "warning"
				}
				a.Value = slog.StringValue(v)
			case slog.MessageKey:
				a.Key = "message"
			case slog.SourceKey:
				a.Key = "caller"
				a.Value = slog.StringValue(getCaller())
			}
			return a
		},
		AddSource: config.ShowCaller,
	}))

	return &loggerSlog{
		rootLogger: l.With(
			slog.String("@version", "1"),
			slog.String("application", appName),
		),
		levelVar:   levelVar,
		fieldNames: make(map[string]struct{}, 16),
	}, nil
}

// SetLevel set the log level accordingly
func (l *loggerSlog) SetLevel(level logging.Level) {
	l.levelVar.Set(getSlogLevel(level))
}

// GetLevel get the log level accordingly
func (l *loggerSlog) GetLevel() logging.Level {
	switch l.levelVar.Level() {
	case slog.LevelDebug:
		return logging.LevelDebug
	case slog.LevelInfo:
		return logging.LevelInfo
	case slog.LevelWarn:
		return logging.LevelWarning
	default:
		return logging.LevelError
	}
}

// Debug wraps slog.Debug
func (l *loggerSlog) Debug(v ...interface{}) {
	l.rootLogger.Debug(fmt.Sprint(v...))
}

// Info wraps slog.Info
func (l *loggerSlog) Info(v ...interface{}) {
	l.rootLogger.Info(fmt.Sprint(v...))
}

// Warning wraps slog.Warn
func (l *loggerSlog) Warning(v ...interface{}) {
	l.rootLogger.Warn(fmt.Sprint(v...))
}

// Error wraps slog.Error
func (l *loggerSlog) Error(v ...interface{}) {
	l.rootLogger.Error(fmt.Sprint(v...))
}

// Debugf wraps slog.Debug
func (l *loggerSlog) Debugf(format string, v ...any) {
	l.rootLogger.Debug(fmt.Sprintf(format, v...))
}

// Infof wraps slog.Info
func (l *loggerSlog) Infof(format string, v ...any) {
	l.rootLogger.Info(fmt.Sprintf(format, v...))
}

// Warningf wraps slog.Warn
func (l *loggerSlog) Warningf(format string, v ...any) {
	l.rootLogger.Warn(fmt.Sprintf(format, v...))
}

// Errorf wraps slog.Error
func (l *loggerSlog) Errorf(format string, v ...any) {
	l.rootLogger.Error(fmt.Sprintf(format, v...))
}

func (l *loggerSlog) hasDuplicateAttr() bool {
	_, found := l.fieldNames[duplicateAttr.Key]
	return found
}

// WithError wraps the error in a field
func (l *loggerSlog) WithError(err error) logging.Logger {
	return l.WithField("error", err)
}

// WithField wraps slog.With
func (l *loggerSlog) WithField(k string, v interface{}) logging.Logger {
	_, fieldFound := l.fieldNames[k]
	if fieldFound && l.hasDuplicateAttr() {
		return l
	}

	var attr slog.Attr
	if !fieldFound {
		attr = slog.Attr{
			Key:   k,
			Value: slog.AnyValue(v),
		}
	} else {
		attr = duplicateAttr
	}

	return l.withAttrs([]slog.Attr{attr})
}

// WithFields wraps slog.With
func (l *loggerSlog) WithFields(fields logging.Fields) logging.Logger {
	var (
		nAttrs           int
		addDuplicateAttr bool
	)

	for k := range fields {
		if _, fieldFound := l.fieldNames[k]; fieldFound {
			if !l.hasDuplicateAttr() {
				addDuplicateAttr = true
				nAttrs++
			}
			continue
		}

		nAttrs++
	}

	var attrs = make([]slog.Attr, 0, nAttrs)
	for k, v := range fields {
		if _, ok := l.fieldNames[k]; ok {
			continue
		}

		attrs = append(attrs, slog.Attr{
			Key:   k,
			Value: slog.AnyValue(v),
		})
	}

	if addDuplicateAttr {
		attrs = append(attrs, duplicateAttr)
	}

	return l.withAttrs(attrs)
}

// SlogLogger returns a new slog.Logger. This logger acts as a bridge from the log API.
func (l *loggerSlog) SlogLogger() *slog.Logger {
	return l.rootLogger
}

func (l *loggerSlog) withAttrs(attrs []slog.Attr) *loggerSlog {
	if len(attrs) == 0 {
		return l
	}

	newFieldNames := make(map[string]struct{}, len(l.fieldNames)+len(attrs))
	maps.Copy(newFieldNames, l.fieldNames)

	for idx := range attrs {
		attr := &attrs[idx]
		newFieldNames[attr.Key] = struct{}{}
	}

	h := l.rootLogger.Handler()
	h = h.WithAttrs(attrs)

	return &loggerSlog{
		rootLogger: slog.New(h),
		levelVar:   l.levelVar,
		fieldNames: newFieldNames,
	}
}

// getSlogLevel returns the slog level from the logging level
func getSlogLevel(lvl logging.Level) slog.Level {
	switch lvl {
	case logging.LevelDebug:
		return slog.LevelDebug
	case logging.LevelInfo:
		return slog.LevelInfo
	case logging.LevelWarning:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}
