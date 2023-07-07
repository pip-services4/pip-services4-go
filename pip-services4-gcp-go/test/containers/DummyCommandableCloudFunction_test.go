package containers_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyCommandableCloudFunctionTest struct {
	fixture        *DummyCloudFunctionFixture
	funcContainers *DummyCloudFunction
}

func newDummyCommandableCloudFunctionTest() *DummyCommandableCloudFunctionTest {
	return &DummyCommandableCloudFunctionTest{}
}

func (c *DummyCommandableCloudFunctionTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
	)

	ctx := context.Background()

	c.funcContainers = NewDummyCloudFunction()
	c.funcContainers.Configure(ctx, config)
	err := c.funcContainers.Open(ctx)
	assert.Nil(t, err)

	c.fixture = NewDummyCloudFunctionFixture(c.funcContainers.GetHandler())
}

func (c *DummyCommandableCloudFunctionTest) teardown(t *testing.T) {
	err := c.funcContainers.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsCommandableCloud(t *testing.T) {
	c := newDummyCommandableCloudFunctionTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
