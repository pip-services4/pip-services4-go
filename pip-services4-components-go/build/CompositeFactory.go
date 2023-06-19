package build

// CompositeFactory aggregates multiple factories into a single factory component.
// When a new component is requested, it iterates through factories to locate
// the one able to create the requested component.
// This component is used to conveniently keep all supported factories in a single place.
//	Example:
//		factory := NewCompositeFactory();
//		factory.Add(NewDefaultLoggerFactory());
//		factory.Add(NewDefaultCountersFactory());
//
//		loggerLocator := NewDescriptor("*", "logger", "*", "*", "1.0");
//		factory.CanCreate(context.Background(), loggerLocator);
//			Result: Descriptor("pip-service", "logger", "null", "default", "1.0")
//		factory.Create(context.Background(), loggerLocator);    // Result: created NullLogger
type CompositeFactory struct {
	_factories []IFactory
}

// NewCompositeFactory creates a new instance of the factory.
// Returns: *CompositeFactory
func NewCompositeFactory() *CompositeFactory {
	return &CompositeFactory{
		_factories: make([]IFactory, 0),
	}
}

// NewCompositeFactoryFromFactories creates a new instance of the factory.
//	Parameters:
//		- factories ...IFactory a list of factories to embed into this factory.
//	Returns: *CompositeFactory
func NewCompositeFactoryFromFactories(factories ...IFactory) *CompositeFactory {
	return &CompositeFactory{
		_factories: factories,
	}
}

// Add a factory into the list of embedded factories.
//	Parameters:
//		- factory IFactory a factory to be added.
func (c *CompositeFactory) Add(factory IFactory) {
	if factory == nil {
		panic("Factory cannot be nil")
	}

	c._factories = append(c._factories, factory)
}

// Remove removes a factory from the list of embedded factories.
//	Parameters:
//		- factory IFactory the factory to remove.
func (c *CompositeFactory) Remove(factory IFactory) {
	removeIndex := -1
	for i, thisFactory := range c._factories {
		if thisFactory == factory {
			removeIndex = i
			break
		}
	}
	if removeIndex >= 0 {
		c._factories = append(c._factories[:removeIndex], c._factories[removeIndex+1:]...)
	}
}

// CanCreate checks if this factory is able to create component by given locator.
// This method searches for all registered components and returns a locator for component
// it is able to create that matches the given locator. If the factory is not able
// to create a requested component is returns null.
//	Parameters:
//		- locator any a locator to identify component to be created.
//	Returns: any  a locator for a component that the factory is able to create.
func (c *CompositeFactory) CanCreate(locator any) any {
	if locator == nil {
		panic("Locator cannot be null")
	}

	// Iterate from the latest factories
	for i := len(c._factories) - 1; i >= 0; i-- {
		if thisLocator := c._factories[i].CanCreate(locator); thisLocator != nil {
			return thisLocator
		}
	}

	return nil
}

// Create creates a component identified by given locator.
//	Parameters:
//		- locator any a locator to identify component to be created.
//	Returns: any, error the created component and a
//		CreateError if the factory is not able to create the component..
func (c *CompositeFactory) Create(locator any) (any, error) {
	if locator == nil {
		panic("Locator cannot be nil")
	}

	// Iterate from the latest _factories
	for i := len(c._factories) - 1; i >= 0; i-- {
		factory := c._factories[i]
		if factory.CanCreate(locator) != nil {
			return factory.Create(locator)
		}
	}

	return nil, NewCreateErrorByLocator("", locator)
}
