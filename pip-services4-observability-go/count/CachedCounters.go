package count

import (
	"context"
	"sync"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// CachedCounters abstract implementation of performance counters that measures and
// stores counters in memory. Child classes implement saving of the counters
// into various destinations.
//
//	Configuration parameters:
//		- options:
//			- interval: interval in milliseconds to save current counters measurements (default: 5 mins)
//			- reset_timeout: timeout in milliseconds to reset the counters. 0 disables the reset (default: 0)
type CachedCounters struct {
	cache         map[string]*AtomicCounter
	updated       bool
	lastDumpTime  time.Time
	lastResetTime time.Time
	mux           sync.RWMutex
	interval      int64
	resetTimeout  int64
	Overrides     ICachedCountersOverrides
}

type ICachedCountersOverrides interface {
	Save(ctx context.Context, counters []Counter) error
}

const (
	DefaultInterval             int64 = 300000
	DefaultResetTimeout         int64 = 300000
	ConfigParameterInterval           = "interval"
	ConfigParameterResetTimeout       = "reset_timeout"
)

// InheritCacheCounters inherit cache counters from saver
//
//	Parameters:
//		- save ICountersSaver
//	Returns: *CachedCounters
func InheritCacheCounters(overrides ICachedCountersOverrides) *CachedCounters {
	return &CachedCounters{
		cache:         make(map[string]*AtomicCounter),
		updated:       false,
		lastDumpTime:  time.Now(),
		lastResetTime: time.Now(),
		interval:      DefaultInterval,
		resetTimeout:  DefaultResetTimeout,
		Overrides:     overrides,
	}
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *CachedCounters) Configure(ctx context.Context, config *config.ConfigParams) {
	c.interval = config.GetAsLongWithDefault(ConfigParameterInterval, c.interval)
	c.resetTimeout = config.GetAsLongWithDefault(ConfigParameterResetTimeout, c.resetTimeout)
}

// Clear clears (resets) a counter specified by its name.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name to clear.
func (c *CachedCounters) Clear(ctx context.Context, name string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.cache, name)
}

// ClearAll clears (resets) all counters.
//
//	Parameters:
//		- ctx context.Context
func (c *CachedCounters) ClearAll(ctx context.Context) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cache = make(map[string]*AtomicCounter)
}

func (c *CachedCounters) isUpdated() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.updated
}

// Dump (saves) the current values of counters.
//
//	Parameters:
//		- ctx context.Context
func (c *CachedCounters) Dump(ctx context.Context) error {
	if !c.isUpdated() {
		return nil
	}

	counters := c.GetAllCountersStats()
	err := c.Overrides.Save(ctx, counters)
	if err != nil {
		return err
	}

	c.mux.Lock()
	c.updated = false
	c.lastDumpTime = time.Now()
	c.mux.Unlock()

	return nil
}

func (c *CachedCounters) update(ctx context.Context) error {
	c.mux.Lock()
	c.updated = true
	newDumpTime := c.lastDumpTime.Add(time.Duration(c.interval) * time.Millisecond)
	c.mux.Unlock()
	if time.Now().After(newDumpTime) {
		return c.Dump(ctx)
	}
	return nil
}

func (c *CachedCounters) resetIfNeeded(ctx context.Context) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.resetTimeout == 0 {
		return
	}

	newResetTime := c.lastResetTime.Add(time.Duration(c.resetTimeout) * time.Millisecond)
	if time.Now().After(newResetTime) {
		c.cache = make(map[string]*AtomicCounter)
		c.updated = false
		c.lastDumpTime = time.Now()
	}
}

// GetAll gets all captured counters.
//
//	Returns: []*AtomicCounter
func (c *CachedCounters) GetAll() []*AtomicCounter {
	c.mux.RLock()
	defer c.mux.RUnlock()

	result := make([]*AtomicCounter, 0, len(c.cache))
	for _, v := range c.cache {
		result = append(result, v)
	}

	return result
}

// GetAllCountersStats gets all captured counters stats.
//
//	Returns: []Counter
func (c *CachedCounters) GetAllCountersStats() []Counter {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.getAllCountersStats()
}

func (c *CachedCounters) getAllCountersStats() []Counter {

	result := make([]Counter, 0, len(c.cache))
	for _, v := range c.cache {
		result = append(result, v.GetCounter())
	}

	return result
}

// Get a counter specified by its name. It counter does not exist or its type doesn't match the
// specified type it creates a new one.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name to retrieve.
//		- typ int a counter type.
//	Returns: *Counter an existing or newly created counter of the specified type.
func (c *CachedCounters) Get(ctx context.Context, name string, typ CounterType) (*AtomicCounter, bool) {
	if name == "" {
		return nil, false
	}

	c.resetIfNeeded(ctx)

	c.mux.Lock()
	defer c.mux.Unlock()

	counter, ok := c.cache[name]
	if !ok || counter.Type() != typ {
		counter = NewAtomicCounter(name, typ)
		c.cache[name] = counter
	}

	return counter, true
}

// BeginTiming begins measurement of execution time interval.
// It returns Timing object which has to be called at
// Timing.EndTiming to end the measurement and update the counter.
//
//	Parameters
//		- ctx context.Context
//		- name string a counter name of Interval type.
//	Returns: *Timing a Timing callback object to end timing.
func (c *CachedCounters) BeginTiming(ctx context.Context, name string) *CounterTiming {
	return NewCounterTiming(name, c)
}

// EndTiming ends measurement of execution elapsed time and updates specified counter.
//
//	see Timing.EndTiming
//	Parameters:
//		- ctx context.Context
//		- name string a counter name
//		- elapsed float64 execution elapsed time in milliseconds to update the counter.
func (c *CachedCounters) EndTiming(ctx context.Context, name string, elapsed float64) {
	if counter, ok := c.Get(ctx, name, Interval); ok {
		counter.CalculateStats(elapsed)
		_ = c.update(ctx)
	}
}

// Stats calculates min/average/max statistics based on the current and previous values.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Statistics type
//		- value float32 a value to update statistics
func (c *CachedCounters) Stats(ctx context.Context, name string, value float64) {
	if counter, ok := c.Get(ctx, name, Statistics); ok {
		counter.CalculateStats(value)
		_ = c.update(ctx)
	}
}

// Last records the last calculated measurement value.
// Usually this method is used by metrics calculated externally.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Last type.
//		- value number a last value to record.
func (c *CachedCounters) Last(ctx context.Context, name string, value float64) {
	if counter, ok := c.Get(ctx, name, LastValue); ok {
		counter.SetLast(value)
		_ = c.update(ctx)
	}
}

// TimestampNow records the current time as a timestamp.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Timestamp type.
func (c *CachedCounters) TimestampNow(ctx context.Context, name string) {
	c.Timestamp(ctx, name, time.Now())
}

// Timestamp records the given timestamp.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Timestamp type.
//		- value time.Time a timestamp to record.
func (c *CachedCounters) Timestamp(ctx context.Context, name string, value time.Time) {
	if counter, ok := c.Get(ctx, name, Timestamp); ok {
		counter.SetTime(value)
		_ = c.update(ctx)
	}
}

// IncrementOne increments counter by 1.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Increment type.
func (c *CachedCounters) IncrementOne(ctx context.Context, name string) {
	c.Increment(ctx, name, 1)
}

// Increment increments counter by given value.
//
//	Parameters:
//		- ctx context.Context
//		- name string a counter name of Increment type.
//		- value int a value to add to the counter.
func (c *CachedCounters) Increment(ctx context.Context, name string, value int64) {
	if counter, ok := c.Get(ctx, name, Increment); ok {
		counter.Inc(value)
		_ = c.update(ctx)
	}
}
