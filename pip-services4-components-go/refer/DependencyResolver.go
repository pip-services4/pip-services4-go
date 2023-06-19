package refer

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	conf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// DependencyResolver helper class for resolving component dependencies.
// The resolver is configured to resolve named dependencies by specific locator.
// During deployment the dependency locator can be changed.
// This mechanism can be used to clarify specific dependency among several alternatives.
// Typically components are configured to retrieve the first dependency that matches
// logical group, type and version. But if container contains more than one instance
// and resolution has to be specific about those instances, they can be given a unique
// name and dependency resolvers can be reconfigured to retrieve dependencies by their name.
//
//	Configuration parameters dependencies:
//		[dependency name 1]: Dependency 1 locator (descriptor)
//		...
//		[dependency name N]: Dependency N locator (descriptor)
//	References must match configured dependencies.
type DependencyResolver struct {
	dependencies map[string]any
	references   IReferences
}

// NewDependencyResolver creates a new instance of the dependency resolver.
//
//	Returns: *DependencyResolver
func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		dependencies: make(map[string]any),
		references:   nil,
	}
}

// NewDependencyResolverWithParams creates a new instance of the dependency resolver.
//
//	see ConfigParams
//	see Configure
//	see IReferences
//	see SetReferences
//	Parameters:
//		- ctx context.Context
//		- config ConfigParams default configuration where key is
//			dependency name and value is locator (descriptor)
//		- references IReferences default component references
//	Returns: *DependencyResolver
func NewDependencyResolverWithParams(ctx context.Context,
	config *conf.ConfigParams, references IReferences) *DependencyResolver {

	c := NewDependencyResolver()

	if config != nil {
		c.Configure(ctx, config)
	}

	if references != nil {
		c.SetReferences(ctx, references)
	}

	return c
}

// NewDependencyResolverFromTuples creates a new DependencyResolver from a list of key-value pairs
// called tuples where key is dependency name and value the dependency locator (descriptor).
//
//	see NewDependencyResolverFromTuplesArray
//	Parameters:
//		- ctx context.Context
//		- tuples ...any a list of values where odd elements are
//		dependency name and the following even elements
//		are dependency locator (descriptor)
//	Returns: *DependencyResolver a newly created DependencyResolver.
func NewDependencyResolverFromTuples(ctx context.Context, tuples ...any) *DependencyResolver {
	result := NewDependencyResolver()
	if len(tuples) == 0 {
		return result
	}

	for index := 0; index < len(tuples); index += 2 {
		if index+1 >= len(tuples) {
			break
		}

		name := convert.StringConverter.ToString(tuples[index])
		locator := tuples[index+1]

		result.Put(ctx, name, locator)
	}

	return result
}

// Configure the component with specified parameters.
//
//	see ConfigParams
//	Parameters:
//		- ctx context.Context
//		- config *conf.ConfigParams configuration parameters to set.
func (c *DependencyResolver) Configure(ctx context.Context, config *conf.ConfigParams) {
	dependencies := config.GetSection("dependencies")
	names := dependencies.Keys()
	for _, name := range names {
		if locator, ok := dependencies.GetAsNullableString(name); ok {
			descriptor, err := ParseDescriptorFromString(locator)
			if err == nil {
				c.dependencies[name] = descriptor
			} else {
				c.dependencies[name] = locator
			}
		}
	}
}

// SetReferences sets the component references. References must match configured dependencies.
//
//	Parameters:
//		- ctx context.Context
//		- references IReferences references to set.
func (c *DependencyResolver) SetReferences(ctx context.Context, references IReferences) {
	c.references = references
}

// Put adds a new dependency into this resolver.
//
//	Parameters:
//		- ctx context.Context
//		- name string the dependency's name.
//		- locator any the locator to find the dependency by.
func (c *DependencyResolver) Put(ctx context.Context, name string, locator any) {
	c.dependencies[name] = locator
}

// Locate dependency by name
//
//	Parameters: name string dependency name
//	Returns: any located dependency
func (c *DependencyResolver) Locate(name string) any {
	if name == "" {
		panic("Dependency name cannot be empty")
	}

	if c.references == nil {
		panic("References shall be set")
	}

	return c.dependencies[name]
}

// GetOptional gets all optional dependencies by their name.
//
//	Parameters: name string the dependency name to locate.
//	Returns:
//		- []any a list with found dependencies or
//		empty list of no dependencies was found.
func (c *DependencyResolver) GetOptional(name string) []any {
	locator := c.Locate(name)
	if locator == nil {
		return make([]any, 0)
	}
	return c.references.GetOptional(locator)
}

// GetRequired gets all required dependencies by their name. At least one dependency must be present.
// If no dependencies was found it throws a ReferenceError
// throws a ReferenceError if no dependencies were found.
//
//	Parameters:
//		- name string the dependency name to locate.
//	Returns: []any a list with found dependencies.
func (c *DependencyResolver) GetRequired(name string) ([]any, error) {
	locator := c.Locate(name)
	if locator == nil {
		err := NewReferenceError(context.Background(), name)
		return make([]any, 0), err
	}

	return c.references.GetRequired(locator)
}

// GetOneOptional gets one optional dependency by its name.
//
//	Parameters:
//		- ctx context.Context
//		- name string the dependency name to locate.
//	Returns: any a dependency reference or nil of the dependency was not found
func (c *DependencyResolver) GetOneOptional(name string) any {
	locator := c.Locate(name)
	if locator == nil {
		return nil
	}
	return c.references.GetOneOptional(locator)
}

// GetOneRequired gets one required dependency by its name. At least one dependency must present.
// If the dependency was found it throws a ReferenceError
// throws a ReferenceError if dependency was not found.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- name string the dependency name to locate.
//	Returns: any, error a dependency reference and error
func (c *DependencyResolver) GetOneRequired(name string) (any, error) {
	locator := c.Locate(name)
	if locator == nil {
		err := NewReferenceError(context.Background(), name)
		return nil, err
	}
	return c.references.GetOneRequired(locator)
}

// Find all matching dependencies by their name.
// throws a ReferenceError of required is true and no dependencies found.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- name string the dependency name to locate.
//		- required bool true to raise an exception when no dependencies are found.
//	Returns: []any, error a list of found dependencies and error
func (c *DependencyResolver) Find(name string, required bool) ([]any, error) {
	if name == "" {
		panic("Name cannot be empty")
	}

	locator := c.Locate(name)
	if locator == nil {
		if required {
			err := NewReferenceError(context.Background(), name)
			return make([]any, 0), err
		}
		return make([]any, 0), nil
	}

	return c.references.Find(locator, required)
}
