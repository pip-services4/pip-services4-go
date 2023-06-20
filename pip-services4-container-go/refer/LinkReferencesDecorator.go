package refer

import (
	"context"

	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// LinkReferencesDecorator references decorator that automatically sets references
// to newly added components that implement IReferenceable
// interface and unsets references from removed components
// that implement IUnreferenceable interface.
type LinkReferencesDecorator struct {
	*ReferencesDecorator
	opened bool
}

// NewLinkReferencesDecorator creates a new instance of the decorator.
//
//	Parameters:
//		- nextReferences crefer.IReferences the next references or decorator in the chain.
//		- topReferences crefer.IReferences the decorator at the top of the chain.
//	Returns: *LinkReferencesDecorator
func NewLinkReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *LinkReferencesDecorator {
	return &LinkReferencesDecorator{
		ReferencesDecorator: NewReferencesDecorator(nextReferences, topReferences),
	}
}

// IsOpen checks if the component is opened.
//
//	Returns: bool true if the component has been opened and false otherwise.
func (c *LinkReferencesDecorator) IsOpen() bool {
	return c.opened
}

// Open the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *LinkReferencesDecorator) Open(ctx context.Context) error {
	if !c.opened {
		c.opened = true
		components := c.GetAll()
		crefer.Referencer.SetReferences(ctx, c.ReferencesDecorator.TopReferences, components)
	}
	return nil
}

// Close closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *LinkReferencesDecorator) Close(ctx context.Context) error {
	if c.opened {
		c.opened = false
		components := c.GetAll()
		crefer.Referencer.UnsetReferences(ctx, components)
	}
	return nil
}

// Put a new reference into this reference map.
//
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to find the reference by.
//		- component any a component reference to be added.
func (c *LinkReferencesDecorator) Put(ctx context.Context, locator any, component any) {
	c.ReferencesDecorator.Put(ctx, locator, component)

	if c.opened {
		crefer.Referencer.SetReferencesForOne(ctx, c.ReferencesDecorator.TopReferences, component)
	}
}

// Remove a previously added reference that matches specified locator.
// If many references match the locator, it removes only the first one.
// When all references shall be removed, use removeAll method instead.
//
//	see RemoveAll
//	Parameters:
//		- ctx context.Context
//		- locator interface a locator to remove reference
//	Returns: any the removed component reference.
func (c *LinkReferencesDecorator) Remove(ctx context.Context, locator any) any {
	component := c.ReferencesDecorator.Remove(ctx, locator)

	if c.opened {
		crefer.Referencer.UnsetReferencesForOne(ctx, component)
	}

	return component
}

// RemoveAll removes all component references that match the specified locator.
//
//	Parameters:
//		- ctx context.Context
//		- locator interface a locator to remove reference
//	Returns: []any a list, containing all removed references.
func (c *LinkReferencesDecorator) RemoveAll(ctx context.Context, locator any) []any {
	components := c.NextReferences.RemoveAll(ctx, locator)

	if c.opened {
		crefer.Referencer.UnsetReferences(ctx, components)
	}

	return components
}
