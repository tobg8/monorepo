package graceful

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_graceful(t *testing.T) {
	type ctxKey string

	const ctxKeyTest = ctxKey("test")

	t.Run("context key are preserved with 0 delay", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), ctxKeyTest, "ok")
		graceful, forceCancel := Graceful(ctx, 0)
		defer forceCancel()
		assert.Equal(t, "ok", graceful.Value(ctxKeyTest))
	})

	t.Run("context has no delay if timeout is 0", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		graceful, forceCancel := Graceful(ctx, 0)
		defer forceCancel()
		select {
		case <-graceful.Done():
		case <-time.After(time.Microsecond):
			t.Fatal("Graceful context not canceled directly")
		}

	})

	t.Run("give some time to graceful stop when context is canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())

		graceful, forceCancel := Graceful(ctx, 10*time.Second)
		defer forceCancel()
		cancel()
		select {
		case <-graceful.Done():
			t.Fatal("Graceful context not canceled directly")
		case <-time.After(time.Microsecond):
		}

	})

	t.Run("graceful context expires after provided duration when original context is done due to canceled reason", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())

		graceful, forceCancel := Graceful(ctx, 1*time.Second)
		defer forceCancel()
		now := time.Now()
		cancel()
		<-graceful.Done()
		assert.Truef(t, time.Since(now) > time.Second, "the graceful didn't delay properly")
	})

	t.Run("graceful context expires instantly when original context is done for other reason", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()

		graceful, forceCancel := Graceful(ctx, 1*time.Second)
		defer forceCancel()
		now := time.Now()
		<-graceful.Done()
		timeAdded := time.Since(now)
		assert.Truef(t, timeAdded < time.Second, "the graceful added too much delay (%s)", timeAdded)
	})
}
