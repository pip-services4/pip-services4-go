package trace

import (
	"context"

	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
)

type InstrumentTiming struct {
	context       context.Context
	name          string
	verb          string
	logger        clog.ILogger
	counters      ccount.ICounters
	counterTiming *ccount.CounterTiming
	traceTiming   *ctrace.TraceTiming
}

func NewInstrumentTiming(context context.Context, name string,
	verb string, logger clog.ILogger, counters ccount.ICounters,
	counterTiming *ccount.CounterTiming, traceTiming *ctrace.TraceTiming) *InstrumentTiming {

	if len(verb) == 0 {
		verb = "call"
	}
	return &InstrumentTiming{
		context:       context,
		name:          name,
		verb:          verb,
		logger:        logger,
		counters:      counters,
		counterTiming: counterTiming,
		traceTiming:   traceTiming,
	}
}

func (c *InstrumentTiming) clear() {
	// Clear references to avoid double processing
	c.counters = nil
	c.logger = nil
	c.counterTiming = nil
	c.traceTiming = nil
}

func (c *InstrumentTiming) EndTiming(ctx context.Context, err error) {
	if err != nil {
		c.EndFailure(ctx, err)
	} else {
		c.EndSuccess(ctx)
	}
}

func (c *InstrumentTiming) EndSuccess(ctx context.Context) {
	if c.counterTiming != nil {
		c.counterTiming.EndTiming(ctx)
	}
	if c.traceTiming != nil {
		c.traceTiming.EndTrace()
	}

	c.clear()
}

func (c *InstrumentTiming) EndFailure(ctx context.Context, err error) {
	if c.counterTiming != nil {
		c.counterTiming.EndTiming(ctx)
	}

	if err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, err, "Failed to call %s method", c.name)
		}
		if c.counters != nil {
			c.counters.IncrementOne(ctx, c.name+"."+c.verb+"_errors")
		}
		if c.traceTiming != nil {
			c.traceTiming.EndFailure(err)
		}
	} else {
		if c.traceTiming != nil {
			c.traceTiming.EndTrace()
		}
	}

	c.clear()
}
