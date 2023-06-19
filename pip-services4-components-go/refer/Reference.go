package refer

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"

// Reference contains a reference to a component and locator to find it. It is used by
// References to store registered component references.
type Reference struct {
	locator   any
	component any
}

// NewReference create a new instance of the reference object and assigns its values.
// Parameters:
//		- locator any a locator to find the reference.
//		- component interface {}
//	Returns: *Reference
func NewReference(locator any, component any) *Reference {
	if component == nil {
		panic("Component cannot be null")
	}

	return &Reference{
		locator:   locator,
		component: component,
	}
}

// Component gets the stored component reference.
//	Returns: any the component's references.
func (c *Reference) Component() any {
	return c.component
}

// Locator gets the stored component locator.
//	Returns: any the component's locator.
func (c *Reference) Locator() any {
	return c.locator
}

// Match locator to this reference locator.
// Descriptors are matched using equal method. All other locator types are matched using direct comparison.
//	see Descriptor
//	Parameters:  locator any the locator to match.
//	Returns: bool true if locators are matching and false it they don't.
func (c *Reference) Match(locator any) bool {
	// Check for nil locator
	if locator == nil {
		return false
	}

	// Locate by direct reference matching
	if c.component == locator {
		return true
	}

	// Locate by direct locator matching
	if equatable, ok := c.locator.(data.IEquatable[any]); ok {
		return equatable.Equals(locator)
	}

	// Locate by direct locator matching
	return c.locator == locator
}
