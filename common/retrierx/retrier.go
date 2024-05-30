// Package retrierx provides extensions for the go-resiliency/retrier package
package retrierx

import (
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/monorepo/common/logging"
)

// LoggingClassifier is a retrier.Classifier that logs errors using the provided logging function.
type LoggingClassifier struct {
	LogFunc func(format string, v ...interface{})
}

// Classify implements the retrier.Classifier interface.
func (c LoggingClassifier) Classify(err error) retrier.Action {
	if err == nil {
		return retrier.Succeed
	}
	if c.LogFunc != nil {
		c.LogFunc(err.Error())
	}
	return retrier.Retry
}

// DebugClassifier is a retrier.Classifier that logs errors using the provided logrus.FieldLogger.
type DebugClassifier struct {
	Logger logging.Logger
}

// Classify implements the retrier.Classifier interface.
func (c DebugClassifier) Classify(err error) retrier.Action {
	if err == nil {
		return retrier.Succeed
	}
	if c.Logger != nil {
		c.Logger.WithError(err).Debug("retrying...")
	}
	return retrier.Retry
}

// LinearGenerator returns a generator function which return
// time.Duration which linearly increase at each call
func LinearGenerator(initialAmount time.Duration) func() time.Duration {
	currentDelay := initialAmount
	return func() time.Duration {
		tmp := currentDelay
		currentDelay += initialAmount
		return tmp
	}
}

// LinearBackoff generates a simple back-off strategy of retrying 'n' times, and increasing the amount of
// time waited by the initial amount each time (ex: 10, 20, 30, 40, 50).
// The maximum amount of time waited will be (n+1)*n*duration/2.
// For example, if n = 100 and initial duration = 5 ms, then it will timeout at (100+1)*100*5/2 ms, ie. ~25 secs.
func LinearBackoff(n int, initialAmount time.Duration) []time.Duration {
	ret := make([]time.Duration, n)
	generator := LinearGenerator(initialAmount)
	for i := range ret {
		ret[i] = generator()
	}
	return ret
}
