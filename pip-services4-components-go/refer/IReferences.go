package refer

import "context"

// IReferences interface for a map that holds component references and passes them to
// components to establish dependencies with each other.
// Together with IReferenceable and IUnreferenceable interfaces it implements a
// Locator pattern that is used by PipServices toolkit for Inversion of Control to assign external
// dependencies to components.
// The IReferences object is a simple map, where keys are locators and values are component references.
// It allows to add, remove and find components by their locators. Locators can be any values like integers,
// strings or component types. But most often PipServices toolkit uses
// Descriptor as locators that match by 5 fields: group, type, kind, name and version.
type IReferences interface {

	// Put a new reference into this reference map.
	//	Parameters:
	//		- ctx context.Context
	//		- locator any a locator to find the reference by.
	//		- component any a component reference to be added.
	//	Returns: any
	Put(ctx context.Context, locator any, component any)

	// Remove a previously added reference that matches specified locator. If many references match the locator, it removes only the first one. When all references shall be removed, use removeAll method instead.
	//	see RemoveAll
	//	Parameters:
	//		- ctx context.Context
	//		- locator any a locator to remove reference
	//	Returns: any the removed component reference.
	Remove(ctx context.Context, locator any) any

	// RemoveAll removes all component references that match the specified locator.
	//	Parameters:
	//		- ctx context.Context
	//		- locator any the locator to remove references by.
	//	Returns: []any a list, containing all removed references.
	RemoveAll(ctx context.Context, locator any) []any

	// GetAllLocators gets locators for all registered component references in this reference map.
	//	Returns: []any a list with component locators.
	GetAllLocators() []any

	// GetAll gets all component references registered in this reference map.
	//	Returns: []any a list with component references.
	GetAll() []any

	// GetOptional gets all component references that match specified locator.
	//	Parameters:
	//		- locator any the locator to find references by.
	//	Returns: []any a list with matching component references or empty list if nothing was found.
	GetOptional(locator any) []any

	// GetRequired gets all component references that match specified locator. At least one component reference must be present. If it doesn't the method throws an error.
	// throws a ReferenceException when no references found.
	//	Parameters:
	//		- locator any the locator to find references by.
	//	Returns []any a list with matching component references.
	GetRequired(locator any) ([]any, error)

	// GetOneOptional gets an optional component reference that matches specified locator.
	//	Parameters:
	//		- locator any the locator to find references by.
	//	Returns: any a matching component reference or nil if nothing was found.
	GetOneOptional(locator any) any

	// GetOneRequired gets a required component reference that matches specified locator.
	// throws a ReferenceError when no references found.
	//	Parameters:
	//		- locator any the locator to find references by.
	//	Returns: any a matching component reference.
	GetOneRequired(locator any) (any, error)

	// Find gets all component references that match specified locator.
	// throws a ReferenceError when required is set to true but no references found.
	//	Parameters:
	//		- locator any the locator to find a reference by.
	//		- required bool forces to raise an exception if no reference is found.
	//	Returns: []any a list with matching component references.
	Find(locator any, required bool) ([]any, error)
}
