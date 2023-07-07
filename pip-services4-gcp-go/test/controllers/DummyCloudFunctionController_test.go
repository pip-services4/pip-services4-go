package services_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyCloudFunctionControllerTest struct {
	fixture       *DummyCloudFunctionFixture
	funcContainer *DummyCloudFunction
}

func newDummyCloudFunctionControllerTest() *DummyCloudFunctionControllerTest {
	return &DummyCloudFunctionControllerTest{}
}

func (c *DummyCloudFunctionControllerTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"controller.descriptor", "pip-services-dummies:controller:cloudfunc:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyCloudFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx)
	assert.Nil(t, err)

	c.fixture = NewDummyCloudFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyCloudFunctionControllerTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsCloudController(t *testing.T) {
	c := newDummyCloudFunctionControllerTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
