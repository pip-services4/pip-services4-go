package test_persistence

import (
	"context"
	"testing"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type DummyPersistenceFixture struct {
	dummy1      Dummy
	dummy2      Dummy
	persistence IDummyPersistence
}

func NewDummyPersistenceFixture(persistence IDummyPersistence) *DummyPersistenceFixture {
	c := DummyPersistenceFixture{}
	c.dummy1 = Dummy{Id: "", Key: "Key 11", Content: "Content 1"}
	c.dummy2 = Dummy{Id: "", Key: "Key 2", Content: "Content 2"}
	c.persistence = persistence
	return &c
}

func (c *DummyPersistenceFixture) TestCrudOperations(t *testing.T) {
	var dummy1 Dummy
	var dummy2 Dummy

	result, err := c.persistence.Create(context.Background(), c.dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1.Id)
	assert.NotEqual(t, dummy1.Id, "")
	assert.Equal(t, c.dummy1.Key, dummy1.Key)
	assert.Equal(t, c.dummy1.Content, dummy1.Content)

	// Create another dummy by send pointer
	result, err = c.persistence.Create(context.Background(), c.dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2.Id)
	assert.NotEqual(t, dummy2.Id, "")
	assert.Equal(t, c.dummy2.Key, dummy2.Key)
	assert.Equal(t, c.dummy2.Content, dummy2.Content)

	page, errp := c.persistence.GetPageByFilter(context.Background(),
		*cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.True(t, page.HasData())
	assert.Len(t, page.Data, 2)
	//Testing default sorting by Key field len

	assert.Equal(t, page.Data[0].Key, dummy2.Key)
	assert.Equal(t, page.Data[1].Key, dummy1.Key)

	page, errp = c.persistence.GetPageByFilter(context.Background(),
		*cquery.NewEmptyFilterParams(), *cquery.NewPagingParams(10, 1, false))
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.False(t, page.HasData())
	assert.Len(t, page.Data, 0)

	// Get count
	count, errc := c.persistence.GetCountByFilter(context.Background(), *cquery.NewEmptyFilterParams())
	assert.Nil(t, errc)
	assert.Equal(t, count, int64(2))

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	result, err = c.persistence.Update(context.Background(), dummy1)
	if err != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, dummy1.Content, result.Content)

	// Partially update the dummy
	result, err = c.persistence.UpdatePartially(context.Background(), dummy1.Id, *cdata.NewAnyValueMapFromTuples("Content", "Partially Updated Content 1"))
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, "Partially Updated Content 1", result.Content)

	// Get the dummy by Id
	result, err = c.persistence.GetOneById(context.Background(), dummy1.Id)
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, "Partially Updated Content 1", result.Content)

	// Delete the dummy
	result, err = c.persistence.DeleteById(context.Background(), dummy1.Id)
	if err != nil {
		t.Errorf("DeleteById method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, "Partially Updated Content 1", result.Content)

	// Get the deleted dummy
	result, err = c.persistence.GetOneById(context.Background(), dummy1.Id)
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item, must be an empty Dummy struct
	var temp Dummy
	assert.Equal(t, temp, result)
}

func (c *DummyPersistenceFixture) TestBatchOperations(t *testing.T) {
	var dummy1 Dummy
	var dummy2 Dummy

	// Create one dummy
	result, err := c.persistence.Create(context.Background(), c.dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1.Id)
	assert.Equal(t, c.dummy1.Key, dummy1.Key)
	assert.Equal(t, c.dummy1.Content, dummy1.Content)

	// Create another dummy
	result, err = c.persistence.Create(context.Background(), c.dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2.Id)
	assert.Equal(t, c.dummy2.Key, dummy2.Key)
	assert.Equal(t, c.dummy2.Content, dummy2.Content)

	// Read batch
	items, err := c.persistence.GetListByIds(context.Background(), []string{dummy1.Id, dummy2.Id})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(t,items)
	assert.NotNil(t, items)
	assert.Len(t, items, 2)

	// Delete batch
	err = c.persistence.DeleteByIds(context.Background(), []string{dummy1.Id, dummy2.Id})
	if err != nil {
		t.Errorf("DeleteByIds method error %v", err)
	}
	assert.Nil(t, err)

	// Read empty batch
	items, err = c.persistence.GetListByIds(context.Background(), []string{dummy1.Id, dummy2.Id})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	assert.Nil(t, items)
	assert.Len(t, items, 0)

}
