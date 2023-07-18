package test_fixture

import (
	"context"
	"testing"
	"time"

	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	"github.com/stretchr/testify/assert"
)

type CountersFixture struct {
	counters *ccount.CachedCounters
}

func NewCountersFixture(counters *ccount.CachedCounters) *CountersFixture {
	cf := CountersFixture{}
	cf.counters = counters
	return &cf
}

func (c *CountersFixture) TestSimpleCounters(t *testing.T) {
	c.counters.Last(context.Background(), "Test.LastValue", 123)
	c.counters.Last(context.Background(), "Test.LastValue", 123456)

	var counter, _ = c.counters.Get(context.Background(), "Test.LastValue", ccount.LastValue)
	assert.NotNil(t, counter)
	assert.NotNil(t, counter.Last())
	assert.Equal(t, counter.Last(), float64(123456), 3)

	c.counters.IncrementOne(context.Background(), "Test.Increment")
	c.counters.Increment(context.Background(), "Test.Increment", 3)

	counter, _ = c.counters.Get(context.Background(), "Test.Increment", ccount.Increment)
	assert.NotNil(t, counter)
	assert.Equal(t, counter.Count(), int64(4))

	c.counters.TimestampNow(context.Background(), "Test.Timestamp")
	c.counters.TimestampNow(context.Background(), "Test.Timestamp")

	counter, _ = c.counters.Get(context.Background(), "Test.Timestamp", ccount.Timestamp)
	assert.NotNil(t, counter)
	assert.NotNil(t, counter.Time())

	c.counters.Stats(context.Background(), "Test.Statistics", 1)
	c.counters.Stats(context.Background(), "Test.Statistics", 2)
	c.counters.Stats(context.Background(), "Test.Statistics", 3)

	counter, _ = c.counters.Get(context.Background(), "Test.Statistics", ccount.Statistics)
	assert.NotNil(t, counter)
	assert.Equal(t, counter.Average(), float64(2), 3)

	c.counters.Dump(context.Background())

	select {
	case <-time.After(1000 * time.Millisecond):
		{
		}
	}
}

func (c *CountersFixture) TestMeasureElapsedTime(t *testing.T) {
	timer := c.counters.BeginTiming(context.Background(), "Test.Elapsed")

	time.AfterFunc(100*time.Millisecond, func() {
		timer.EndTiming(context.Background())

		counter, _ := c.counters.Get(context.Background(), "Test.Elapsed", ccount.Interval)
		assert.True(t, counter.Last() > 50)
		assert.True(t, counter.Last() < 5000)

		c.counters.Dump(context.Background())
		select {
		case <-time.After(1000 * time.Millisecond):
			{
			}
		}
	})
}
