package httputils

import (
	"net/http"
	"testing"

	"github.com/f2prateek/train"
	"github.com/monorepo/common/httputils/interceptors"
	"github.com/monorepo/common/monitoring/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Observe(t *testing.T) {
	trMock := &transportMock{}
	trMock.On("AppendInterceptors", mock.MatchedBy(func(is []train.Interceptor) bool {
		require.Equal(t, 1, len(is))
		_, ok := is[0].(*interceptors.Monitoring)
		assert.True(t, ok)
		return ok
	})).Once()
	trMock.On("PrependInterceptors", mock.MatchedBy(func(is []train.Interceptor) bool {
		require.Equal(t, 1, len(is))
		_, ok := is[0].(*interceptors.Tracing)
		assert.True(t, ok)
		return true
	})).Once()

	rp := NewReverseProxy(func(req *http.Request) {}, trMock)

	rp.Observe(metrics.NoopStatsdHandler)

	assert.True(t, rp.isMonitored)
	assert.True(t, rp.isTraced)
	trMock.AssertExpectations(t)
}

func Test_WithMonitor(t *testing.T) {
	trMock := &transportMock{}
	trMock.On("AppendInterceptors", mock.MatchedBy(func(is []train.Interceptor) bool {
		require.Equal(t, 1, len(is))
		_, ok := is[0].(*interceptors.Monitoring)
		assert.True(t, ok)
		return ok
	}))

	rp := NewReverseProxy(func(req *http.Request) {}, trMock)

	rp.WithMonitor(metrics.NoopStatsdHandler)

	assert.True(t, rp.isMonitored)
	assert.False(t, rp.isTraced)
	trMock.AssertExpectations(t)
}

func Test_WithTracer(t *testing.T) {
	trMock := &transportMock{}
	trMock.On("PrependInterceptors", mock.MatchedBy(func(is []train.Interceptor) bool {
		require.Equal(t, 1, len(is))
		_, ok := is[0].(*interceptors.Tracing)
		assert.True(t, ok)
		return true
	})).Once()

	rp := NewReverseProxy(func(req *http.Request) {}, trMock)

	rp.WithTracer()

	assert.False(t, rp.isMonitored)
	assert.True(t, rp.isTraced)
	trMock.AssertExpectations(t)
}

type transportMock struct {
	mock.Mock
}

func (m *transportMock) AppendInterceptors(is ...train.Interceptor) {
	m.Called(is)
}

func (m *transportMock) PrependInterceptors(is ...train.Interceptor) {
	m.Called(is)
}

func (m *transportMock) RoundTrip(req *http.Request) (*http.Response, error) {
	called := m.Called(req)
	return called.Get(0).(*http.Response), called.Error(1)
}
