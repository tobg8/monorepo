package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Int64ObservableUpDownCounter is an autogenerated mock type for the `Int64ObservableUpDownCounter` interface.
type Int64ObservableUpDownCounter struct {
	embedded.Int64ObservableUpDownCounter
	metric.Int64Observable
	mock.Mock
}

// NewInt64ObservableUpDownCounter creates a new mock instance of `Int64ObservableUpDownCounter`.
// It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInt64ObservableUpDownCounter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Int64ObservableUpDownCounter {
	m := new(Int64ObservableUpDownCounter)
	m.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}
