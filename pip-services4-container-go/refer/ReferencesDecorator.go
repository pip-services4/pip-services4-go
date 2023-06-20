package refer

import (
	"context"

	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// ReferencesDecorator chainable decorator for IReferences that allows
// to inject additional capabilities such as
// automatic component creation, automatic registration and opening.
type ReferencesDecorator struct {
	NextReferences crefer.IReferences
	TopReferences  crefer.IReferences
}

// NewReferencesDecorator creates a new instance of the decorator.
//
//	Parameters:
//		- nextReferences crefer.IReferences the next references or decorator in the chain.
//		- topReferences crefer.IReferences the decorator at the top of the chain.
//	Returns: *ReferencesDecorator
func NewReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *ReferencesDecorator {

	c := &ReferencesDecorator{
		NextReferences: nextReferences,
		TopReferences:  topReferences,
	}

	if c.NextReferences == nil {
		c.NextReferences = topReferences
	}
	if c.TopReferences == nil {
		c.TopReferences = nextReferences
	}

	return c
}

// Put a new reference into this reference map.
//
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to find the reference by.
//		- component any a component reference to be added.
func (c *ReferencesDecorator) Put(ctx context.Context, locator any, component any) {
	c.NextReferences.Put(ctx, locator, component)
}

// Remove a previously added reference that matches specified locator.
// If many references match the locator, it removes only the first one.
// When all references shall be removed, use removeAll method instead.
// see RemoveAll
//
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to remove reference
//	Returns: any the removed component reference.
func (c *ReferencesDecorator) Remove(ctx context.Context, locator any) any {
	return c.NextReferences.Remove(ctx, locator)
}

// RemoveAll all component references that match the specified locator.
//
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to remove reference
//	Returns: []any a list, containing all removed references.
func (c *ReferencesDecorator) RemoveAll(ctx context.Context, locator any) []any {
	return c.NextReferences.RemoveAll(ctx, locator)
}

// GetAllLocators locators for all registered component references in this reference map.
//
//	Returns: []any a list with component locators.
func (c *ReferencesDecorator) GetAllLocators() []any {
	return c.NextReferences.GetAllLocators()
}

// GetAll all component references registered in this reference map.
//
//	Returns: []any a list with component references.
func (c *ReferencesDecorator) GetAll() []any {
	return c.NextReferences.GetAll()
}

// GetOneOptional gets an optional component reference that matches specified locator.
//
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: any a matching component reference or null if nothing was found.
func (c *ReferencesDecorator) GetOneOptional(locator any) any {
	var component any

	defer func() {
		recover()
	}()

	components, err := c.Find(locator, false)
	if err == nil && len(components) > 0 {
		component = components[0]
	}

	return component
}

// GetOneRequired a required component reference that matches specified locator.
//
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: any, error a matching component reference, a ReferenceError when no references found.
func (c *ReferencesDecorator) GetOneRequired(locator any) (any, error) {
	components, err := c.Find(locator, true)
	if err != nil || len(components) == 0 {
		return nil, err
	}
	return components[0], nil
}

// GetOptional all component references that match specified locator.
//
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: []any a list with matching component references or empty list if nothing was found.
func (c *ReferencesDecorator) GetOptional(locator any) []any {
	components := make([]any, 0)

	defer func() {
		recover()
	}()

	components, _ = c.Find(locator, false)

	return components
}

// GetRequired all component references that match specified locator.
// At least one component reference must be present. If it doesn't the method throws an error.
//
//	Parameters:
//		- locator any a locator to remove reference
//	Returns []any a list with matching component references and
//		error a ReferenceError when no references found.
func (c *ReferencesDecorator) GetRequired(locator any) ([]any, error) {
	return c.Find(locator, true)
}

// Find all component references that match specified locator.
//
//	Parameters:
//		- locator any the locator to find a reference by.
//		- required bool forces to raise an exception if no reference is found.
//	Returns: []any, error a list with matching component references and
//		a ReferenceError when required is set to true but no references found
func (c *ReferencesDecorator) Find(locator any, required bool) ([]any, error) {
	return c.NextReferences.Find(locator, required)
}
