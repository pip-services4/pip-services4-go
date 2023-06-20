package refer

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// BuildReferencesDecorator references decorator that automatically creates missing components using available
// component factories upon component retrival.
type BuildReferencesDecorator struct {
	*ReferencesDecorator
}

// NewBuildReferencesDecorator creates a new instance of the decorator.
//
//	Parameters:
//		- nextReferences crefer.IReferences the next references or decorator in the chain.
//		- topReferences IReferences the decorator at the top of the chain.
//	Returns: *BuildReferencesDecorator
func NewBuildReferencesDecorator(nextReferences crefer.IReferences,
	topReferences crefer.IReferences) *BuildReferencesDecorator {

	return &BuildReferencesDecorator{
		ReferencesDecorator: NewReferencesDecorator(nextReferences, topReferences),
	}
}

// FindFactory finds a factory capable creating component by given descriptor
// from the components registered in the references.
//
//	Parameters:
//		- locator any a locator of component to be created.
//	Returns: build.IFactory found factory or nil if factory was not found.
func (c *BuildReferencesDecorator) FindFactory(locator any) build.IFactory {
	components := c.GetAll()

	for _, component := range components {
		if factory, ok := component.(build.IFactory); ok && factory.CanCreate(locator) != nil {
			return factory
		}
	}

	return nil
}

// Create creates a component identified by given locator.
//
//	throws a CreateEerror if the factory is not able to create the component.
//	see FindFactory
//	Parameters:
//		- locator any a locator to identify component to be created.
//		- factory build.IFactory a factory that shall create the component.
//	Returns: any the created component.
func (c *BuildReferencesDecorator) Create(locator any,
	factory build.IFactory) any {

	if factory == nil {
		return nil
	}

	var result any

	defer func() {
		recover()
	}()

	result, _ = factory.Create(locator)

	return result
}

// ClarifyLocator a component locator by merging two descriptors into one to replace missing fields.
// That allows to get a more complete descriptor that includes all possible fields.
//
//	Parameters:
//		- locator any a component locator to clarify.
//		- factory build.IFactory a factory that shall create the component.
//	Returns: any clarified component descriptor (locator)
func (c *BuildReferencesDecorator) ClarifyLocator(locator any,
	factory build.IFactory) any {

	if factory == nil {
		return nil
	}

	descriptor, ok := locator.(*crefer.Descriptor)
	if !ok {
		return locator
	}

	anotherLocator := factory.CanCreate(locator)
	anotherDescriptor, ok := anotherLocator.(*crefer.Descriptor)
	if !ok {
		return locator
	}

	group := descriptor.Group()
	if group == "" {
		group = anotherDescriptor.Group()
	}
	typ := descriptor.Type()
	if typ == "" {
		typ = anotherDescriptor.Type()
	}
	kind := descriptor.Kind()
	if kind == "" {
		kind = anotherDescriptor.Kind()
	}
	name := descriptor.Name()
	if name == "" {
		name = anotherDescriptor.Name()
	}
	version := descriptor.Version()
	if version == "" {
		version = anotherDescriptor.Version()
	}

	return crefer.NewDescriptor(group, typ, kind, name, version)
}

// GetOneOptional gets an optional component reference that matches specified locator.
//
//	Parameters:
//		- locator any the locator to find references by.
//	Returns: any a matching component reference or nil if nothing was found.
func (c *BuildReferencesDecorator) GetOneOptional(locator any) any {
	components, err := c.Find(locator, false)
	if err != nil || len(components) == 0 {
		return nil
	}
	return components[0]
}

// GetOneRequired a required component reference that matches specified locator.
//
//	throws a ReferenceException when no references found.
//	Parameters:
//		- locator any the locator to find a reference by.
//	Returns: any, error a matching component reference and error.
func (c *BuildReferencesDecorator) GetOneRequired(locator any) (any, error) {
	components, err := c.Find(locator, true)
	if err != nil || len(components) == 0 {
		return nil, err
	}
	return components[0], nil
}

// GetOptional all component references that match specified locator.
//
//	Parameters:
//		- locator any the locator to find references by.
//	Returns: []any a list with matching component references or empty list if nothing was found.
func (c *BuildReferencesDecorator) GetOptional(locator any) []any {
	components, _ := c.Find(locator, false)
	return components
}

// GetRequired all component references that match specified locator.
// At least one component reference must be present.
// If it doesn't the method throws an error.
//
//	throws a ReferenceException when no references found.
//	Parameters:
//		- locator any the locator to find references by.
//	Returns: []any, erorr a list with matching component references and error.
func (c *BuildReferencesDecorator) GetRequired(locator any) ([]any, error) {
	return c.Find(locator, true)
}

// Find all component references that match specified locator.
//
//	throws a ReferenceError when required is set to true but no references found.
//	Parameters:
//		- locator interface the locator to find a reference by.
//		- required bool forces to raise an exception if no reference is found.
//	Returns: []interface, error a list with matching component references and error.
func (c *BuildReferencesDecorator) Find(locator any, required bool) ([]any, error) {
	components, _ := c.ReferencesDecorator.Find(locator, required)

	if required && len(components) == 0 {
		factory := c.FindFactory(locator)
		component := c.Create(locator, factory)
		if component != nil {
			locator = c.ClarifyLocator(locator, factory)
			// TODO:: check ctx propagation
			c.ReferencesDecorator.TopReferences.Put(context.TODO(), locator, component)
			components = append(components, component)
		}
	}

	if required && len(components) == 0 {
		err := crefer.NewReferenceError(context.Background(), locator)
		return nil, err
	}

	return components, nil
}
