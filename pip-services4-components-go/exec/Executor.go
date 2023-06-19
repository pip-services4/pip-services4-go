package exec

import "context"

// Executor Helper class that executes components.
var Executor = &_TExecutor{}

type _TExecutor struct{}

// ExecuteOne executes specific component.
// To be executed components must implement IExecutable interface.
// If they don't the call to this method has no effect.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component any the component that is to be executed.
//		- args: *Parameters execution arguments.
//	Returns: []any, error execution result or error
func (c *_TExecutor) ExecuteOne(ctx context.Context, component any, args *Parameters) (any, error) {
	if v, ok := component.(IExecutable); ok {
		return v.Execute(ctx, args)
	}
	return nil, nil
}

// Executes multiple components.

// Execute to be executed components must implement IExecutable interface.
// If they don't the call to this method has no effect.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- components []any a list of components that are to be executed.
//		- args *Parameters execution arguments.
//	Returns: []any, error execution result or error
func (c *_TExecutor) Execute(ctx context.Context, components []any, args *Parameters) ([]any, error) {
	results := make([]any, 0, 5)

	for _, component := range components {
		result, err := c.ExecuteOne(ctx, component, args)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}
