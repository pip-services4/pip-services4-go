package services_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyCommandableCloudFunctionControllerTest struct {
	fixture       *DummyCloudFunctionFixture
	funcContainer *DummyCloudFunction
}

func newDummyCommandableCloudFunctionControllerTest() *DummyCommandableCloudFunctionControllerTest {
	return &DummyCommandableCloudFunctionControllerTest{}
}

func (c *DummyCommandableCloudFunctionControllerTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"controller.descriptor", "pip-services-dummies:controller:commandable-cloudfunc:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyCloudFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx)
	assert.Nil(t, err)

	c.fixture = NewDummyCloudFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyCommandableCloudFunctionControllerTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsCommandableController(t *testing.T) {
	c := newDummyCommandableCloudFunctionControllerTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
