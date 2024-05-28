package slog

import (
	"fmt"
	"runtime"
	"strings"
)

// getCaller returns the caller of the logger
func getCaller() string {
	pc := make([]uintptr, 20)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return ""
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	for {
		f, more := frames.Next()

		if !strings.Contains(f.Function, "log/slog") &&
			!strings.Contains(f.Function, "github.mpi-internal.com/leboncoin/go/common/logging") {
			return fmt.Sprintf("%s:%d", f.File, f.Line)
		}

		if !more {
			return ""
		}
	}
}
