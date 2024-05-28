package logrus

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/logging/loggingtest"
)

func TestLogrus(t *testing.T) {
	t.Run("field-logger", func(t *testing.T) {
		rec, logger := initTestLogger(t, logging.LevelInfo)

		l := logger.WithFields(logging.Fields{
			"nil":        nil,
			"string":     "value",
			"nil_string": (*string)(nil),
			"int":        1,
			"nil_int":    (*int)(nil),
			"bool":       true,
			"nil_bool":   (*bool)(nil),
			"float":      32.12,
			"nil_float":  (*float32)(nil),
			"struct":     struct{}{},
			"nilStruct":  (*struct{})(nil),
			"array":      []int{1, 2},
			"nil_array":  ([]int)(nil),
			"map": map[string]bool{
				"true_value": true,
			},
			"nil_map": (map[string]bool)(nil),
		})
		l.Warningf("test %d", 1)
		l.Debug("debug log line")

		rec.AssertLog(t, logging.Fields{
			"level":      "warning",
			"message":    "test 1",
			"nil":        nil,
			"string":     "value",
			"nil_string": nil,
			"int":        float64(1),
			"nil_int":    nil,
			"bool":       true,
			"nil_bool":   nil,
			"float":      32.12,
			"nil_float":  nil,
			"struct":     map[string]interface{}{},
			"nilStruct":  nil,
			"array":      []interface{}{float64(1), float64(2)},
			"nil_array":  nil,
			"map": map[string]interface{}{
				"true_value": true,
			},
			"nil_map": nil,
		})
	})

	t.Run("format", func(t *testing.T) {
		rec, logger := initTestLogger(t, logging.LevelInfo)

		logger.Infof("Test format %d", 100)

		rec.AssertLog(t, logging.Fields{
			"level":   "info",
			"message": "Test format 100",
		})
	})

	t.Run("writer-level", func(t *testing.T) {
		rec, logger := initTestLogger(t, logging.LevelInfo)
		logging.OverrideDefaultStandardLogger(logger)

		log.Println("Test with std logger")

		deadline := time.Now().Add(5 * time.Second)
		for {
			if time.Now().After(deadline) {
				t.Fatal("timeout, the log buffer is always empty")
			}

			if rec.Len() > 0 {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		rec.AssertLog(t, logging.Fields{
			"level":   "error",
			"message": "Test with std logger",
			"std":     loggingtest.IgnoreFieldValue,
		})
	})

	t.Run("with-error", func(t *testing.T) {
		rec, logger := initTestLogger(t, logging.LevelInfo)

		testError := fmt.Errorf("foo: %w", errors.New("test error"))
		logger.
			WithError(testError).
			Error("test WithError()")

		rec.AssertLog(t, logging.Fields{
			"level":   "error",
			"message": "test WithError()",
			"error":   "foo: test error",
		})
	})

	t.Run("set-level", func(t *testing.T) {
		t.Run("using-logger", func(t *testing.T) {
			rec, logger := initTestLogger(t, logging.LevelInfo)
			require.Equal(t, logging.LevelInfo, logger.GetLevel())

			logger.SetLevel(logging.LevelError)
			require.Equal(t, logging.LevelError, logger.GetLevel())

			logger.Info("test")

			rec.AssertLogs(t, []logging.Fields{})
		})

		t.Run("using-entry", func(t *testing.T) {
			rec, logger := initTestLogger(t, logging.LevelInfo)
			l := logger.WithField("foo", 100)
			require.Equal(t, logging.LevelInfo, logger.GetLevel())

			logger.SetLevel(logging.LevelError)
			require.Equal(t, logging.LevelError, logger.GetLevel())

			l.Info("test")

			rec.AssertLogs(t, []logging.Fields{})
		})
	})

	t.Run("std-logger", func(t *testing.T) {
		const (
			logMessage = "Log entry from stdlib"
		)

		rec, l := initTestLogger(t, logging.LevelInfo)

		oldStdLogger := log.Default()
		t.Cleanup(func() {
			log.SetFlags(oldStdLogger.Flags())
			log.SetPrefix(oldStdLogger.Prefix())
			log.SetOutput(oldStdLogger.Writer())
		})

		logging.OverrideDefaultStandardLogger(l)

		log.Print(logMessage)

		// logrus use a goroutine to read asynchronously all log data.
		// So we need to manually delete the logger and force the garbage collector
		// to stop this goroutine and flush the log buffer.
		l = nil //nolint:ineffassign // This assignment is required to free the memory
		runtime.GC()

		rec.AssertLog(t, logging.Fields{
			"level":   "error",
			"message": logMessage,
			"std":     loggingtest.IgnoreFieldValue,
		})
	})
}

func initTestLogger(tb testing.TB, lvl logging.Level) (*loggingtest.Recorder, logging.LoggerLevel) {
	const appName = "test"

	rec := loggingtest.NewRecorder("@timestamp")
	rec.ExpectFields(logging.Fields{
		"@version":    "1",
		"application": appName,
	})

	l, err := New(rec.Writer(), &logging.Config{
		Level: lvl.String(),
	}, appName)
	require.NoError(tb, err)

	return rec, l
}
