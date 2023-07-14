package test_persistence

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
	c.dummy1 = map[string]any{"Id": "", "Key": "Key 11", "Content": "Content 1"}
	c.dummy2 = map[string]any{"Id": "", "Key": "Key 2", "Content": "Content 2"}
	c.persistence = persistence
	return &c
}

func (c *DummyMapPersistenceFixture) TestCrudOperations(t *testing.T) {
	var dummy1 map[string]any
	var dummy2 map[string]any

	result, err := c.persistence.Create(context.Background(), c.dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1["Id"])
	assert.Equal(t, c.dummy1["Key"], dummy1["Key"])
	assert.Equal(t, c.dummy1["Content"], dummy1["Content"])

	// Create another dummy by set pointer
	result, err = c.persistence.Create(context.Background(), c.dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2["Id"])
	assert.Equal(t, c.dummy2["Key"], dummy2["Key"])
	assert.Equal(t, c.dummy2["Content"], dummy2["Content"])

	page, errp := c.persistence.GetPageByFilter(context.Background(), *cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)
	//Testing default sorting by Key field len

	item1 := page.Data[0]
	assert.Equal(t, item1["Key"], dummy2["Key"])
	item2 := page.Data[1]
	assert.Equal(t, item2["Key"], dummy1["Key"])

	// Update the dummy
	dummy1["Content"] = "Updated Content 1"
	result, err = c.persistence.Update(context.Background(), dummy1)
	if err != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], result["Id"])
	assert.Equal(t, dummy1["Key"], result["Key"])
	assert.Equal(t, dummy1["Content"], result["Content"])

	// Partially update the dummy
	updateMap := *cdata.NewAnyValueMapFromTuples("Content", "Partially Updated Content 1")
	result, err = c.persistence.UpdatePartially(context.Background(), dummy1["Id"].(string), updateMap)
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], result["Id"])
	assert.Equal(t, dummy1["Key"], result["Key"])
	assert.Equal(t, "Partially Updated Content 1", result["Content"])

	// Get the dummy by Id
	result, err = c.persistence.GetOneById(context.Background(), dummy1["Id"].(string))
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], result["Id"])
	assert.Equal(t, dummy1["Key"], result["Key"])
	assert.Equal(t, "Partially Updated Content 1", result["Content"])

	// Delete the dummy
	result, err = c.persistence.DeleteById(context.Background(), dummy1["Id"].(string))
	if err != nil {
		t.Errorf("DeleteById method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], result["Id"])
	assert.Equal(t, dummy1["Key"], result["Key"])
	assert.Equal(t, "Partially Updated Content 1", result["Content"])

	// Get the deleted dummy
	result, err = c.persistence.GetOneById(context.Background(), dummy1["Id"].(string))
	assert.Nil(t, err)
	// Try to get item
	assert.Nil(t, result)
}

func (c *DummyMapPersistenceFixture) TestBatchOperations(t *testing.T) {
	var dummy1 map[string]any
	var dummy2 map[string]any

	// Create one dummy
	result, err := c.persistence.Create(context.Background(), c.dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1["Id"])
	assert.Equal(t, c.dummy1["Key"], dummy1["Key"])
	assert.Equal(t, c.dummy1["Content"], dummy1["Content"])

	// Create another dummy
	result, err = c.persistence.Create(context.Background(), c.dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2["Id"])
	assert.Equal(t, c.dummy2["Key"], dummy2["Key"])
	assert.Equal(t, c.dummy2["Content"], dummy2["Content"])

	// Read batch
	items, err := c.persistence.GetListByIds(context.Background(), []string{dummy1["Id"].(string), dummy2["Id"].(string)})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(t,items)
	assert.NotNil(t, items)
	assert.Len(t, items, 2)

	// Delete batch
	err = c.persistence.DeleteByIds(context.Background(), []string{dummy1["Id"].(string), dummy2["Id"].(string)})
	if err != nil {
		t.Errorf("DeleteByIds method error %v", err)
	}
	assert.Nil(t, err)

	// Read empty batch
	items, err = c.persistence.GetListByIds(context.Background(), []string{dummy1["Id"].(string), dummy2["Id"].(string)})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	assert.NotNil(t, items)
	assert.Len(t, items, 0)

}
