package count

import (
	"math"
	"sync"
	"time"
)

// AtomicCounter data object to store measurement for a performance counter.
// This object is used by CachedCounters to store counters.
type AtomicCounter struct {
	_mtx sync.RWMutex

	_name    string
	_type    CounterType
	_time    time.Time
	_last    float64
	_min     float64
	_max     float64
	_average float64
	_count   int64
}

// NewAtomicCounter creates an instance of the data obejct
//	Parameters:
//		- name string a counter name.
//		- type CounterType a counter type.
//	Returns: *Counter
func NewAtomicCounter(name string, typ CounterType) *AtomicCounter {
	return &AtomicCounter{
		_mtx:  sync.RWMutex{},
		_name: name,
		_type: typ,

		_min: math.MaxFloat64,
		_max: -math.MaxFloat64,
	}
}

// SetLast is a setter for the _last
//	Parameters: value float64
func (c *AtomicCounter) SetLast(value float64) {
	c._mtx.Lock()
	defer c._mtx.Unlock()
	c._last = value
}

// SetTime is a setter for the _time
//	Parameters: value time.Time
func (c *AtomicCounter) SetTime(value time.Time) {
	c._mtx.Lock()
	defer c._mtx.Unlock()
	c._time = value
}

// CalculateStats set up _last and calculates stats
//	Parameters: value float64
func (c *AtomicCounter) CalculateStats(value float64) {
	c._mtx.Lock()
	defer c._mtx.Unlock()

	c._last = value
	c._count++
	c._max = math.Max(c._max, value)
	c._min = math.Min(c._min, value)
	c._average = ((c._average * float64(c._count-1)) + value) / float64(c._count)
}

// Inc increments _count for the provided value
//	Parameters: value float64
func (c *AtomicCounter) Inc(value int64) {
	c._mtx.Lock()
	defer c._mtx.Unlock()

	c._count += value
}

// GetCounter converts AtomicCounter to Counter struct
//	Returns: Counter
func (c *AtomicCounter) GetCounter() Counter {
	c._mtx.RLock()
	defer c._mtx.RUnlock()

	return Counter{
		Name:    c._name,
		Type:    c._type,
		Last:    c._last,
		Count:   c._count,
		Min:     c._min,
		Max:     c._max,
		Average: c._average,
		Time:    c._time,
	}
}

// Name gets counter _name
//	Returns: string
func (c *AtomicCounter) Name() string {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._name
}

// Type gets counter _type
//	Returns: int
func (c *AtomicCounter) Type() CounterType {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._type
}

// Time gets counter _time
//	Returns: time.Time
func (c *AtomicCounter) Time() time.Time {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._time
}

// Last gets counter _last
//	Returns: float64
func (c *AtomicCounter) Last() float64 {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._last
}

// Count gets counter _count
//	Returns: int64
func (c *AtomicCounter) Count() int64 {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._count
}

// Min gets counter _min
//	Returns: float64
func (c *AtomicCounter) Min() float64 {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._min
}

// Max gets counter _max
//	Returns: float64
func (c *AtomicCounter) Max() float64 {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._max
}

// Average gets counter _average
//	Returns: float64
func (c *AtomicCounter) Average() float64 {
	c._mtx.RLock()
	defer c._mtx.RUnlock()
	return c._average
}
