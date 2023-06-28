package test_clients

import (
	"context"
	"testing"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/test/sample"
	"github.com/stretchr/testify/assert"
)

type DummyClientFixture struct {
	client IDummyClient
}

func NewDummyClientFixture(client IDummyClient) *DummyClientFixture {
	dcf := DummyClientFixture{client: client}
	return &dcf
}

func (c *DummyClientFixture) TestCrudOperations(t *testing.T) {
	dummy1 := tsample.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	dummy2 := tsample.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	ctx := cctx.NewContextWithTraceId(context.Background(), "ClientFixture")

	// Create one dummy
	dummy, err := c.client.CreateDummy(ctx, dummy1)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, dummy1.Content)
	assert.Equal(t, dummy.Key, dummy1.Key)
	dummy1 = dummy

	// Create another dummy
	dummy, err = c.client.CreateDummy(ctx, dummy2)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, dummy2.Content)
	assert.Equal(t, dummy.Key, dummy2.Key)
	dummy2 = dummy

	// Get all dummies
	dummies, err := c.client.GetDummies(ctx, *cquery.NewEmptyFilterParams(), *cquery.NewPagingParams(0, 5, false))
	assert.Nil(t, err)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	dummy, err = c.client.UpdateDummy(ctx, dummy1)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, dummy1.Key)
	dummy1 = dummy

	// Delete dummy
	dummy, err = c.client.DeleteDummy(ctx, dummy1.Id)
	assert.Nil(t, err)

	// Try to get delete dummy
	dummy, err = c.client.GetDummyById(ctx, dummy1.Id)
	assert.Nil(t, err)
	assert.Equal(t, tsample.Dummy{}, dummy)

	// Check trace id propagation
	ctx = cctx.NewContextWithTraceId(context.Background(), "test_trace_id")

	values, err := c.client.CheckTraceId(ctx)
	assert.Nil(t, err)
	assert.Equal(t, values["traceId"], "test_trace_id")

	values, err = c.client.CheckTraceId(ctx)
	assert.Nil(t, err)
	assert.Equal(t, values["traceId"], "test_trace_id")

	// Check error propagation
	ctx = cctx.NewContextWithTraceId(context.Background(), "test_error_propagation")

	err = c.client.CheckErrorPropagation(ctx)
	appErr, ok := err.(*cerr.ApplicationError)

	assert.True(t, ok)
	assert.Equal(t, appErr.TraceId, "test_error_propagation")
	assert.Equal(t, appErr.Status, 404)
	assert.Equal(t, appErr.Code, "NOT_FOUND_TEST")
	assert.Equal(t, appErr.Message, "Not found error")
}
