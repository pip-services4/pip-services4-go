package trace

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// CompositeTracer aggregates all tracers from component references under a single component.
// It allows to record traces and conveniently send them to multiple destinations.
//	References:
//		- *:tracer:*:*:1.0 (optional) ITracer components to pass operation traces
//	See ITracer

// Example:
//
//	type MyComponent struct {
//		tracer CompositeTracer
//	}
//	func NewMyComponent() *MyComponent{
//		return &MyComponent{
//			tracer: NewCompositeTracer(nil);
//		}
//	}
//	func (c* MyComponent) SetReferences(ctx context.Context, references IReferences) {
//		c.tracer.SetReferences(references)
//		...
//	}
//	public MyMethod(ctx context.Context) {
//		timing := c.tracer.BeginTrace(ctx, "mycomponent", "mymethod");
//		...
//		timing.EndTrace(ctx);
//		if err != nil {
//			timing.EndFailure(ctx, err);
//		}
//	}
type CompositeTracer struct {
	Tracers []ITracer
}

// NewCompositeTracer creates a new instance of the tracer.
//
//	Parameters:
//		- references to locate the component dependencies.
func NewCompositeTracer() *CompositeTracer {
	return &CompositeTracer{}
}

// NewCompositeTracerFromReferences creates a new instance of the tracer.
//
//	Parameters:
//		- ctx context.Context
//		- refer.IReferences references to locate the component dependencies.
//	Returns: CompositeLogger
func NewCompositeTracerFromReferences(ctx context.Context, references cref.IReferences) *CompositeTracer {
	c := NewCompositeTracer()
	c.SetReferences(ctx, references)
	return c
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references to locate the component dependencies.
func (c *CompositeTracer) SetReferences(ctx context.Context, references cref.IReferences) {

	if c.Tracers == nil {
		c.Tracers = []ITracer{}
	}

	tracers := references.GetOptional(cref.NewDescriptor("*", "tracer", "*", "*", "*"))
	for _, l := range tracers {
		if l == c {
			continue
		}

		if tracer, ok := l.(ITracer); ok {
			c.Tracers = append(c.Tracers, tracer)
		}
	}

}

// Trace records an operation trace with its name and duration
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- duration      execution duration in milliseconds.
func (c *CompositeTracer) Trace(ctx context.Context, component string, operation string, duration int64) {
	for _, tracer := range c.Tracers {
		tracer.Trace(ctx, component, operation, duration)
	}
}

// Failure records an operation failure with its name, duration and error
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- error         an error object associated with this trace.
//		- duration      execution duration in milliseconds.
func (c *CompositeTracer) Failure(ctx context.Context, component string, operation string, err error, duration int64) {
	for _, tracer := range c.Tracers {
		tracer.Failure(ctx, component, operation, err, duration)
	}
}

// BeginTrace begins recording an operation trace
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//	Returns: a trace timing object.
func (c *CompositeTracer) BeginTrace(ctx context.Context, component string, operation string) *TraceTiming {
	return NewTraceTiming(ctx, component, operation, c)
}
