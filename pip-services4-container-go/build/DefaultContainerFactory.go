package build

// Creates default container components (loggers, counters, caches, locks, etc.) by their descriptors.
import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cfgbuild "github.com/pip-services4/pip-services4-go/pip-services4-config-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/test"
	lbuild "github.com/pip-services4/pip-services4-go/pip-services4-logic-go/build"
	obuild "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/build"
)

// NewDefaultContainerFactory create a new instance of the factory and sets nested factories.
//
//	Returns: *DefaultContainerFactory
func NewDefaultContainerFactory() *cbuild.CompositeFactory {
	c := cbuild.NewCompositeFactory()

	c.Add(context.NewDefaultContextFactory())
	c.Add(obuild.NewDefaultObservabilityFactory())
	c.Add(lbuild.NewDefaultLogicFactory())
	c.Add(cfgbuild.NewDefaultConfigFactory())

	c.Add(test.NewDefaultTestFactory())

	return c
}

// NewDefaultContainerFactoryFromFactories create a new instance of the factory and sets nested factories.
//
//	Parameters:
//		- factories ...cbuild.IFactory a list of nested factories
//	Returns: *cbuild.CompositeFactory
func NewDefaultContainerFactoryFromFactories(factories ...cbuild.IFactory) *cbuild.CompositeFactory {
	c := NewDefaultContainerFactory()

	for _, factory := range factories {
		c.Add(factory)
	}

	return c
}
