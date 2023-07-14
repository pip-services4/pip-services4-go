package test_persistence

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-mongodb-go/persistence"
	"go.mongodb.org/mongo-driver/bson"
)

type DummyRefMongoDbPersistence struct {
	*persist.IdentifiableMongoDbPersistence[*Dummy, string]
}

func NewDummyRefMongoDbPersistence() *DummyRefMongoDbPersistence {
	c := &DummyRefMongoDbPersistence{}
	c.IdentifiableMongoDbPersistence = persist.InheritIdentifiableMongoDbPersistence[*Dummy, string](c, "dummies")
	return c
}

func (c *DummyRefMongoDbPersistence) GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[*Dummy], err error) {

	filterObj := bson.M{}

	if key, ok := filter.GetAsNullableString("Key"); ok {
		filterObj = bson.M{"key": key}
	}

	sorting := bson.M{"key": -1}

	return c.IdentifiableMongoDbPersistence.GetPageByFilter(ctx, filterObj, paging,
		sorting, nil)
}

func (c *DummyRefMongoDbPersistence) GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error) {

	filterObj := bson.M{}

	if key, ok := filter.GetAsNullableString("Key"); ok {
		filterObj = bson.M{"key": key}
	}
	return c.IdentifiableMongoDbPersistence.GetCountByFilter(ctx, filterObj)
}
