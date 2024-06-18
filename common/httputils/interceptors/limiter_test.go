package interceptors

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func try(ctx context.Context, f func()) {
	done := make(chan struct{})
	go func() {
		f()
		done <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		}
	}
}

func Test_limiter_acquire_zero_size(t *testing.T) {
	p := NewLimiter(0)
	assert.Nil(t, p)

	var acquired int
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		acquired++
	})
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		acquired++
	})

	assert.Equal(t, 2, acquired)
}

func Test_limiter_release_zero_size(t *testing.T) {
	p := NewLimiter(0)
	assert.Nil(t, p)

	var acquired int
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	try(ctx, func() {
		p.release()
		acquired++
	})
	try(ctx, func() {
		p.release()
		acquired++
	})

	assert.Equal(t, 2, acquired)
}

func Test_limiter_acquire_full(t *testing.T) {
	p := NewLimiter(2)
	assert.NotNil(t, p)

	var acquired int
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// fill up the limiter
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		acquired++
	})
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		acquired++
	})

	// This one should not work since we never released.
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		acquired++
	})

	assert.Equal(t, 2, acquired)
}

func Test_limiter_release(t *testing.T) {
	p := NewLimiter(2)
	assert.NotNil(t, p)

	var acquired int
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// fill up the limiter
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		defer p.release()
		acquired++
	})
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		defer p.release()
		acquired++
	})
	try(ctx, func() {
		err := p.acquire(context.Background())
		require.NoError(t, err)
		defer p.release()
		acquired++
	})

	assert.Equal(t, 3, acquired)
}

func Test_limiter_abort(t *testing.T) {
	p := NewLimiter(1)
	assert.NotNil(t, p)

	// Fill the concurrent requests
	_ = p.acquire(context.Background())

	// Add a request with a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := p.acquire(ctx)

	assert.EqualError(t, err, "context canceled")
}
