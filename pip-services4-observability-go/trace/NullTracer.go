package trace

import "context"

// NullTracer dummy implementation of tracer that doesn't do anything.
// It can be used in testing or in situations when tracing is required
// but shall be disabled.
//	See ITracer
type NullTracer struct {
}

// NewNullTracer creates a new instance of the tracer.
func NewNullTracer() *NullTracer {
	return &NullTracer{}
}

// Trace records an operation trace with its name and duration
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- duration      execution duration in milliseconds.
func (c *NullTracer) Trace(ctx context.Context, component string, operation string, duration int64) {
	// Do nothing...
}

// Failure records an operation failure with its name, duration and error
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component         a name of called component
//		- operation         a name of the executed operation.
//		- error             an error object associated with this trace.
//		- duration          execution duration in milliseconds.
func (c *NullTracer) Failure(ctx context.Context, component string, operation string, err error, duration int64) {
	// Do nothing...
}

// BeginTrace begins recording an operation trace
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component         a name of called component
//		- operation         a name of the executed operation.
//	Returns: a trace timing object.
func (c *NullTracer) BeginTrace(ctx context.Context, component string, operation string) *TraceTiming {
	return NewTraceTiming(ctx, component, operation, c)
}
