package controllers_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyCommandableAzureFunctionControllerTest struct {
	fixture       *DummyAzureFunctionFixture
	funcContainer *DummyAzureFunction
}

func newDummyCommandableAzureFunctionControllerTest() *DummyCommandableAzureFunctionControllerTest {
	return &DummyCommandableAzureFunctionControllerTest{}
}

func (c *DummyCommandableAzureFunctionControllerTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"controller.descriptor", "pip-services-dummies:controller:commandable-azurefunc:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyAzureFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx)
	assert.Nil(t, err)

	c.fixture = NewDummyAzureFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyCommandableAzureFunctionControllerTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsCommandableService(t *testing.T) {
	c := newDummyCommandableAzureFunctionControllerTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
