package retrierx

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/stretchr/testify/assert"
)

func TestLoggingClassifier_Classify(t *testing.T) {
	messages := make([]string, 0)
	f := func(format string, v ...interface{}) {
		messages = append(messages, fmt.Sprintf(format, v...))
	}
	c := LoggingClassifier{LogFunc: f}

	res := c.Classify(errors.New("foo"))
	assert.Equal(t, retrier.Retry, res)
	assert.Len(t, messages, 1)
	assert.Equal(t, messages[0], "foo")

	res = c.Classify(nil)
	assert.Equal(t, retrier.Succeed, res)
	assert.Len(t, messages, 1)
}

func TestLoggingClassifier_ClassifyNilFunction(t *testing.T) {
	c := LoggingClassifier{}
	res := c.Classify(errors.New("foo"))
	assert.Equal(t, retrier.Retry, res)
	res = c.Classify(nil)
	assert.Equal(t, retrier.Succeed, res)
}

func TestLinearBackoff(t *testing.T) {
	backoff := LinearBackoff(4, 10)
	assert.EqualValues(t, []time.Duration{10, 20, 30, 40}, backoff)
}
