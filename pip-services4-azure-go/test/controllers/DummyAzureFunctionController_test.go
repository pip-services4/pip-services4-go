package controllers_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyAzureFunctionControllerTest struct {
	fixture       *DummyAzureFunctionFixture
	funcContainer *DummyAzureFunction
}

func newDummyAzureFunctionControllerTest() *DummyAzureFunctionControllerTest {
	return &DummyAzureFunctionControllerTest{}
}

func (c *DummyAzureFunctionControllerTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:controller:azurefunc:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyAzureFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx)
	assert.Nil(t, err)

	c.fixture = NewDummyAzureFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyAzureFunctionControllerTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsAzureService(t *testing.T) {
	c := newDummyAzureFunctionControllerTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
