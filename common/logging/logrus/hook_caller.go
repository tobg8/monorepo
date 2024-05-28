package logrus

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type callerHook struct{}

func (*callerHook) Levels() []logrus.Level { return logrus.AllLevels }

func (*callerHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 20)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return nil
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	for {
		f, more := frames.Next()

		if !strings.Contains(f.Function, "github.com/sirupsen/logrus") &&
			!strings.Contains(f.Function, "github.mpi-internal.com/leboncoin/go/common/logging") {
			entry.Data["caller"] = fmt.Sprintf("%s:%d", f.File, f.Line)
			break
		}

		if !more {
			break
		}
	}

	return nil
}
