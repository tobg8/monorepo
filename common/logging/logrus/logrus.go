package logrus

import (
	"fmt"
	"io"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/monorepo/common/logging"
)

// ErrorKey defines the key when adding errors using WithError.
const ErrorKey = "error"

// New creates a Logger using a new logrus logger as backend.
func New(output io.Writer, cfg *logging.Config, appName string) (logging.LoggerLevel, error) {
	logger := logrus.New()
	logger.SetOutput(output)

	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("can't parse the log level %q: %v", cfg.Level, err)
	}

	logger.SetLevel(lvl)

	if appName != "" {
		logger.AddHook(&appNameHook{appName})
	}

	logger.AddHook(&versionHook{"1"})

	if cfg.ShowCaller {
		logger.AddHook(&callerHook{})
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	return &loggerLogrus{
		Entry:      logrus.NewEntry(logger),
		rootLogger: logger,
	}, nil
}

// loggerLogrus is the logrus implementation of the Logger interface.
type loggerLogrus struct {
	*logrus.Entry
	rootLogger *logrus.Logger
}

// SetLevel sets the logger level.
func (l *loggerLogrus) SetLevel(lvl logging.Level) {
	l.rootLogger.SetLevel(getLogrusLevel(lvl))
}

// GetLevel returns the logger level.
func (l *loggerLogrus) GetLevel() logging.Level {
	switch l.rootLogger.GetLevel() {
	case logrus.TraceLevel, logrus.DebugLevel:
		return logging.LevelDebug
	case logrus.InfoLevel:
		return logging.LevelInfo
	case logrus.WarnLevel:
		return logging.LevelWarning
	default:
		return logging.LevelError
	}
}

// WithField adds a key value field to the log message
func (l *loggerLogrus) WithField(fn string, value interface{}) logging.Logger {
	newEntry := l.Entry.WithField(fn, value)
	return &loggerLogrus{
		Entry:      newEntry,
		rootLogger: l.rootLogger,
	}
}

// WithFields adds a key value field to the log message
func (l *loggerLogrus) WithFields(newFields logging.Fields) logging.Logger {
	newEntry := l.Entry.WithFields(logrus.Fields(newFields))
	return &loggerLogrus{
		Entry:      newEntry,
		rootLogger: l.rootLogger,
	}
}

// WithError adds an error as single field (using the key defined in ErrorKey).
func (l *loggerLogrus) WithError(err error) logging.Logger {
	return l.WithField(ErrorKey, err)
}

// LogLogger returns a new log.Logger. This logger acts as a bridge from the log API.
func (l *loggerLogrus) LogLogger() *log.Logger {
	wl := l.WriterLevel(logrus.ErrorLevel)
	return log.New(wl, "", 0)
}

func getLogrusLevel(lvl logging.Level) logrus.Level {
	var logrusLevel logrus.Level
	switch lvl {
	case logging.LevelDebug:
		logrusLevel = logrus.DebugLevel
	case logging.LevelInfo:
		logrusLevel = logrus.InfoLevel
	case logging.LevelWarning:
		logrusLevel = logrus.WarnLevel
	default:
		logrusLevel = logrus.ErrorLevel
	}

	return logrusLevel
}
