package fixtures

import (
	"context"
	"testing"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type DummyMapPersistenceFixture struct {
	dummy1      map[string]any
	dummy2      map[string]any
	persistence IDummyMapPersistence
}

func NewDummyMapPersistenceFixture(persistence IDummyMapPersistence) *DummyMapPersistenceFixture {
	c := DummyMapPersistenceFixture{}
	c.dummy1 = map[string]any{"id": "", "key": "Key 11", "content": "Content 1"}
	c.dummy2 = map[string]any{"id": "", "key": "Key 2", "content": "Content 2"}
	c.persistence = persistence
	return &c
}

func (c *DummyMapPersistenceFixture) TestCrudOperations(t *testing.T) {
	var dummy1 map[string]any
	var dummy2 map[string]any

	result, err := c.persistence.Create(context.Background(), c.dummy1)
	assert.Nil(t, err)

	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1["id"])
	assert.Equal(t, c.dummy1["key"], dummy1["key"])
	assert.Equal(t, c.dummy1["content"], dummy1["content"])

	// Create another dummy by set pointer
	result, err = c.persistence.Create(context.Background(), c.dummy2)
	assert.Nil(t, err)

	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2["id"])
	assert.Equal(t, c.dummy2["key"], dummy2["key"])
	assert.Equal(t, c.dummy2["content"], dummy2["content"])

	page, err := c.persistence.GetPageByFilter(context.Background(), *cquery.NewEmptyFilterParams(), *cquery.NewPagingParams(0, 5, true))
	assert.Nil(t, err)

	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
	assert.True(t, page.HasTotal())
	assert.Equal(t, page.Total, 2)

	// Testing default sorting by Key field len
	// Note: may be different order
	assert.True(t, page.Data[0]["key"] == dummy1["key"] || page.Data[0]["key"] == dummy2["key"])
	assert.True(t, page.Data[1]["key"] == dummy1["key"] || page.Data[1]["key"] == dummy2["key"])

	// Update the dummy
	dummy1["content"] = "Updated Content 1"
	result, err = c.persistence.Update(context.Background(), dummy1)
	assert.Nil(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, dummy1["id"], result["id"])
	assert.Equal(t, dummy1["key"], result["key"])
	assert.Equal(t, dummy1["content"], result["content"])

	// Set the dummy (updating)
	dummy1["content"] = "Updated Content 2"
	result, err = c.persistence.Set(context.Background(), dummy1)
	assert.Nil(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, dummy1["id"], result["id"])
	assert.Equal(t, dummy1["key"], result["key"])
	assert.Equal(t, dummy1["content"], result["content"])

	// Set the dummy (creating)
	dummy2["id"] = "New_id"
	dummy2["key"] = "New_key"
	result, err = c.persistence.Set(context.Background(), dummy2)
	assert.Nil(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, dummy2["id"], result["id"])
	assert.Equal(t, dummy2["key"], result["key"])
	assert.Equal(t, dummy2["content"], result["content"])

	// Partially update the dummy
	updateMap := cdata.NewAnyValueMapFromTuples("content", "Partially Updated Content 1")
	result, err = c.persistence.UpdatePartially(context.Background(), dummy1["id"].(string), *updateMap)
	assert.Nil(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, dummy1["id"], result["id"])
	assert.Equal(t, dummy1["key"], result["key"])
	assert.Equal(t, "Partially Updated Content 1", result["content"])

	// Get the dummy by Id
	result, err = c.persistence.GetOneById(context.Background(), dummy1["id"].(string))
	assert.Nil(t, err)

	// Try to get item
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["id"], result["id"])
	assert.Equal(t, dummy1["key"], result["key"])
	assert.Equal(t, "Partially Updated Content 1", result["content"])

	// Delete the dummy
	result, err = c.persistence.DeleteById(context.Background(), dummy1["id"].(string))
	assert.Nil(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, dummy1["id"], result["id"])
	assert.Equal(t, dummy1["key"], result["key"])
	assert.Equal(t, "Partially Updated Content 1", result["content"])

	// Get the deleted dummy
	result, err = c.persistence.GetOneById(context.Background(), dummy1["id"].(string))
	assert.Nil(t, result)
}

func (c *DummyMapPersistenceFixture) TestBatchOperations(t *testing.T) {
	var dummy1 map[string]any
	var dummy2 map[string]any

	// Create one dummy
	result, err := c.persistence.Create(context.Background(), c.dummy1)
	assert.Nil(t, err)

	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1["id"])
	assert.Equal(t, c.dummy1["key"], dummy1["key"])
	assert.Equal(t, c.dummy1["content"], dummy1["content"])

	// Create another dummy
	result, err = c.persistence.Create(context.Background(), c.dummy2)
	assert.Nil(t, err)

	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2["id"])
	assert.Equal(t, c.dummy2["key"], dummy2["key"])
	assert.Equal(t, c.dummy2["content"], dummy2["content"])

	// Read batch
	items, err := c.persistence.GetListByIds(context.Background(), []string{dummy1["id"].(string), dummy2["id"].(string)})
	assert.Nil(t, err)

	assert.NotNil(t, items)
	assert.Len(t, items, 2)

	// Delete batch
	err = c.persistence.DeleteByIds(context.Background(), []string{dummy1["id"].(string), dummy2["id"].(string)})
	assert.Nil(t, err)

	// Read empty batch
	items, err = c.persistence.GetListByIds(context.Background(), []string{dummy1["id"].(string), dummy2["id"].(string)})
	assert.Nil(t, err)

	assert.Len(t, items, 0)
}
