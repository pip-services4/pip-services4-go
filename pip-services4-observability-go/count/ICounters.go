package count

import (
	"context"
	"time"
)

// ICounters interface for performance counters that measure execution metrics.
// The performance counters measure how code is performing: how fast or slow,
// how many transactions performed, how many objects are stored, what was the latest transaction time and so on.
// They are critical to monitor and improve performance, scalability and reliability of code in production.
type ICounters interface {
	// BeginTiming begins measurement of execution time interval.
	// It returns Timing object which has to be called at
	// Timing.endTiming to end the measurement and update the counter.
	BeginTiming(ctx context.Context, name string) *CounterTiming

	// Stats calculates min/average/max statistics based on the current and previous values.
	Stats(ctx context.Context, name string, value float64)

	// Last records the last calculated measurement value.
	// Usually this method is used by metrics calculated externally.
	Last(ctx context.Context, name string, value float64)

	// TimestampNow records the given timestamp.
	TimestampNow(ctx context.Context, name string)

	// Timestamp records the current time as a timestamp.
	Timestamp(ctx context.Context, name string, value time.Time)

	// IncrementOne increments counter by 1.
	IncrementOne(ctx context.Context, name string)

	// Increment increments counter by given value.
	Increment(ctx context.Context, name string, value int64)
}
