package loggingtest

import (
	"testing"

	"github.com/monorepo/common/logging"
)

type testingLogger struct {
	logger      logging.Logger
	loggerLevel logging.LoggerLevel
	tb          testing.TB
}

// ForTest returns a logging.Logger using t.Logf to print all logs.
func ForTest(tb testing.TB) logging.LoggerLevel {
	logger := logging.FromFuncHelper(tb.Log, tb.Helper)
	return &testingLogger{
		logger:      logger,
		loggerLevel: logger,
		tb:          tb,
	}
}

func (tl *testingLogger) SetLevel(lvl logging.Level) {
	tl.loggerLevel.SetLevel(lvl)
}

func (tl *testingLogger) GetLevel() logging.Level {
	return tl.loggerLevel.GetLevel()
}

func (tl *testingLogger) Debugf(format string, v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Debugf(format, v...)
}

func (tl *testingLogger) Infof(format string, v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Infof(format, v...)
}

func (tl *testingLogger) Warningf(format string, v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Warningf(format, v...)
}

func (tl *testingLogger) Errorf(format string, v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Errorf(format, v...)
}

func (tl *testingLogger) Debug(v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Debug(v...)
}

func (tl *testingLogger) Info(v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Info(v...)
}

func (tl *testingLogger) Warning(v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Warning(v...)
}

func (tl *testingLogger) Error(v ...interface{}) {
	tl.tb.Helper()
	tl.logger.Error(v...)
}

func (tl *testingLogger) WithField(fn string, fv interface{}) logging.Logger {
	return &testingLogger{
		logger: tl.logger.WithField(fn, fv),
		tb:     tl.tb,
	}
}

func (tl *testingLogger) WithFields(newFields logging.Fields) logging.Logger {
	return &testingLogger{
		logger: tl.logger.WithFields(newFields),
		tb:     tl.tb,
	}
}

func (tl *testingLogger) WithError(err error) logging.Logger {
	return tl.WithField("error", err)
}
