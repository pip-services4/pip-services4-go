package refer

import (
	"context"
	"fmt"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/config"
)

// ContainerReferences container managed references that can be
// created from container configuration.
type ContainerReferences struct {
	*ManagedReferences
}

// NewContainerReferences creates a new instance of the references
// Returns *ContainerReferences
func NewContainerReferences() *ContainerReferences {
	return &ContainerReferences{
		ManagedReferences: NewEmptyManagedReferences(),
	}
}

// PutFromConfig puts components into the references from container configuration.
//
//	Parameters:
//		- ctx context.Context
//		- config config.ContainerConfig a container
//			configuration with information of components to be added.
//	Returns: error CreateError when one of component cannot be created.
func (c *ContainerReferences) PutFromConfig(ctx context.Context, config config.ContainerConfig) error {
	var err error
	var locator any
	var component any

	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()

	for _, componentConfig := range config {
		if componentConfig.Type != nil {
			// Create component dynamically
			locator = componentConfig.Type
			component, err = reflect.TypeReflector.CreateInstanceByDescriptor(componentConfig.Type)
		} else if componentConfig.Descriptor != nil {
			// Or create component statically
			locator = componentConfig.Descriptor
			factory := c.ManagedReferences.Builder.FindFactory(locator)
			component = c.ManagedReferences.Builder.Create(locator, factory)
			if component == nil {
				return refer.NewReferenceError(ctx, locator)
			}
			locator = c.ManagedReferences.Builder.ClarifyLocator(locator, factory)
		}

		// Check that component was created
		if component == nil {
			return build.NewCreateError(
				"CANNOT_CREATE_COMPONENT",
				"Cannot create component",
			).WithDetails("config", config)
		}

		fmt.Printf("Created component %v\n", locator)

		// Add component to the list
		c.ManagedReferences.References.Put(ctx, locator, component)

		// Configure component
		if configurable, ok := component.(cconfig.IConfigurable); ok {
			configurable.Configure(ctx, componentConfig.Config)
		}

		// Set references to factories
		if _, ok := component.(build.IFactory); ok {
			if referenceable, ok := component.(refer.IReferenceable); ok {
				referenceable.SetReferences(ctx, c)
			}
		}
	}

	return err
}
