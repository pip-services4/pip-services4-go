package containers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/data"
	"github.com/stretchr/testify/assert"
)

type DummyCloudFunctionFixture struct {
	DUMMY1  tdata.Dummy
	DUMMY2  tdata.Dummy
	handler http.HandlerFunc
}

func NewDummyCloudFunctionFixture(handler http.HandlerFunc) *DummyCloudFunctionFixture {
	return &DummyCloudFunctionFixture{
		handler: handler,
		DUMMY1:  *tdata.NewDummy("", "key 1", "content 1"),
		DUMMY2:  *tdata.NewDummy("", "key 2", "content 2"),
	}
}

func (c *DummyCloudFunctionFixture) TestCrudOperations(t *testing.T) {
	var dummy1 tdata.Dummy
	var dummy2 tdata.Dummy

	// Create one dummy
	res := c.invokeFunc(t, map[string]any{
		"cmd":   "create_dummy",
		"dummy": c.DUMMY1,
	})

	err := json.Unmarshal(res, &dummy1)
	assert.Nil(t, err)

	assert.Equal(t, dummy1.Key, c.DUMMY1.Key)
	assert.Equal(t, dummy1.Content, c.DUMMY1.Content)

	// Create another dummy
	res = c.invokeFunc(t, map[string]any{
		"cmd":   "create_dummy",
		"dummy": c.DUMMY2,
	})

	err = json.Unmarshal(res, &dummy2)
	assert.Nil(t, err)

	assert.Equal(t, dummy2.Key, c.DUMMY2.Key)
	assert.Equal(t, dummy2.Content, c.DUMMY2.Content)

	// Update the dummy
	dummy1.Content = "Updated Content 1"

	res = c.invokeFunc(t, map[string]any{
		"cmd":   "update_dummy",
		"dummy": dummy1,
	})

	var updatedDummy1 tdata.Dummy
	err = json.Unmarshal(res, &updatedDummy1)
	assert.Nil(t, err)

	assert.Equal(t, updatedDummy1.Id, dummy1.Id)
	assert.Equal(t, updatedDummy1.Content, dummy1.Content)
	assert.Equal(t, updatedDummy1.Key, dummy1.Key)
	dummy1 = updatedDummy1

	// Delete dummy
	res = c.invokeFunc(t, map[string]any{
		"cmd":      "delete_dummy",
		"dummy_id": dummy1.Id,
	})

	var deleted tdata.Dummy
	err = json.Unmarshal(res, &deleted)
	assert.Nil(t, err)

	assert.Equal(t, deleted.Id, dummy1.Id)
	assert.Equal(t, deleted.Content, dummy1.Content)
	assert.Equal(t, deleted.Key, dummy1.Key)

	res = c.invokeFunc(t, map[string]any{
		"cmd":      "get_dummy_by_id",
		"dummy_id": dummy1.Id,
	})

	var empty tdata.Dummy
	err = json.Unmarshal(res, &empty)
	assert.Nil(t, err)

	assert.Equal(t, empty.Id, "")
	assert.Equal(t, empty.Content, "")
	assert.Equal(t, empty.Key, "")

	// Failed validation
	res = c.invokeFunc(t, map[string]any{
		"cmd":   "create_dummy",
		"dummy": nil,
	})

	var validErr cerr.ApplicationError
	err = json.Unmarshal(res, &validErr)
	assert.Nil(t, err)

	assert.Equal(t, validErr.Code, "INVALID_DATA")
}

func (c *DummyCloudFunctionFixture) invokeFunc(t *testing.T, data any) []byte {
	body, err := cconv.JsonConverter.ToJson(data)
	assert.Nil(t, err)

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	c.handler(rr, req)

	return rr.Body.Bytes()
}
