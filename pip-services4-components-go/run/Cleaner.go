package run

import "context"

// Cleaner helper class that cleans stored object state.
var Cleaner = &_TCleaner{}

type _TCleaner struct{}

// ClearOne clears state of specific component.
// To be cleaned state components must implement ICleanable interface. If they don't the call to this method has no effect.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component any the component that is to be cleaned.
//	Returns: error
func (c *_TCleaner) ClearOne(ctx context.Context, component any) error {
	if v, ok := component.(ICleanable); ok {
		return v.Clear(ctx)
	}
	return nil
}

// Clear clears state of multiple components.
// To be cleaned state components must implement ICleanable interface.
// If they don't the call to this method has no effect.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- components []any the list of components that are to be cleaned.
//	Returns: error
func (c *_TCleaner) Clear(ctx context.Context, components []any) error {
	for _, component := range components {
		err := c.ClearOne(ctx, component)
		if err != nil {
			return err
		}
	}
	return nil
}
