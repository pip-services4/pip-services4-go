package exec

import (
	"context"
	"sync"
	"time"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/util"
)

// FixedRateTimer timer that is triggered in equal time intervals.
// It has summetric cross-language implementation and is often used by Pip.Services toolkit
// to perform periodic processing and cleanup in microservices.
//
//	see INotifiable
//	Example:
//		type MyComponent {
//			timer FixedRateTimer
//		}
//		...
//		func (mc* MyComponent) Open(ctx, context.Context) {
//			...
//			mc.timer = NewFixedRateTimerFromCallback(func(ctx context.Context){ this.cleanup }, 60000, 0, 5);
//			mc.timer.Start(ctx);
//			...
//		}
//		func (mc* MyComponent) Close(ctx, context.Context) {
//			...
//			mc.timer.Stop(ctx);
//			...
//		}
type FixedRateTimer struct {
	task        INotifiable
	callback    func(ctx context.Context)
	delay       int
	interval    int
	ticker      *time.Ticker
	mtx         sync.Mutex
	workerCount int
	exit        chan bool
}

// NewFixedRateTimer creates new instance of the timer and sets its values.
//
//	Returns: *FixedRateTimer
func NewFixedRateTimer() *FixedRateTimer {
	return &FixedRateTimer{
		workerCount: 5,
		exit:        make(chan bool),
	}
}

// NewFixedRateTimerFromCallback creates new instance of the timer and sets its values.
//
//	Parameters:
//		- callback func() callback function to call when timer is triggered.
//		- interval int an interval to trigger timer in milliseconds.
//		- delay int a delay before the first triggering in milliseconds.
//		- workerCount int a count of parallel running workers.
//	Returns: *FixedRateTimer
func NewFixedRateTimerFromCallback(callback func(ctx context.Context),
	interval int, delay int, workerCount int) *FixedRateTimer {

	return &FixedRateTimer{
		workerCount: workerCount,
		exit:        make(chan bool),
		callback:    callback,
		interval:    interval,
		delay:       delay,
	}
}

// NewFixedRateTimerFromTask creates new instance of the timer and sets its values.
//
//	Parameters:
//		- callback INotifiable Notifiable object to call when timer is triggered.
//		- interval int an interval to trigger timer in milliseconds.
//		- delay int a delay before the first triggering in milliseconds.
//		- workerCount int a count of parallel running workers.
//	Returns: *FixedRateTimer
func NewFixedRateTimerFromTask(task INotifiable,
	interval int, delay int, workerCount int) *FixedRateTimer {

	c := &FixedRateTimer{
		workerCount: 5,
		exit:        make(chan bool),
		interval:    interval,
		delay:       delay,
	}
	c.SetTask(task)
	return c
}

// Task gets the INotifiable object that receives notifications from this timer.
//
//	Returns: INotifiable the INotifiable object or null if it is not set.
func (c *FixedRateTimer) Task() INotifiable {
	return c.task
}

// SetTask sets a new INotifiable object to receive notifications from this timer.
//
//	Parameters: value INotifiable a INotifiable object to be triggered.
func (c *FixedRateTimer) SetTask(value INotifiable) {
	if c.IsStarted() {
		return
	}
	c.task = value
	c.callback = func(ctx context.Context) {
		c.task.Notify(util.ContextHelper.NewContextWithTraceId(ctx, "timer"), NewEmptyParameters())
	}
}

// WorkerCount gets the worker count.
//
//	Returns: int worker count.
func (c *FixedRateTimer) WorkerCount() int {
	return c.workerCount
}

// SetWorkerCount sets a new worker count.
//
//	Parameters: workerCount int.
func (c *FixedRateTimer) SetWorkerCount(workerCount int) {
	if workerCount < 1 || c.IsStarted() {
		return
	}
	c.workerCount = workerCount
}

// Callback gets the callback function that is called when this timer is triggered.
//
//	Returns: function the callback function or null if it is not set.
func (c *FixedRateTimer) Callback() func(ctx context.Context) {
	return c.callback
}

// SetCallback sets the callback function that is called when this timer is triggered.
//
//	Parameters: value func() the callback function to be called.
func (c *FixedRateTimer) SetCallback(value func(ctx context.Context)) {
	if c.IsStarted() {
		return
	}
	c.callback = value
	c.task = nil
}

// Delay gets initial delay before the timer is triggered for the first time.
//
//	Returns: number the delay in milliseconds.
func (c *FixedRateTimer) Delay() int {
	return c.delay
}

// SetDelay sets initial delay before the timer is triggered for the first time.
//
//	Parameters: value int a delay in milliseconds.
func (c *FixedRateTimer) SetDelay(value int) {
	if c.IsStarted() {
		return
	}
	c.delay = value
}

// Interval gets periodic timer triggering interval.
//
//	Returns: number the interval in milliseconds
func (c *FixedRateTimer) Interval() int {
	return c.interval
}

// SetInterval sets periodic timer triggering interval.
//
//	Parameters: value int an interval in milliseconds.
func (c *FixedRateTimer) SetInterval(value int) {
	if c.IsStarted() {
		return
	}
	c.interval = value
}

// IsStarted checks if the timer is started.
//
//	Returns: bool true if the timer is started and false if it is stopped.
func (c *FixedRateTimer) IsStarted() bool {
	return c.ticker != nil
}

// Start starts the timer. Initially the timer is triggered after delay.
// After that it is triggered after interval until it is stopped.
func (c *FixedRateTimer) Start(ctx context.Context) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// Stop previously set timer
	c.stop(ctx)

	// Exit if interval is not defined
	if c.interval <= 0 {
		return
	}

	// Introducing delay
	delay := c.delay - c.interval
	ticker := time.NewTicker(time.Millisecond * time.Duration(c.interval))
	c.ticker = ticker

	callback := c.callback
	exit := c.exit

	if delay > 0 {
		time.Sleep(time.Millisecond * time.Duration(delay))
	}

	for i := 0; i < c.workerCount; i++ {
		go func() {
			defer cctx.DefaultErrorHandlerWithShutdown(ctx)

			for {
				select {
				case <-ticker.C:
					if callback != nil {
						callback(ctx)
					}
					break
				case _, ok := <-exit:
					if !ok {
						return
					}
				}
			}
		}()
	}
}

// Stop the timer.
func (c *FixedRateTimer) Stop(ctx context.Context) {
	c.mtx.Lock()
	c.stop(ctx)
	c.mtx.Unlock()
}

// stop is a private function to implement thread save
func (c *FixedRateTimer) stop(ctx context.Context) {
	ticker := c.ticker
	if ticker != nil {
		ticker.Stop()
		c.ticker = nil
		close(c.exit)
	}
}

// Close closes the timer.
// This is required by ICloseable interface, but besides that it is identical to stop().
//
//	Parameters: ctx context.Context a context to trace execution through call chain.
//	Returns: error
func (c *FixedRateTimer) Close(ctx context.Context) error {
	c.mtx.Lock()
	c.stop(ctx)
	c.mtx.Unlock()
	return nil
}
