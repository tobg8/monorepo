package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Float64ObservableCounter is an autogenerated mock type for the `Float64ObservableCounter` interface.
type Float64ObservableCounter struct {
	embedded.Float64ObservableCounter
	metric.Float64Observable
	mock.Mock
}

// NewFloat64ObservableCounter creates a new mock instance of `Float64ObservableCounter`.
// It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFloat64ObservableCounter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Float64ObservableCounter {
	m := new(Float64ObservableCounter)
	m.Test(t)

	t.Cleanup(func() {
		m.AssertExpectations(t)
	})

	return m
}