package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Float64ObservableUpDownCounter is an autogenerated mock type for the `Float64ObservableUpDownCounter` interface.
type Float64ObservableUpDownCounter struct {
	embedded.Float64ObservableUpDownCounter
	metric.Float64Observable
	mock.Mock
}

// NewFloat64ObservableUpDownCounter creates a new mock instance of `Float64ObservableUpDownCounter`.
// It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFloat64ObservableUpDownCounter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Float64ObservableUpDownCounter {
	m := new(Float64ObservableUpDownCounter)
	m.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}