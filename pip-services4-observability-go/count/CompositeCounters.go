package count

import (
	"context"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// CompositeCounters aggregates all counters from component references under a single component.
// It allows capturing metrics and conveniently send them to multiple destinations.
//
//		References:
//			- *:counters:*:*:1.0 (optional) ICounters components to pass collected measurements
//		see ICounters
//
//		Example:
//	 	type MyComponent {
//				_counters CompositeCounters = new CompositeCounters();
//	 	}
//			func (mc *MyConponent) SetReferences(ctx context.Context, references refer.IReferences) {
//				mc._counters.SetReferences(ctx, references);
//			}
//
//			func (mc * MyConponent) myMethod() {
//				mc._counters.Increment(context.Background(), "mycomponent.mymethod.calls");
//				timing := mc._counters.BeginTiming(context.Background(), "mycomponent.mymethod.exec_time");
//				defer timing.EndTiming(context.Background());
//	 		// do something
//			}
//			var mc MyComponent{};
//			mc._counters = NewCompositeCounters();
type CompositeCounters struct {
	counters []ICounters
}

// NewCompositeCounters creates a new instance of the counters.
//
//	Returns: *CompositeCounters
func NewCompositeCounters() *CompositeCounters {
	c := &CompositeCounters{
		counters: []ICounters{},
	}
	return c
}

// NewCompositeCountersFromReferences creates a new instance of the counters.
//
//	Parameters:
//		- ctx context.Context
//		- references is a refer.IReferences to locate the component dependencies.
//	Returns: *CompositeCounters
func NewCompositeCountersFromReferences(ctx context.Context, references refer.IReferences) *CompositeCounters {
	c := NewCompositeCounters()
	c.SetReferences(ctx, references)
	return c
}

// SetReferences references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references refer.IReferences references to locate the component dependencies.
func (c *CompositeCounters) SetReferences(ctx context.Context, references refer.IReferences) {
	if c.counters == nil {
		c.counters = []ICounters{}
	}

	counters := references.GetOptional(
		refer.NewDescriptor("*", "counters", "*", "*", "*"),
	)
	for _, l := range counters {
		if l == c {
			continue
		}

		counter, ok := l.(ICounters)
		if ok {
			c.counters = append(c.counters, counter)
		}
	}
}

// BeginTiming begins measurement of execution time interval.
// It returns Timing object which has to be called at
// Timing.endTiming to end the measurement and update the counter.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Interval type.
//	Returns: *Timing a Timing callback object to end timing.
func (c *CompositeCounters) BeginTiming(ctx context.Context, name string) *CounterTiming {
	return NewCounterTiming(name, c)
}

// EndTiming ends measurement of execution elapsed time and updates specified counter.
//
//	see Timing.EndTiming
//	Parameters:
//		- ctx context.Context
//		- name string a counter name
//		- elapsed float64 execution elapsed time in milliseconds to update the counter.
func (c *CompositeCounters) EndTiming(ctx context.Context, name string, elapsed float64) {
	for _, counter := range c.counters {
		if counter != nil {
			if callback, ok := counter.(ITimingCallback); ok {
				callback.EndTiming(ctx, name, elapsed)
			}
		}
	}
}

// Stats calculates min/average/max statistics based on the current and previous values.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Statistics type
//		- value float64 a value to update statistics
func (c *CompositeCounters) Stats(ctx context.Context, name string, value float64) {
	for _, counter := range c.counters {
		if counter != nil {
			counter.Stats(ctx, name, value)
		}
	}
}

// Last records the last calculated measurement value.
// Usually this method is used by metrics calculated externally.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Last type.
//		- value float64 a last value to record.
func (c *CompositeCounters) Last(ctx context.Context, name string, value float64) {
	for _, counter := range c.counters {
		if counter != nil {
			counter.Last(ctx, name, value)
		}
	}
}

// TimestampNow records the current time as a timestamp.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Timestamp type.
func (c *CompositeCounters) TimestampNow(ctx context.Context, name string) {
	c.Timestamp(ctx, name, time.Now())
}

// Timestamp records the given timestamp.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Timestamp type.
//		- value time.Time a timestamp to record.
func (c *CompositeCounters) Timestamp(ctx context.Context, name string, value time.Time) {
	for _, counter := range c.counters {
		if counter != nil {
			counter.Timestamp(ctx, name, value)
		}
	}
}

// IncrementOne increments counter by 1.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Increment type.
func (c *CompositeCounters) IncrementOne(ctx context.Context, name string) {
	c.Increment(ctx, name, 1)
}

// Increment increments counter by given value.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Increment type.
//		- value number a value to add to the counter.
func (c *CompositeCounters) Increment(ctx context.Context, name string, value int64) {
	for _, counter := range c.counters {
		if counter != nil {
			counter.Increment(ctx, name, value)
		}
	}
}
