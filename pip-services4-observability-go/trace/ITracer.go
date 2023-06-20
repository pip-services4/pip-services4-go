package trace

import "context"

// ITracer interface for tracer components that capture operation traces.
type ITracer interface {

	// Trace records an operation trace with its name and duration
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- component         a name of called component
	//		- operation         a name of the executed operation.
	//		- duration          execution duration in milliseconds.
	Trace(ctx context.Context, component string, operation string, duration int64)

	// Failure records an operation failure with its name, duration and error
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- component         a name of called component
	//		- operation         a name of the executed operation.
	//		- error             an error object associated with this trace.
	//		- duration          execution duration in milliseconds.
	Failure(ctx context.Context, component string, operation string, err error, duration int64)

	//BeginTrace recording an operation trace
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- component         a name of called component
	//		- operation         a name of the executed operation.
	//	Returns: a trace timing object.
	BeginTrace(ctx context.Context, component string, operation string) *TraceTiming
}
