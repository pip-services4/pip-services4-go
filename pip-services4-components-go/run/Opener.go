package run

import "context"

// Opener Helper class that opens components.
var Opener = &_TOpener{}

type _TOpener struct{}

// IsOpenOne checks if specified component is opened.
// To be checked components must implement IOpenable interface.
// If they don't the call to this method returns true.
//	see IOpenable
//	Parameters: component any the component that is to be checked.
//	Returns: bool true if component is opened and false otherwise.
func (c *_TOpener) IsOpenOne(component any) bool {
	if v, ok := component.(IOpenable); ok {
		return v.IsOpen()
	}
	return true
}

// IsOpen checks if all components are opened.
// To be checked components must implement IOpenable interface.
// If they don't the call to this method returns true.
//	see IsOpenOne
//	see IOpenable
//	Parameters: components []any a list of components that are to be checked.
//	Returns: bool true if all components are opened and false if at least one component is closed.
func (c *_TOpener) IsOpen(components []any) bool {
	result := true

	for _, component := range components {
		if result = result && c.IsOpenOne(component); !result {
			return result
		}
	}
	return result
}

// OpenOne opens specific component.
// To be opened components must implement IOpenable interface.
// If they don't the call to this method has no effect.
//	see IOpenable
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component any the component that is to be opened.
//	Returns: error
func (c *_TOpener) OpenOne(ctx context.Context, component any) error {
	if v, ok := component.(IOpenable); ok {
		return v.Open(ctx)
	}
	return nil
}

// Open opens multiple components.
// To be opened components must implement IOpenable interface.
// If they don't the call to this method has no effect.
//	see OpenOne
//	see IOpenable
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- components []any the list of components that are to be closed.
//	Returns: error
func (c *_TOpener) Open(ctx context.Context, components []any) error {
	for _, component := range components {
		if err := c.OpenOne(ctx, component); err != nil {
			return err
		}
	}
	return nil
}
