package glog

import (
	"log"

	"github.com/monorepo/common/logging"
)

var defaultLogger logging.LoggerLevel

func init() {
	defaultLogger = logging.FromFunc(log.Print)
}

// SetupLogger setup the logger used by the package methods.
func SetupLogger(logger logging.LoggerLevel) {
	defaultLogger = logger
}

// GetLogger returns the global logger.
func GetLogger() logging.LoggerLevel {
	return defaultLogger
}

// SetLevel sets the logger level.
func SetLevel(lvl logging.Level) {
	defaultLogger.SetLevel(lvl)
}

// GetLevel returns the logger level.
func GetLevel() logging.Level {
	return defaultLogger.GetLevel()
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Warningf logs a formatted warning message
func Warningf(format string, v ...interface{}) {
	defaultLogger.Warningf(format, v...)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Debug logs a formatted debug message
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Info logs a formatted info message
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Warning logs a formatted warning message
func Warning(v ...interface{}) {
	defaultLogger.Warning(v...)
}

// Error logs a formatted error message
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// WithField adds a key value field to the log message
func WithField(fn string, fv interface{}) logging.Logger {
	return defaultLogger.WithField(fn, fv)
}

// WithFields adds a map of fields to the log message
func WithFields(newFields logging.Fields) logging.Logger {
	return defaultLogger.WithFields(newFields)
}

// WithError adds an error as single field error to the log message
func WithError(err error) logging.Logger {
	return defaultLogger.WithError(err)
}
