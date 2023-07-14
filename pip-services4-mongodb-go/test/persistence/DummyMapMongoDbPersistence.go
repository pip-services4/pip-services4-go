package test_persistence

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-mongodb-go/persistence"
	"go.mongodb.org/mongo-driver/bson"
)

type DummyMapMongoDbPersistence struct {
	*persist.IdentifiableMongoDbPersistence[map[string]any, string]
}

func NewDummyMapMongoDbPersistence() *DummyMapMongoDbPersistence {
	c := &DummyMapMongoDbPersistence{}
	c.IdentifiableMongoDbPersistence = persist.InheritIdentifiableMongoDbPersistence[map[string]any, string](c, "dummies")
	return c
}

func (c *DummyMapMongoDbPersistence) GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[map[string]any], err error) {

	filterObj := bson.M{}

	if key, ok := filter.GetAsNullableString("Key"); ok {
		filterObj = bson.M{"Key": key}
	}

	sorting := bson.M{"Key": -1}

	return c.IdentifiableMongoDbPersistence.GetPageByFilter(ctx, filterObj, paging,
		sorting, nil)
}

func (c *DummyMapMongoDbPersistence) GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error) {

	filterObj := bson.M{}

	if key, ok := filter.GetAsNullableString("Key"); ok {
		filterObj = bson.M{"key": key}
	}
	return c.IdentifiableMongoDbPersistence.GetCountByFilter(ctx, filterObj)
}
