package loggingtest

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/monorepo/common/logging"
)

// Mock is the mock implementation of a logging.Logger.
type Mock struct {
	mock.Mock
}

// NewMock creates a new mock implementation of a logging.Logger.
func NewMock(t *testing.T) *Mock {
	var m Mock

	m.Test(t)

	return &m
}

// SetLevel implements logging.Logger.
func (m *Mock) SetLevel(lvl logging.Level) {
	m.Called(lvl)
}

// GetLevel implements logging.Logger.
func (m *Mock) GetLevel() logging.Level {
	args := m.Called()
	return args.Get(0).(logging.Level)
}

// Debugf implements logging.Logger.
func (m *Mock) Debugf(format string, v ...interface{}) {
	m.Called(format, v)
}

// Infof implements logging.Logger.
func (m *Mock) Infof(format string, v ...interface{}) {
	m.Called(format, v)
}

// Warningf implements logging.Logger.
func (m *Mock) Warningf(format string, v ...interface{}) {
	m.Called(format, v)
}

// Errorf implements logging.Logger.
func (m *Mock) Errorf(format string, v ...interface{}) {
	m.Called(format, v)
}

// Debug implements logging.Logger.
func (m *Mock) Debug(v ...interface{}) {
	m.Called(v)
}

// Info implements logging.Logger.
func (m *Mock) Info(v ...interface{}) {
	m.Called(v)
}

// Warning implements logging.Logger.
func (m *Mock) Warning(v ...interface{}) {
	m.Called(v)
}

// Error implements logging.Logger.
func (m *Mock) Error(v ...interface{}) {
	m.Called(v)
}

// WithField implements logging.Logger.
func (m *Mock) WithField(fn string, fv interface{}) logging.Logger {
	args := m.Called(fn, fv)

	var logger logging.Logger
	if value := args.Get(0); value != nil {
		logger = value.(logging.Logger)
	}

	return logger
}

// WithFields implements logging.Logger.
func (m *Mock) WithFields(newFields logging.Fields) logging.Logger {
	args := m.Called(newFields)

	var logger logging.Logger
	if value := args.Get(0); value != nil {
		logger = value.(logging.Logger)
	}

	return logger
}

// WithError implements logging.Logger.
func (m *Mock) WithError(err error) logging.Logger {
	args := m.Called(err)

	var logger logging.Logger
	if value := args.Get(0); value != nil {
		logger = value.(logging.Logger)
	}

	return logger
}
