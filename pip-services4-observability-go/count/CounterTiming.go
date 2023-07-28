package count

import (
	"context"
	"time"
)

// CounterTiming callback object returned by ICounters.beginTiming to end timing of
// execution block and update the associated counter.
//
//	Example:
//		timing := counters.BeginTiming(contex.Background(), "mymethod.exec_time")
//		defer  timing.EndTiming()
type CounterTiming struct {
	start    time.Time
	callback ICounterTimingCallback
	counter  string
}

// NewEmptyCounterTiming creates a new instance of the timing callback object.
//
//	Returns: *CounterTiming
func NewEmptyCounterTiming() *CounterTiming {
	return &CounterTiming{
		start: time.Now(),
	}
}

// NewCounterTiming creates a new instance of the timing callback object.
//
//	Parameters:
//		- counter string an associated counter name
//		- callback ITimingCallback a callback that shall be called when EndTiming is called.
//	Returns: *Timing
func NewCounterTiming(counter string, callback ICounterTimingCallback) *CounterTiming {
	return &CounterTiming{
		start:    time.Now(),
		callback: callback,
		counter:  counter,
	}
}

// EndTiming ends timing of an execution block, calculates
// elapsed time and updates the associated counter.
func (c *CounterTiming) EndTiming(ctx context.Context) {
	if c.callback == nil {
		return
	}

	elapsed := time.Since(c.start).Seconds() * 1000
	c.callback.EndTiming(ctx, c.counter, elapsed)
}
