package slog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monorepo/common/logging"
	"github.com/monorepo/common/logging/loggingtest"
)

func TestNominal(t *testing.T) {
	var b bytes.Buffer

	l, err := New(&b, &logging.Config{
		Level: "debug",
	}, "test")
	require.NoError(t, err)

	l.
		WithField("f1", "v1").
		WithFields(logging.Fields{"f2": 42, "f3": []string{"a", "b", "c"}, "f4": nil}).
		WithError(errors.New("unexpected error")).
		Error("test")

	compareJSONLogs(t, b.String(), `{
	    "level": "error",
	    "message": "test",
	    "@version": "1",
	    "application": "test",
	    "f1": "v1",
	    "f3": ["a","b","c"],
	    "f4": null,
	    "f2": 42,
	    "error": "unexpected error"
	}`)
}

func TestDuplicateEntries(t *testing.T) {
	var b bytes.Buffer

	l, err := New(&b, &logging.Config{
		Level: "info",
	}, "test")
	require.NoError(t, err)

	l.
		WithField("foo", "value1").
		WithField("foo", "value2").
		Info("msg")

	compareJSONLogs(t, b.String(), `{
	    "@timestamp": "2024-02-21T15:28:15Z",
	    "level": "info",
	    "message": "msg",
	    "@version": "1",
	    "application": "test",
	    "foo": "value1",
	    "__duplicated_field": "This log entry contains duplicated fields"
	}`)
}

func compareJSONLogs(t *testing.T, got string, want string) {
	var actualMap map[string]any
	err := json.Unmarshal([]byte(got), &actualMap)
	require.NoError(t, err)

	require.NotEmpty(t, actualMap["@timestamp"])
	delete(actualMap, "@timestamp")

	var expectedMap map[string]any
	err = json.Unmarshal([]byte(want), &expectedMap)
	require.NoError(t, err)

	delete(expectedMap, "@timestamp")

	require.Equal(t, expectedMap, actualMap)
}

func TestWithFieldInParallel(t *testing.T) {
	const (
		sizeInParallel = 100
	)

	l, err := New(io.Discard, &logging.Config{
		Level: logging.LevelDebug.String(),
	}, "test")
	require.NoError(t, err)

	logger := l.WithField("user", "user-1")

	for i := 0; i < sizeInParallel; i++ {
		t.Run("test_"+fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			logger.WithField("ads", 50).
				WithField("field1", "value1").
				WithField("field2", "value2").
				WithField("field3", "value3").
				Info("test")

			logger.WithFields(logging.Fields{
				"field1": 1,
				"field2": 2,
				"field3": 3,
				"field4": 4,
				"field5": 5,
				"field6": 6,
				"field7": 7,
			}).Debug("test")
		})
	}
}

func TestSetLevel(t *testing.T) {
	type fields struct {
		levelVar *slog.LevelVar
	}
	type args struct {
		level logging.Level
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantLevel slog.Level
	}{
		{
			name: "setting to debug",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelInfo)
					return l
				}(),
			},
			args: args{
				level: logging.LevelDebug,
			},
			wantLevel: slog.LevelDebug,
		},
		{
			name: "setting to info",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelInfo)
					return l
				}(),
			},
			args: args{
				level: logging.LevelInfo,
			},
			wantLevel: slog.LevelInfo,
		},
		{
			name: "setting to warning",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelInfo)
					return l
				}(),
			},
			args: args{
				level: logging.LevelWarning,
			},
			wantLevel: slog.LevelWarn,
		},
		{
			name: "setting to error",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelInfo)
					return l
				}(),
			},
			args: args{
				level: logging.LevelError,
			},
			wantLevel: slog.LevelError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loggerSlog{
				levelVar: tt.fields.levelVar,
				rootLogger: slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{
					Level: tt.fields.levelVar,
				})),
			}
			l.SetLevel(tt.args.level)
			assert.Equal(t, tt.wantLevel, l.levelVar.Level())
		})
	}
}

func TestGetLevel(t *testing.T) {
	type fields struct {
		levelVar *slog.LevelVar
	}
	tests := []struct {
		name   string
		fields fields
		want   logging.Level
	}{
		{
			name: "getting debug",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelDebug)
					return l
				}(),
			},
			want: logging.LevelDebug,
		},
		{
			name: "getting info",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelInfo)
					return l
				}(),
			},
			want: logging.LevelInfo,
		},
		{
			name: "getting warning",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelWarn)
					return l
				}(),
			},
			want: logging.LevelWarning,
		},
		{
			name: "getting error",
			fields: fields{
				levelVar: func() *slog.LevelVar {
					l := new(slog.LevelVar)
					l.Set(slog.LevelError)
					return l
				}(),
			},
			want: logging.LevelError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loggerSlog{
				levelVar: tt.fields.levelVar,
				rootLogger: slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{
					Level: tt.fields.levelVar,
				})),
			}
			assert.Equal(t, tt.want, l.GetLevel())
		})
	}
}

func TestStdLogger(t *testing.T) {
	const (
		logMessage  = "Log entry from stdlib"
		slogMessage = "Slog entry from stdlib"
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
	rec.AssertLog(t, logging.Fields{
		"level":   "info",
		"message": logMessage,
		"std":     loggingtest.IgnoreFieldValue,
	})

	slog.Info(slogMessage)
	rec.AssertLog(t, logging.Fields{
		"level":   "info",
		"message": slogMessage,
		"std":     loggingtest.IgnoreFieldValue,
	})
}

func initTestLogger(tb testing.TB, lvl logging.Level) (*loggingtest.Recorder, logging.LoggerLevel) {
	const (
		appName = "test"
	)

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
