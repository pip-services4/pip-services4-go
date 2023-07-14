package persistence

import (
	"context"
	"errors"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mngoptions "go.mongodb.org/mongo-driver/mongo/options"
)

// IdentifiableMongoDbPersistence is abstract persistence component that stores data in MongoDB
// and implements a number of CRUD operations over data items with unique ids.
// The data items must implement IIdentifiable interface.
//
// In basic scenarios child classes shall only override GetPageByFilter,
// GetListByFilter or DeleteByFilter operations with specific filter function.
// All other operations can be used out of the box.
//
// In complex scenarios child classes can implement additional operations by
// accessing c.Collection properties.
//
//	Configuration parameters:
//		- collection:                  (optional) MongoDB collection name
//		- connection(s):
//			- discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//			- host:                      host name or IP address
//			- port:                      port number (default: 27017)
//			- database:                  database name
//			- uri:                       resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                 (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                  (optional) user name
//			- password:                  (optional) user password
//		- options:
//			- max_pool_size:             (optional) maximum connection pool size (default: 2)
//			- keep_alive:                (optional) enable connection keep alive (default: true)
//			- connect_timeout:           (optional) connection timeout in milliseconds (default: 5000)
//			- socket_timeout:            (optional) socket timeout in milliseconds (default: 360000)
//			- auto_reconnect:            (optional) enable auto reconnection (default: true) (not used)
//			- reconnect_interval:        (optional) reconnection interval in milliseconds (default: 1000) (not used)
//			- max_page_size:             (optional) maximum page size (default: 100)
//			- replica_set:               (optional) name of replica set
//			- ssl:                       (optional) enable SSL connection (default: false) (not implements in this release)
//			- auth_source:               (optional) authentication source
//			- debug:                     (optional) enable debug output (default: false). (not used)
//
//	References:
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
//
// Example:
//
//	type MyIdentifiableMongoDbPersistence struct {
//		*persist.IdentifiableMongoDbPersistence[test_persistence.Dummy, string]
//	}
//
//	func NewMyIdentifiableMongoDbPersistence() *MyIdentifiableMongoDbPersistence {
//		c := &MyIdentifiableMongoDbPersistence{}
//		c.IdentifiableMongoDbPersistence = persist.InheritIdentifiableMongoDbPersistence[test_persistence.Dummy, string](c, "dummies")
//		return c
//	}
//
//	func composeFilter(filter cdata.FilterParams) any {
//		filterObj := bson.M{}
//
//		if name, ok := filter.GetAsNullableString("name"); ok {
//			filterObj = bson.M{"name": name}
//		}
//
//		return filterObj
//	}
//
//	func (c *MyIdentifiableMongoDbPersistence) GetPageByFilter(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[test_persistence.Dummy], err error) {
//		return c.IdentifiableMongoDbPersistence.GetPageByFilter(ctx, composeFilter(filter), paging,
//			bson.M{"key": -1}, nil)
//	}
//
//	func main() {
//		persistence := NewMyIdentifiableMongoDbPersistence()
//		persistence.Configure(context.Background(), config.NewConfigParamsFromTuples(
//			"host", "localhost",
//			"port", 27017,
//		))
//
//		_ = persistence.Open(context.Background(), "123")
//		page, err := persistence.GetPageByFilter(context.Background(), *cdata.NewFilterParamsFromTuples("name", "ABC"), *cdata.NewEmptyPagingParams())
//		fmt.Println(page) // Result: { id: "1", name: "ABC" }
//
//		err = persistence.DeleteByFilter(context.Background(), "1")
//	}
type IdentifiableMongoDbPersistence[T any, K any] struct {
	*MongoDbPersistence[T]

	// Flag to turn on automated string ID generation
	_autoGenerateId bool
}

// InheritIdentifiableMongoDbPersistence is creates a new instance of the persistence component.
//
//	Parameters:
//		- collection string (optional) a collection name.
//	Returns: *IdentifiableMongoDbPersistence[T, K] new created IdentifiableMongoDbPersistence component
func InheritIdentifiableMongoDbPersistence[T any, K any](overrides IMongoDbPersistenceOverrides[T], collection string) *IdentifiableMongoDbPersistence[T, K] {
	if collection == "" {
		panic("Collection name could not be nil")
	}
	c := IdentifiableMongoDbPersistence[T, K]{}
	c.MongoDbPersistence = InheritMongoDbPersistence(overrides, collection)
	c.maxPageSize = 100
	c._autoGenerateId = true
	return &c
}

// Configure is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams configuration parameters to be set.
func (c *IdentifiableMongoDbPersistence[T, K]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.MongoDbPersistence.Configure(ctx, config)
	c.maxPageSize = (int32)(config.GetAsIntegerWithDefault("options.max_page_size", (int)(c.maxPageSize)))
}

// GetListByIds is gets a list of data items retrieved by given unique ids.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- ids  []K ids of data items to be retrieved
//	Returns: items []T, err error a data list and error, if they are occurred.
func (c *IdentifiableMongoDbPersistence[T, K]) GetListByIds(ctx context.Context,
	ids []K) (items []T, err error) {

	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}
	return c.GetListByFilter(ctx, filter, nil, nil)
}

// GetOneById is gets a data item by its unique id.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- id                an id of data item to be retrieved.
//	Returns: item T, err error a data and error, if they are occurred.
func (c *IdentifiableMongoDbPersistence[T, K]) GetOneById(ctx context.Context,
	id K) (item T, err error) {

	filter := bson.M{"_id": id}
	var docPointer map[string]any

	res := c.Collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return item, nil
		}
		return item, err
	}

	if err := res.Decode(&docPointer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return item, nil
		}
		return item, err
	}
	c.Logger.Trace(ctx, "Retrieved from %s by id = %s", c.CollectionName, id)
	return c.Overrides.ConvertToPublic(docPointer)
}

// Create was creates a data item.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- item any an item to be created.
//	Returns: result any, err error created item and error, if they are occurred
func (c *IdentifiableMongoDbPersistence[T, K]) Create(ctx context.Context,
	item T) (result T, err error) {
	var defaultValue T

	newItem, err := c.Overrides.ConvertFromPublic(item)
	if err != nil {
		return defaultValue, err
	}

	// Auto generate unique id
	val, ok := newItem["_id"]
	if (!ok || val == nil || val == "") && c._autoGenerateId {
		newItem["_id"] = keys.IdGenerator.NextLong()
	}

	res, err := c.Collection.InsertOne(ctx, newItem)
	if err != nil {
		return result, err
	}

	result, err = c.Overrides.ConvertToPublic(newItem)
	if err != nil {
		return defaultValue, err
	}

	c.Logger.Trace(ctx, "Created in %s with id = %s", c.Collection, res.InsertedID)

	return result, nil
}

// Set is sets a data item. If the data item exists it updates it,
// otherwise it create a new data item.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- item T an item to be set.
//	Returns: result any, err error updated item and error, if they occurred
func (c *IdentifiableMongoDbPersistence[T, K]) Set(ctx context.Context,
	item T) (result T, err error) {
	var defaultValue T

	newItem, err := c.Overrides.ConvertFromPublic(item)
	if err != nil {
		return defaultValue, err
	}

	// Auto unique generate id
	val, ok := newItem["_id"]
	if (!ok || val == nil || val == "") && c._autoGenerateId {
		newItem["_id"] = keys.IdGenerator.NextLong()
	}

	id := newItem["_id"]
	filter := bson.M{"_id": id}
	var options mngoptions.FindOneAndReplaceOptions
	retDoc := mngoptions.After
	options.ReturnDocument = &retDoc
	upsert := true
	options.Upsert = &upsert

	res := c.Collection.FindOneAndReplace(ctx, filter, newItem, &options)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
		return result, err
	}

	c.Logger.Trace(ctx, "Set in %s with id = %s", c.CollectionName, id)
	var docPointer map[string]any
	if err := res.Decode(&docPointer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
		return result, err
	}

	return c.Overrides.ConvertToPublic(docPointer)
}

// Update is updates a data item.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- item T an item to be updated.
//	Returns: result any, err error updated item and error, if they are occurred
func (c *IdentifiableMongoDbPersistence[T, K]) Update(ctx context.Context,
	item T) (result T, err error) {

	newItem, err := c.Overrides.ConvertFromPublic(item)
	if err != nil {
		return result, err
	}
	id := newItem["_id"]

	filter := bson.M{"_id": id}
	update := bson.D{{"$set", newItem}}

	var options mngoptions.FindOneAndUpdateOptions
	retDoc := mngoptions.After
	options.ReturnDocument = &retDoc

	res := c.Collection.FindOneAndUpdate(ctx, filter, update, &options)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
		return result, err
	}

	c.Logger.Trace(ctx, "Updated in %s with id = %s", c.CollectionName, id)

	var docPointer map[string]any
	if err := res.Decode(&docPointer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
		return result, err
	}

	return c.Overrides.ConvertToPublic(docPointer)
}

// UpdatePartially is updates only few selected fields in a data item.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- id K an id of data item to be updated.
//		- data cdata.AnyValueMap a map with fields to be updated.
//	Returns: item any, err error updated item and error, if they are occurred
func (c *IdentifiableMongoDbPersistence[T, K]) UpdatePartially(ctx context.Context,
	id K, data cdata.AnyValueMap) (item T, err error) {

	newItem := bson.M{}
	for k, v := range data.Value() {
		newItem[k] = v
	}
	filter := bson.M{"_id": id}
	update := bson.D{{"$set", newItem}}

	var options mngoptions.FindOneAndUpdateOptions
	retDoc := mngoptions.After
	options.ReturnDocument = &retDoc

	res := c.Collection.FindOneAndUpdate(ctx, filter, update, &options)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return item, nil
		}
		return item, err
	}
	c.Logger.Trace(ctx, "Updated partially in %s with id = %s", c.Collection, id)

	var docPointer map[string]any
	if err := res.Decode(&docPointer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return item, nil
		}
		return item, err
	}

	return c.Overrides.ConvertToPublic(docPointer)
}

// DeleteById is deleted a data item by it's unique id.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- id K id of the item to be deleted
//	Returns: item T, err error deleted item and error, if they are occurred
func (c *IdentifiableMongoDbPersistence[T, K]) DeleteById(ctx context.Context,
	id K) (item T, err error) {

	filter := bson.M{"_id": id}

	res := c.Collection.FindOneAndDelete(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return item, nil
		}
		return item, err
	}

	c.Logger.Trace(ctx, "Deleted from %s with id = %s", c.CollectionName, id)

	var docPointer map[string]any
	if err := res.Decode(&docPointer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return item, nil
		}
		return item, err
	}

	return c.Overrides.ConvertToPublic(docPointer)
}

// DeleteByIds is deletes multiple data items by their unique ids.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- ids []K ids of data items to be deleted.
//	Returns: error or nil for success.
func (c *IdentifiableMongoDbPersistence[T, K]) DeleteByIds(ctx context.Context,
	ids []K) error {

	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}
	return c.DeleteByFilter(ctx, filter)
}
