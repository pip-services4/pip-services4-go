package refer

import (
	"context"

	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// ManagedReferences managed references that in addition to keeping and locating
// references can also manage their lifecycle:
//
//	Auto-creation of missing component using available factories
//	Auto-linking newly added components
//	Auto-opening newly added components
//	Auto-closing removed components
type ManagedReferences struct {
	*ReferencesDecorator
	References *crefer.References
	Builder    *BuildReferencesDecorator
	Linker     *LinkReferencesDecorator
	Runner     *RunReferencesDecorator
}

// NewManagedReferences creates a new instance of the references
//
//	Parameters:
//		- ctx context.Context
//		- tuples []any tuples where odd values are component locators
//			(descriptors) and even values are component references
//	Returns: *ManagedReferences
func NewManagedReferences(ctx context.Context, tuples []any) *ManagedReferences {
	c := &ManagedReferences{
		ReferencesDecorator: NewReferencesDecorator(nil, nil),
	}

	c.References = crefer.NewReferences(ctx, tuples)
	c.Builder = NewBuildReferencesDecorator(c.References, c)
	c.Linker = NewLinkReferencesDecorator(c.Builder, c)
	c.Runner = NewRunReferencesDecorator(c.Linker, c)

	c.ReferencesDecorator.NextReferences = c.Runner

	return c
}

// NewEmptyManagedReferences creates a new instance of the references
//
//	Returns: *ManagedReferences
func NewEmptyManagedReferences() *ManagedReferences {
	return NewManagedReferences(context.Background(), []any{})
}

// NewManagedReferencesFromTuples creates a new ManagedReferences object
// filled with provided key-value pairs called tuples. Tuples parameters contain a
// sequence of locator1, component1, locator2, component2, ... pairs.
//
//	Parameters:
//		- ctx context.Context
//		- tuples ...any the tuples to fill a new ManagedReferences object.
//	Returns: *ManagedReferences a new ManagedReferences object.
func NewManagedReferencesFromTuples(ctx context.Context, tuples ...any) *ManagedReferences {
	return NewManagedReferences(ctx, tuples)
}

// IsOpen checks if the component is opened.
//
//	Returns: bool true if the component has been opened and false otherwise.
func (c *ManagedReferences) IsOpen() bool {
	return c.Linker.IsOpen() && c.Runner.IsOpen()
}

// Open the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *ManagedReferences) Open(ctx context.Context) error {
	err := c.Linker.Open(ctx)
	if err == nil {
		err = c.Runner.Open(ctx)
	}
	return err
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *ManagedReferences) Close(ctx context.Context) error {
	err := c.Runner.Close(ctx)
	if err == nil {
		err = c.Linker.Close(ctx)
	}
	return err
}
