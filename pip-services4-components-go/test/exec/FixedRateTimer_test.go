package test_exec

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/stretchr/testify/assert"
)

func TestTimerWithCallback(t *testing.T) {
	var counter int32

	timer := exec.NewFixedRateTimerFromCallback(
		func(ctx context.Context) {
			atomic.AddInt32(&counter, 1)
		},
		100, 0, 5,
	)

	ctx := context.Background()
	timer.Start(ctx)
	time.Sleep(time.Millisecond * 500)
	timer.Stop(ctx)

	assert.True(t, atomic.LoadInt32(&counter) > 3)
}

func TestTimerWithCancelCallback(t *testing.T) {
	var counter, counterCanceled int32

	timer := exec.NewFixedRateTimerFromCallback(
		func(ctx context.Context) {
			atomic.AddInt32(&counter, 1)
			select {
			case <-ctx.Done():
				atomic.AddInt32(&counterCanceled, 1)
				break
			}
		},
		100, 0, 5,
	)

	ctx, cancel := context.WithCancel(context.Background())
	timer.Start(ctx)
	time.Sleep(time.Millisecond * 500)
	cancel()
	timer.Stop(ctx)
	time.Sleep(time.Millisecond * 100)

	assert.True(t, atomic.LoadInt32(&counter) > 3)
	assert.True(t, atomic.LoadInt32(&counterCanceled) > 3)
}
