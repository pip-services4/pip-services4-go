package test_count

import (
	"context"
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	"github.com/stretchr/testify/assert"
)

type CountersFixture struct {
	counters *count.CachedCounters
}

func NewCountersFixture(counters *count.CachedCounters) *CountersFixture {
	return &CountersFixture{
		counters: counters,
	}
}

func (c *CountersFixture) TestSimpleCounters(t *testing.T) {
	c.counters.Last(context.Background(), "Test.LastValue", 123)
	c.counters.Last(context.Background(), "Test.LastValue", 123456)

	counter, ok := c.counters.Get(context.Background(), "Test.LastValue", count.LastValue)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.Equal(t, float64(123456), counter.Last())

	c.counters.IncrementOne(context.Background(), "Test.Increment")
	c.counters.Increment(context.Background(), "Test.Increment", 3)

	counter, ok = c.counters.Get(context.Background(), "Test.Increment", count.Increment)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.Equal(t, int64(4), counter.Count())

	c.counters.TimestampNow(context.Background(), "Test.Timestamp")
	c.counters.TimestampNow(context.Background(), "Test.Timestamp")

	counter, ok = c.counters.Get(context.Background(), "Test.Timestamp", count.Timestamp)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.NotNil(t, counter.Time())

	c.counters.Stats(context.Background(), "Test.Statistics", 1)
	c.counters.Stats(context.Background(), "Test.Statistics", 2)
	c.counters.Stats(context.Background(), "Test.Statistics", 3)

	counter, ok = c.counters.Get(context.Background(), "Test.Statistics", count.Statistics)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.Equal(t, float64(2), counter.Average())

	_ = c.counters.Dump(context.Background())
}

func (c *CountersFixture) TestMeasureElapsedTime(t *testing.T) {
	timing := c.counters.BeginTiming(context.Background(), "Test.Elapsed")

	time.Sleep(100 * time.Millisecond)

	timing.EndTiming(context.Background())

	counter, ok := c.counters.Get(context.Background(), "Test.Elapsed", count.Interval)
	assert.True(t, ok)
	assert.NotNil(t, counter)
	assert.True(t, counter.Last() > 50)
	assert.True(t, counter.Last() < 5000)

	c.counters.Dump(context.Background())
}
