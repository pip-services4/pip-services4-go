package containers_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

type DummyCloudFunctionTest struct {
	fixture        *DummyCloudFunctionFixture
	funcContainers *DummyCloudFunction
}

func newDummyCloudFunctionTest() *DummyCloudFunctionTest {
	return &DummyCloudFunctionTest{}
}

func (c *DummyCloudFunctionTest) setup(t *testing.T) {
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

func (c *DummyCloudFunctionTest) teardown(t *testing.T) {
	err := c.funcContainers.Close(context.Background())
	assert.Nil(t, err)
}

func TestCrudOperationsCloud(t *testing.T) {
	c := newDummyCloudFunctionTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
