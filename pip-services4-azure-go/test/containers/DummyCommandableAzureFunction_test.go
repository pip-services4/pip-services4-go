package containers_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyCommandableAzureFunctionTest struct {
	fixture        *DummyAzureFunctionFixture
	funcContainers *DummyAzureFunction
}

func newDummyCommandableAzureFunctionTest() *DummyCommandableAzureFunctionTest {
	return &DummyCommandableAzureFunctionTest{}
}

func (c *DummyCommandableAzureFunctionTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
	)

	ctx := context.Background()

	c.funcContainers = NewDummyAzureFunction()
	c.funcContainers.Configure(ctx, config)
	err := c.funcContainers.Open(ctx)
	assert.Nil(t, err)

	c.fixture = NewDummyAzureFunctionFixture(c.funcContainers.GetHandler())
}

func (c *DummyCommandableAzureFunctionTest) teardown(t *testing.T) {
	err := c.funcContainers.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsCommandableAzure(t *testing.T) {
	c := newDummyCommandableAzureFunctionTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
