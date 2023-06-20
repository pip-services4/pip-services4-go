package count

import (
	"context"
	"time"
)

// NullCounters dummy implementation of performance counters that doesn't do anything.
// It can be used in testing or in situations when counters is required but shall be disabled.
type NullCounters struct{}

// NewNullCounters creates a new instance of the counter.
//	Returns: *NullCounters
func NewNullCounters() *NullCounters {
	return &NullCounters{}
}

// BeginTiming begins measurement of execution time interval.
// It returns Timing object which has to be called at Timing.EndTiming
// to end the measurement and update the counter.
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Interval type.
//	Returns: *Timing a Timing callback object to end timing.
func (c *NullCounters) BeginTiming(ctx context.Context, name string) *CounterTiming {
	return NewEmptyCounterTiming()
}

// Stats calculates min/average/max statistics based on the current and previous values.
// Parameters:
//		- ctx context.Context
//		- name string a counter name of Statistics type
//		- value float64 a value to update statistics
func (c *NullCounters) Stats(ctx context.Context, name string, value float64) {}

// Last records the last calculated measurement value.
// Usually this method is used by metrics calculated externally.
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Last type.
//		- value float64 a last value to record.
func (c *NullCounters) Last(ctx context.Context, name string, value float64) {}

// TimestampNow records the current time as a timestamp.
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Timestamp type.
func (c *NullCounters) TimestampNow(ctx context.Context, name string) {}

// Timestamp records the given timestamp.
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Timestamp type.
//		- value time.Time a timestamp to record.
func (c *NullCounters) Timestamp(ctx context.Context, name string, value time.Time) {}

// IncrementOne increments counter by 1.
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Increment type.
func (c *NullCounters) IncrementOne(ctx context.Context, name string) {}

// Increment increments counter by given value.
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Increment type.
//		- value int64 a value to add to the counter.
func (c *NullCounters) Increment(ctx context.Context, name string, value int64) {}
