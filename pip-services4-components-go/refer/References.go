package refer

import "context"

// References the most basic implementation of IReferences to store and locate component references.
//	see IReferences
//	Example:
//		type MyController  {
//			_persistence IMyPersistence;
//		}
//		func (mc *MyController) setReferences(references IReferences) {
//			mc._persistence = references.GetOneRequired(
//				NewDescriptor("mygroup", "persistence", "*", "*", "1.0"),
//			);
//		}
//
//		persistence := NewMyMongoDbPersistence();
//		controller := MyController();
//
//		references := NewReferencesFromTuples(
//			new Descriptor("mygroup", "persistence", "mongodb", "default", "1.0"), persistence,
//			new Descriptor("mygroup", "controller", "default", "default", "1.0"), controller
//		);
//		controller.setReferences(references);
type References struct {
	references []*Reference
}

// NewEmptyReferences creates a new instance of references and initializes it with references.
//	Returns: *References
func NewEmptyReferences() *References {
	return &References{
		references: make([]*Reference, 0, 10),
	}
}

// NewReferences creates a new instance of references and initializes it with references.
//	Parameters:
//		- ctx context.Context
//		- tuples []any a list of values where odd
//			elements are locators and the following even elements are component references
//	Returns: *References
func NewReferences(ctx context.Context, tuples []any) *References {
	c := NewEmptyReferences()

	for index := 0; index < len(tuples); index += 2 {
		if index+1 >= len(tuples) {
			break
		}
		c.Put(ctx, tuples[index], tuples[index+1])
	}

	return c
}

// NewReferencesFromTuples creates a new References from a list of key-value pairs called tuples.
//	Parameters:
//		- ctx context.Context
//		- tuples  ...any a list of values where
//			odd elements are locators and the following even elements
//			are component references
//	Returns: *References a newly created References.
func NewReferencesFromTuples(ctx context.Context, tuples ...any) *References {
	return NewReferences(ctx, tuples)
}

// Put a new reference into this reference map.
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to find the reference by.
//		- component any a component reference to be added.
func (c *References) Put(ctx context.Context, locator any, component any) {
	if component == nil {
		panic("Component cannot be null")
	}

	reference := NewReference(locator, component)
	c.references = append(c.references, reference)
}

// Remove a previously added reference that matches specified locator.
// If many references match the locator, it removes only the first one.
// When all references shall be removed, use removeAll method instead.
//	see RemoveAll
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to remove reference
//	Returns: any the removed component reference.
func (c *References) Remove(ctx context.Context, locator any) any {
	if locator == nil {
		return nil
	}

	for index := len(c.references) - 1; index >= 0; index-- {
		reference := c.references[index]
		if reference.Match(locator) {
			c.references = append(c.references[:index], c.references[index+1:]...)
			return reference.Component()
		}
	}

	return nil
}

// RemoveAll removes all component references that match the specified locator.
//	Parameters:
//		- ctx context.Context
//		- locator any a locator to remove reference
//	Returns: []any a list, containing all removed references.
func (c *References) RemoveAll(ctx context.Context, locator any) []any {
	components := make([]any, 0, 5)

	if locator == nil {
		return components
	}

	for index := len(c.references) - 1; index >= 0; index-- {
		reference := c.references[index]
		if reference.Match(locator) {
			c.references = append(c.references[:index], c.references[index+1:]...)
			components = append(components, reference.Component())
		}
	}

	return components
}

// GetAllLocators gets locators for all registered component references in this reference map.
//	Returns: []any a list with component locators.
func (c *References) GetAllLocators() []any {
	components := make([]any, len(c.references))

	for index, reference := range c.references {
		components[index] = reference.Locator()
	}

	return components
}

// GetAll gets all component references registered in this reference map.
//	Returns: []any a list with component references.
func (c *References) GetAll() []any {
	components := make([]any, len(c.references))

	for index, reference := range c.references {
		components[index] = reference.Component()
	}

	return components
}

// GetOneOptional gets an optional component reference that matches specified locator.
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: any a matching component reference or nil if nothing was found.
func (c *References) GetOneOptional(locator any) any {
	components, err := c.Find(locator, false)
	if err != nil || len(components) == 0 {
		return nil
	}
	return components[0]
}

// GetOneRequired gets a required component reference that matches specified locator.
// throws a ReferenceError when no references found.
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: any a matching component reference.
func (c *References) GetOneRequired(locator any) (any, error) {
	components, err := c.Find(locator, true)
	if err != nil || len(components) == 0 {
		return nil, err
	}
	return components[0], nil
}

// GetOptional gets all component references that match specified locator.
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: []any a list with matching component references or
//		empty list if nothing was found.
func (c *References) GetOptional(locator any) []any {
	components, _ := c.Find(locator, false)
	return components
}

// GetRequired gets all component references that match specified locator.
// At least one component reference must be present.
// If it doesn't the method throws an error.
// throws a ReferenceError when no references found.
//	Parameters:
//		- locator any a locator to remove reference
//	Returns: []any a list with matching component references.
func (c *References) GetRequired(locator any) ([]any, error) {
	return c.Find(locator, true)
}

// Find gets all component references that match specified locator.
// throws a ReferenceError when required is set to true but no references found.
//	Parameters:
//		- locator any the locator to find a reference by.
//		- required bool forces to raise an exception if no reference is found.
//	Returns: []any a list with matching component references.
func (c *References) Find(locator any, required bool) ([]any, error) {
	if locator == nil {
		panic("Locator cannot be null")
	}

	components := make([]any, 0, 2)

	// Search all references
	for index := len(c.references) - 1; index >= 0; index-- {
		reference := c.references[index]
		if reference.Match(locator) {
			component := reference.Component()
			components = append(components, component)
		}
	}

	if len(components) == 0 && required {
		err := NewReferenceError(context.Background(), locator)
		return components, err
	}

	return components, nil
}
