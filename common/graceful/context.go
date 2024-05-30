package graceful

import (
	"context"
	"errors"
	"time"
)

// Graceful returns a context which keeps the original context
// but expand it's duration by a given delay if it's canceled
// (but not in other cases)
// Still provides a force cancel
func Graceful(parent context.Context, delay time.Duration) (graceful context.Context, forceCancel func()) {
	if delay <= 0 {
		return context.WithCancel(parent)
	}
	graceful, fc := context.WithCancel(context.WithoutCancel(parent))

	stop := context.AfterFunc(parent, func() {
		if errors.Is(parent.Err(), context.Canceled) {
			sleep, clean := context.WithTimeout(graceful, delay)
			defer clean()
			<-sleep.Done()
		}
		fc()
	})

	return graceful, func() {
		fc()
		stop()
	}
}
