package benchmarks

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/monorepo/common/logging"
	ll "github.com/monorepo/common/logging/logrus"
	ls "github.com/monorepo/common/logging/slog"
)

type loggerBench struct {
	Name        string
	BuildLogger func(out io.Writer, config *logging.Config, appName string) (logging.Logger, error)
}

func BenchmarkLogger(b *testing.B) {
	for _, lb := range []loggerBench{
		{
			Name: "logrus",
			BuildLogger: func(out io.Writer, config *logging.Config, appName string) (logging.Logger, error) {
				return ll.New(out, config, appName)
			},
		},
		{
			Name: "slog",
			BuildLogger: func(out io.Writer, config *logging.Config, appName string) (logging.Logger, error) {
				return ls.New(out, config, appName)
			},
		},
	} {
		if globalBenchFlags.LoggerName == "" {
			b.Run(lb.Name, func(b *testing.B) {
				runLoggerBench(b, &lb)
			})
		} else if globalBenchFlags.LoggerName == lb.Name {
			runLoggerBench(b, &lb)
		}
	}
}

func runLoggerBench(b *testing.B, lb *loggerBench) {
	b.Run("WithFields", func(b *testing.B) {
		duplicateEntriesBenchmark(b, func(b *testing.B, makeKey makeKeyFunc) {
			logger, err := lb.BuildLogger(io.Discard, &logging.Config{
				Level:      "info",
				ShowCaller: false,
			}, "test")
			require.NoError(b, err)

			for n := 0; n < b.N; n++ {
				for i := 0; i < globalBenchFlags.LoopIterations; i++ {
					logger = logger.WithFields(logging.Fields{
						makeKey("user_id", n, i):   1,
						makeKey("ads_count", n, i): 200,
						makeKey("is_admin", n, i):  true,
						makeKey("input", n, i):     "lorem ipsum",
						makeKey("output", n, i):    "",
					})
				}
				logger.Info("hello world")
			}
		})
	})

	b.Run("WithField", func(b *testing.B) {
		duplicateEntriesBenchmark(b, func(b *testing.B, makeKey makeKeyFunc) {
			logger, err := lb.BuildLogger(io.Discard, &logging.Config{
				Level:      "info",
				ShowCaller: false,
			}, "test")
			require.NoError(b, err)

			for n := 0; n < b.N; n++ {
				for i := 0; i < globalBenchFlags.LoopIterations; i++ {
					logger = logger.
						WithField(makeKey("user_id", n, i), 1).
						WithField(makeKey("ads_count", n, i), 200).
						WithField(makeKey("is_admin", n, i), true).
						WithField(makeKey("input", n, i), "lorem ipsum").
						WithField(makeKey("output", n, i), "")
				}
				logger.Info("hello world")
			}
		})
	})

	b.Run("Info", func(b *testing.B) {
		b.Run("level=ok", func(b *testing.B) {
			logger, err := lb.BuildLogger(io.Discard, &logging.Config{
				Level:      "info",
				ShowCaller: false,
			}, "test")
			require.NoError(b, err)

			for n := 0; n < b.N; n++ {
				logger.Info("hello world")
				logger.Info("test", 1, []string{"a", "b", "c", "d", "e", "f", "g", "h"})
				logger.Infof("test %d", 1)
			}
		})

		// Check an insuffisant logging level avoid cpu and memory usage.
		b.Run("level=ko", func(b *testing.B) {
			logger, err := lb.BuildLogger(io.Discard, &logging.Config{
				Level:      "warning",
				ShowCaller: false,
			}, "test")
			require.NoError(b, err)

			for n := 0; n < b.N; n++ {
				logger.Info("hello world")
				logger.Info("test", 1, []string{"a", "b", "c", "d", "e", "f", "g", "h"})
				logger.Infof("test %d", 1)
			}
		})
	})
}

type makeKeyFunc func(s string, ids ...int) string

func duplicateEntriesBenchmark(b *testing.B, bf func(b *testing.B, makeKey makeKeyFunc)) {
	for _, sbb := range []struct {
		Name    string
		MakeKey makeKeyFunc
	}{
		{
			Name: "with-dupl",
			MakeKey: func(s string, _ ...int) string {
				return s
			},
		},
		{
			Name: "without-dupl",
			MakeKey: func(s string, ids ...int) string {
				var (
					formats = []string{"%s"}
					args    = []any{s}
				)

				for _, id := range ids {
					formats = append(formats, "%d")
					args = append(args, id)
				}

				return fmt.Sprintf(strings.Join(formats, "_"), args...)
			},
		},
	} {
		b.Run(sbb.Name, func(b *testing.B) {
			bf(b, sbb.MakeKey)
		})
	}
}
