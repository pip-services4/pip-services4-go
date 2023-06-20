package trace

import (
	"context"
	"time"
)

// TraceTiming timing object returned by {ITracer.BeginTrace} to end timing
// of execution block and record the associated trace.
//
//	Example:
//		timing := tracer.BeginTrace(context.Background(), "123", "my_component","mymethod.exec_time");
//		...
//		timing.EndTrace(context.Background());
//		if err != nil {
//			timing.EndFailure(context.Background(), err);
//		}
type TraceTiming struct {
	context   context.Context
	start     int64
	tracer    ITracer
	component string
	operation string
}

// NewTraceTiming creates a new instance of the timing callback object.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component 	an associated component name
//		- operation 	an associated operation name
//		- callback 		a callback that shall be called when endTiming is called.
func NewTraceTiming(ctx context.Context, component string, operation string, tracer ITracer) *TraceTiming {
	return &TraceTiming{
		context:   ctx,
		component: component,
		operation: operation,
		tracer:    tracer,
		start:     time.Now().UTC().UnixNano(),
	}
}

// EndTrace ends timing of an execution block, calculates elapsed time
// and records the associated trace.
//
//	Parameters:
//		- ctx context.Context
func (c *TraceTiming) EndTrace() {
	if c.tracer != nil {
		elapsed := time.Now().UTC().UnixNano() - c.start
		c.tracer.Trace(c.context, c.component, c.operation, elapsed/int64(time.Millisecond))
	}
}

// EndFailure ends timing of a failed block, calculates elapsed time
// and records the associated trace.
//
//	Parameters:
//		- ctx context.Context
func (c *TraceTiming) EndFailure(err error) {
	if c.tracer != nil {
		elapsed := time.Now().UTC().UnixNano() - c.start
		c.tracer.Failure(c.context, c.component, c.operation, err, elapsed/int64(time.Millisecond))
	}
}
