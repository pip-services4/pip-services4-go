package refer

import (
	"context"

	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/run"
)

// RunReferencesDecorator References decorator that automatically opens
// to newly added components that implement IOpenable interface and
// closes removed components that implement ICloseable interface.
type RunReferencesDecorator struct {
	*ReferencesDecorator
	opened bool
}

// NewRunReferencesDecorator creates a new instance of the decorator.
//
//	Parameters:
//		- nextReferences crefer.IReferences the next references or decorator in the chain.
//		- topReferences crefer.IReferences the decorator at the top of the chain.
//	Returns: *RunReferencesDecorator
func NewRunReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *RunReferencesDecorator {

	return &RunReferencesDecorator{
		ReferencesDecorator: NewReferencesDecorator(nextReferences, topReferences),
	}
}

// IsOpen checks if the component is opened.
//
//	Returns: bool true if the component has been opened and false otherwise.
func (c *RunReferencesDecorator) IsOpen() bool {
	return c.opened
}

// Open the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *RunReferencesDecorator) Open(ctx context.Context) error {
	if !c.opened {
		components := c.GetAll()
		err := run.Opener.Open(ctx, components)
		c.opened = err == nil
		return err
	}
	return nil
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *RunReferencesDecorator) Close(ctx context.Context) error {
	components := c.GetAll()
	err := run.Closer.Close(ctx, components)
	c.opened = false
	return err
}

// Put a new reference into this reference map.
//
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to find the reference by.
//		- component any a component reference to be added.
func (c *RunReferencesDecorator) Put(ctx context.Context, locator any, component any) {
	c.ReferencesDecorator.Put(ctx, locator, component)

	if c.opened {
		_ = run.Opener.OpenOne(ctx, component)
	}
}

// Remove a previously added reference that matches specified locator.
// If many references match the locator, it removes only the first one.
// When all references shall be removed, use removeAll method instead.
//
//	see RemoveAll
//	Parameters:
//		- ctx context.Context
//		- locator any the locator to remove references by.
//	Returns: any the removed component reference.
func (c *RunReferencesDecorator) Remove(ctx context.Context, locator any) any {
	component := c.ReferencesDecorator.Remove(ctx, locator)

	if c.opened {
		_ = run.Closer.CloseOne(ctx, component)
	}

	return component
}

// RemoveAll all component references that match the specified locator.
//
//	Parameters:
//		- ctx context.Context
//		- locator any the locator to remove references by.
//	Returns: []any a list, containing all removed references.
func (c *RunReferencesDecorator) RemoveAll(ctx context.Context, locator any) []any {
	components := c.NextReferences.RemoveAll(ctx, locator)

	if c.opened {
		_ = run.Closer.Close(ctx, components)
	}

	return components
}
