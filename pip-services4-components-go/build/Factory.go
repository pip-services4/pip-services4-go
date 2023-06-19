package build

// Basic component factory that creates components using registered types and factory functions.
//	Example:
//		factory := NewFactory();
//		factory.RegisterType(
//			NewDescriptor("mygroup", "mycomponent1", "default", "*", "1.0"),
//			MyComponent1
//		);
//		factory.Register(
//			NewDescriptor("mygroup", "mycomponent2", "default", "*", "1.0"),
//			(locator){
//				return NewMyComponent2();
//			}
//		);
//
//		res, err := factory.Create(NewDescriptor("mygroup", "mycomponent1", "default", "name1", "1.0"))
//		res, err := factory.Create(NewDescriptor("mygroup", "mycomponent2", "default", "name2", "1.0"))

import (
	refl "reflect"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

type registration struct {
	locator any
	factory func(any) any
}

type Factory struct {
	_registrations []*registration
}

// NewFactory create new factory
//
//	Returns: *Factory
func NewFactory() *Factory {
	return &Factory{
		_registrations: []*registration{},
	}
}

// Register registers a component using a factory method.
//
//	Parameters:
//		- locator any a locator to identify component to be created.
//		- factory func(locator any) any a factory function that receives a locator and returns a created component.
func (c *Factory) Register(locator any, factory func(locator any) any) {
	if locator == nil {
		panic("Locator cannot be nil")
	}
	if factory == nil {
		panic("Factory cannot be nil")
	}

	c._registrations = append(c._registrations, &registration{
		locator: locator,
		factory: factory,
	})
}

// RegisterType registers a component using its type (a constructor function).
//
//	Parameters:
//		- locator any a locator to identify component to be created.
//		- factory any a factory.
func (c *Factory) RegisterType(locator any, factory any) {
	if locator == nil {
		panic("Locator cannot be nil")
	}
	if factory == nil {
		panic("Factory cannot be nil")
	}

	val := refl.ValueOf(factory)
	if val.Kind() != refl.Func {
		panic("Factory must be parameterless function")
	}

	c.Register(locator, func(locator any) any {
		return val.Call([]refl.Value{})[0].Interface()
	})
}

// CanCreate checks if this factory is able to create component by given locator.
// This method searches for all registered components and returns a locator for
// component it is able to create that matches the given locator.
// If the factory is not able to create a requested component is returns null.
//
//	Parameters:
//		- locator any a locator to identify component to be created.
//	Returns: any a locator for a component that the factory is able to create.
func (c *Factory) CanCreate(locator any) any {
	for _, registration := range c._registrations {
		thisLocator := registration.locator

		if equatable, ok := thisLocator.(data.IEquatable[any]); ok && equatable.Equals(locator) {
			return thisLocator
		}

		if thisLocator == locator {
			return thisLocator
		}
	}
	return nil
}

// Create a component identified by given locator.
//
//	Parameters:
//		- locator any a locator to identify component to be created.
//	Returns: any, error the created component and a CreateError if the factory
//		is not able to create the component.
func (c *Factory) Create(locator any) (any, error) {
	var factory func(any) any

	for _, registration := range c._registrations {
		thisLocator := registration.locator

		if equatable, ok := thisLocator.(data.IEquatable[any]); ok && equatable.Equals(locator) {
			factory = registration.factory
			break
		}

		if thisLocator == locator {
			factory = registration.factory
			break
		}
	}

	if factory == nil {
		return nil, NewCreateErrorByLocator("", locator)
	}

	var err error

	obj := func() any {
		defer func() {
			if r := recover(); r != nil {
				tempMessage := convert.StringConverter.ToString(r)
				tempError := NewCreateError("", tempMessage)

				if cause, ok := r.(error); ok {
					_ = tempError.WithCause(cause)
				}

				err = tempError
			}
		}()

		return factory(locator)
	}()

	return obj, err
}
