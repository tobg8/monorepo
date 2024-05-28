package logging_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/logging/logrus"
	"github.com/monorepo/common/logging/slog"
)

func TestLoggers(t *testing.T) {
	type test struct {
		name    string
		logFunc func(l logging.LoggerLevel)
	}

	tests := []test{
		{
			name: "info with fields",
			logFunc: func(l logging.LoggerLevel) {
				l.WithFields(logging.Fields{
					"user_id":    33,
					"account_id": "ACC1234",
				}).Info("logger: hello world", 42)
			},
		},
		{
			name: "error with error",
			logFunc: func(l logging.LoggerLevel) {
				l.WithError(fmt.Errorf("panic there is an error")).Error("something wrong happened")
			},
		},
		{
			name: "warning with formatf",
			logFunc: func(l logging.LoggerLevel) {
				l.Warningf("warning: %s", "dangerous things happened here")
			},
		},
		{
			name: "debug",
			logFunc: func(l logging.LoggerLevel) {
				l.Debug("shoudn't appear")
				l.SetLevel(logging.LevelDebug)
				l.Debug("should appear")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var logrusBuf, slogBuf bytes.Buffer

			cfg := logging.Config{
				Level: "info",
			}

			logrusLog, err := logrus.New(&logrusBuf, &cfg, "testlogs")
			require.NoError(t, err)
			test.logFunc(logrusLog)

			slogLog, err := slog.New(&slogBuf, &cfg, "testlogs")
			require.NoError(t, err)
			test.logFunc(slogLog)

			assert.JSONEq(t, logrusBuf.String(), slogBuf.String())
		})
	}
}
