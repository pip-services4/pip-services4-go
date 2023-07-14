package persistence

import (
	"context"
	"math/rand"
	"time"

	"github.com/jinzhu/copier"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-mongodb-go/connect"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"go.mongodb.org/mongo-driver/bson"
	mongodrv "go.mongodb.org/mongo-driver/mongo"
	mongoopt "go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoDbPersistenceOverrides[T any] interface {
	DefineSchema()
	ConvertFromPublic(item T) (map[string]any, error)
	ConvertFromPublicPartial(item T) (map[string]any, error)
	ConvertToPublic(item any) (T, error)
}

// MongoDbPersistence abstract persistence component that stores data in MongoDB using plain driver.
//
// This is the most basic persistence component that is only
// able to store data items of any type. Specific CRUD operations
// over the data items must be implemented in child classes by
// accessing c.Db or c.Collection properties.
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
//	References:
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
//
// Example:
//
//	type MyMongoDbPersistence struct {
//		*persistence.MongoDbPersistence[MyData]
//	}
//
//	func NewMyMongoDbPersistence() *MyMongoDbPersistence {
//		c := &MyMongoDbPersistence{}
//		c.MongoDbPersistence = persistence.InheritMongoDbPersistence[MyData](c, "my_data")
//		return c
//	}
//
//	func (c *MyMongoDbPersistence) GetByName(ctx context.Context, name string) (count int64, err error) {
//		return c.MongoDbPersistence.GetCountByFilter(ctx, bson.M{"name": name})
//	}
//
//	func (c *MyMongoDbPersistence) Set(ctx context.Context,
//		item MyData) (result MyData, err error) {
//		var defaultValue MyData
//
//		newItem, err := c.Overrides.ConvertFromPublic(item)
//		if err != nil {
//			return defaultValue, err
//		}
//
//		id := newItem["_id"]
//		filter := bson.M{"_id": id}
//		var options mngoptions.FindOneAndReplaceOptions
//		retDoc := mngoptions.After
//		options.ReturnDocument = &retDoc
//		upsert := true
//		options.Upsert = &upsert
//
//		res := c.Collection.FindOneAndReplace(ctx, filter, newItem, &options)
//		if err := res.Err(); err != nil {
//			if errors.Is(err, mongo.ErrNoDocuments) {
//				return result, nil
//			}
//			return result, err
//		}
//
//		c.Logger.Trace(ctx, "Set in %s with id = %s", c.CollectionName, id)
//		var docPointer map[string]any
//		if err := res.Decode(&docPointer); err != nil {
//			if errors.Is(err, mongo.ErrNoDocuments) {
//				return result, nil
//			}
//			return result, err
//		}
//
//		return c.Overrides.ConvertToPublic(docPointer)
//	}
//
//	func main() {
//		persistence := NewMyMongoDbPersistence()
//		persistence.Configure(context.Background(), config.NewConfigParamsFromTuples(
//			"host", "localhost",
//			"port", 27017,
//		))
//
//		_ = persistence.Open(context.Background(), "123")
//		persistence.Set(context.Background(), "123", MyData{Id: "123", Name: "ABC"})
//		item, err := persistence.GetByName(context.Background(), "123", "ABC")
//		fmt.Println(item) // Result: { name: "ABC" }
//	}
type MongoDbPersistence[T any] struct {
	Overrides IMongoDbPersistenceOverrides[T]

	defaultConfig   *cconf.ConfigParams
	config          *cconf.ConfigParams
	references      crefer.IReferences
	opened          bool
	localConnection bool
	indexes         []mongodrv.IndexModel
	maxPageSize     int32

	// The dependency resolver.
	DependencyResolver *crefer.DependencyResolver
	// The logger.
	Logger clog.CompositeLogger
	// The MongoDB connection component.
	Connection *conn.MongoDbConnection
	// The MongoDB connection object.
	Client *mongodrv.Client
	// The MongoDB database name.
	DatabaseName string
	// The MongoDB colleciton object.
	CollectionName string
	//  The MongoDb database object.
	Db *mongodrv.Database
	// The MongoDb collection object.
	Collection *mongodrv.Collection

	// Defines channel which closed before closing persistence and signals about terminating
	// all going processes
	//	!IMPORTANT if you do not Close existing query response the persistence can not be closed
	//	see IsTerminated method
	isTerminated chan struct{}

	// Defines general JSON convertors
	JsonConvertor    cconv.IJSONEngine[T]
	JsonMapConvertor cconv.IJSONEngine[map[string]any]
}

// InheritMongoDbPersistence are creates a new instance of the persistence component.
//
//	Parameters:
//		- overrides IMongoDbPersistenceOverrides overrided mongodb persistence
//		- collection  string a collection name.
//
// Returns: *MongoDbPersistence new created MongoDbPersistence component
func InheritMongoDbPersistence[T any](overrides IMongoDbPersistenceOverrides[T], collection string) *MongoDbPersistence[T] {
	c := MongoDbPersistence[T]{
		Overrides:    overrides,
		isTerminated: make(chan struct{}),
	}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"collection", "",
		"dependencies.connection", "*:connection:mongodb:*:1.0",
		"options.max_pool_size", "2",
		"options.keep_alive", "1000",
		"options.connect_timeout", "5000",
		"options.auto_reconnect", "true",
		"options.max_page_size", "100",
		"options.debug", "true",
	)
	c.DependencyResolver = crefer.NewDependencyResolverWithParams(context.Background(), c.defaultConfig, c.references)
	c.Logger = *clog.NewCompositeLogger()
	c.CollectionName = collection
	c.indexes = make([]mongodrv.IndexModel, 0, 10)
	c.config = cconf.NewEmptyConfigParams()
	c.JsonConvertor = cconv.NewDefaultCustomTypeJsonConvertor[T]()
	c.JsonMapConvertor = cconv.NewDefaultCustomTypeJsonConvertor[map[string]any]()

	return &c
}

// Configure method is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams configuration parameters to be set.
func (c *MongoDbPersistence[T]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config
	c.DependencyResolver.Configure(ctx, config)
	c.CollectionName = config.GetAsStringWithDefault("collection", c.CollectionName)
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences references to locate the component dependencies.
func (c *MongoDbPersistence[T]) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.references = references
	c.Logger.SetReferences(ctx, references)

	// try to get a connection
	c.DependencyResolver.SetReferences(ctx, references)
	if conn, ok := c.DependencyResolver.GetOneOptional("connection").(*conn.MongoDbConnection); ok && conn != nil {
		c.Connection = conn
		c.localConnection = false
		return
	}
	// or create a local one
	if c.Connection == nil {
		c.Connection = c.createConnection(ctx)
		c.localConnection = true
	}
}

// UnsetReferences method is unsets (clears) previously set references to dependent components.
func (c *MongoDbPersistence[T]) UnsetReferences() {
	c.Connection = nil
}

func (c *MongoDbPersistence[T]) createConnection(ctx context.Context) *conn.MongoDbConnection {
	connection := conn.NewMongoDbConnection()
	connection.Configure(ctx, c.config)
	if c.references != nil {
		connection.SetReferences(ctx, c.references)
	}
	return connection
}

// DefineSchema for the collection.
// This method shall be overloaded in child classes
func (c *MongoDbPersistence[T]) DefineSchema() {
	// Overload this implementation in child classes
}

// EnsureIndex method are adds index definition to create it on opening
//
//	Parameters:
//		- keys any index keys (fields)
//		- options *mongoopt.IndexOptions index options
func (c *MongoDbPersistence[T]) EnsureIndex(keys any, options *mongoopt.IndexOptions) {
	if keys == nil {
		return
	}
	index := mongodrv.IndexModel{
		Keys:    keys,
		Options: options,
	}
	c.indexes = append(c.indexes, index)
}

// ConvertFromPublic method help convert object (map) from public view by replaced "Id" to "_id" field
//
//	Parameters:
//		- item *any converted item
func (c *MongoDbPersistence[T]) ConvertFromPublic(value T) (map[string]any, error) {
	buf, toBsonErr := bson.Marshal(value)
	if toBsonErr != nil {
		return nil, toBsonErr
	}

	var item map[string]any

	fromBsonErr := bson.Unmarshal(buf, &item)
	if fromBsonErr != nil {
		return nil, fromBsonErr
	}

	if _, ok := any(value).(map[string]any); ok {
		if _, ok := item["Id"]; ok {
			item["_id"] = item["Id"]
			copier.Copy(item["_id"], item["Id"])
			delete(item, "Id")
		}
	}

	return item, nil
}

// ConvertFromPublicPartial method help convert object (map) from public view by replaced "Id" to "_id" field
//
//	Parameters:
//		- item *any converted item
func (c *MongoDbPersistence[T]) ConvertFromPublicPartial(item T) (map[string]any, error) {
	return c.ConvertFromPublic(item)
}

// ConvertToPublic method is convert object (map) to public view by replaced "_id" to "Id" field
//
//	Parameters:
//		- item *any converted item
func (c *MongoDbPersistence[T]) ConvertToPublic(value any) (T, error) {
	var item T

	_, itemOk := any(item).(map[string]any)

	if m, ok := value.(map[string]any); ok && itemOk {
		m["Id"] = m["_id"]
		copier.Copy(m["Id"], m["_id"])
		delete(m, "_id")
	}

	bsonBuf, toBsonErr := bson.Marshal(value)
	if toBsonErr != nil {
		return item, toBsonErr
	}

	fromBsonErr := bson.Unmarshal(bsonBuf, &item)

	return item, fromBsonErr
}

// IsOpen method is checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *MongoDbPersistence[T]) IsOpen() bool {
	return c.opened
}

// IsTerminated checks if the wee need to terminate process before close component.
//
//	Returns: true if you need terminate your processes.
func (c *MongoDbPersistence[T]) IsTerminated() bool {
	select {
	case _, ok := <-c.isTerminated:
		if !ok {
			return true
		}
	default:
		return false
	}
	return false
}

// Open method is opens the component.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occured.
func (c *MongoDbPersistence[T]) Open(ctx context.Context) error {
	if c.opened {
		return nil
	}

	traceId := cctx.GetTraceId(ctx)

	c.isTerminated = make(chan struct{})

	if c.Connection == nil {
		c.Connection = c.createConnection(ctx)
		c.localConnection = true
	}

	c.opened = false
	if c.localConnection {
		err := c.Connection.Open(ctx)
		if err == nil && c.Connection == nil {
			return cerr.NewInvalidStateError(traceId, "NO_CONNECTION", "MongoDB connection is missing")
		}
	}

	if !c.Connection.IsOpen() {
		return cerr.NewConnectionError(traceId, "CONNECT_FAILED", "MongoDB connection is not opened")
	}

	c.Client = c.Connection.GetConnection()
	c.Db = c.Connection.GetDatabase()
	c.DatabaseName = c.Connection.GetDatabaseName()
	if c.Collection = c.Db.Collection(c.CollectionName); c.Collection == nil {
		c.Db = nil
		c.Client = nil
		return cerr.NewConnectionError(traceId, "CONNECT_FAILED", "Connection to mongodb failed")
	}

	// Define database schema
	c.Overrides.DefineSchema()

	// Recreate indexes
	if len(c.indexes) > 0 {
		keys, err := c.Collection.Indexes().CreateMany(ctx, c.indexes, mongoopt.CreateIndexes())
		if err != nil {
			c.Db = nil
			c.Client = nil
			return cerr.NewConnectionError(traceId, "CREATE_IDX_FAILED", "Recreate indexes failed").WithCause(err)
		}
		for _, v := range keys {
			c.Logger.Debug(ctx, traceId, "Created index %s for collection %s", v, c.CollectionName)
		}
	}
	c.opened = true
	c.Logger.Debug(ctx, traceId, "Connected to mongodb database %s, collection %s", c.DatabaseName, c.CollectionName)
	return nil
}

// Close methods closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occured.
func (c *MongoDbPersistence[T]) Close(ctx context.Context) error {
	if !c.opened {
		return nil
	}
	if c.Connection == nil {
		return cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "MongoDb connection is missing")
	}

	defer c.cleanUpConnection()

	if c.localConnection {
		if err := c.Connection.Close(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *MongoDbPersistence[T]) cleanUpConnection() {
	c.opened = false
	c.Client = nil
	c.Db = nil
	c.Collection = nil
	close(c.isTerminated)
}

// Clear method are clears component state.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occurred.
func (c *MongoDbPersistence[T]) Clear(ctx context.Context) error {
	// Return error if collection is not set
	if c.CollectionName == "" {
		return cerr.NewError("Collection name is not defined")
	}

	if err := c.Collection.Drop(ctx); err != nil {
		return cerr.NewConnectionError(cctx.GetTraceId(ctx), "CLEAR_FAILED", "Clear collection failed.").WithCause(err)
	}
	return nil
}

// GetPageByFilter is gets a page of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c *IdentifiableMongoDbPersistence) GetPageByFilter method from child type that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- filter any (optional) a filter JSON object
//		- paging cdata.PagingParams (optional) paging parameters
//		- sort any (optional) sorting BSON object
//		- select  any (optional) projection BSON object
//	Returns: page cdata.DataPage[T], err error a data page or error, if they are occurred
func (c *MongoDbPersistence[T]) GetPageByFilter(ctx context.Context,
	filter any, paging cquery.PagingParams, sort any, sel any) (page cquery.DataPage[T], err error) {
	// Adjust max item count based on configuration

	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c.maxPageSize))
	pagingEnabled := paging.Total
	// Configure options
	var options mongoopt.FindOptions
	if skip >= 0 {
		options.Skip = &skip
	}
	options.Limit = &take
	if sort != nil {
		options.Sort = sort
	}
	if sel != nil {
		options.Projection = sel
	}

	cursor, err := c.Collection.Find(ctx, filter, &options)
	if err != nil {
		return *cquery.NewEmptyDataPage[T](), err
	}
	defer cursor.Close(ctx)

	if c.IsTerminated() {
		return *cquery.NewEmptyDataPage[T](), cerr.
			NewError("query terminated").
			WithTraceId(cctx.GetTraceId(ctx))
	}

	items := make([]T, 0, 1)
	for cursor.Next(ctx) {
		if c.IsTerminated() {
			return *cquery.NewEmptyDataPage[T](), cerr.
				NewError("query terminated").
				WithTraceId(cctx.GetTraceId(ctx))
		}
		var item T

		curErr := cursor.Decode(&item)

		if curErr != nil {
			continue
		}

		items = append(items, item)
	}
	if items != nil {
		c.Logger.Trace(ctx, "Retrieved %d from %s", len(items), c.CollectionName)
	}
	if pagingEnabled {
		if c.IsTerminated() {
			return *cquery.NewEmptyDataPage[T](), cerr.
				NewError("query terminated").
				WithTraceId(cctx.GetTraceId(ctx))
		}
		docCount, _ := c.Collection.CountDocuments(ctx, filter)
		return *cquery.NewDataPage(items, int(docCount)), nil
	}
	return *cquery.NewDataPage(items, cquery.EmptyTotalValue), nil
}

// GetListByFilter is gets a list of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c *IdentifiableMongoDbPersistence) GetListByFilter method from child type that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- filter any (optional) a filter BSON object
//		- sort any (optional) sorting BSON object
//		- select any (optional) projection BSON object
//	Returns: items []any, err error data list and error, if they are occurred
func (c *MongoDbPersistence[T]) GetListByFilter(ctx context.Context,
	filter any, sort any, sel any) (items []T, err error) {

	// Configure options
	var options mongoopt.FindOptions

	if sort != nil {
		options.Sort = sort
	}
	if sel != nil {
		options.Projection = sel
	}

	cursor, err := c.Collection.Find(ctx, filter, &options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items = make([]T, 0)
	if c.IsTerminated() {
		return nil, cerr.
			NewError("query terminated").
			WithTraceId(cctx.GetTraceId(ctx))
	}

	for cursor.Next(ctx) {
		if c.IsTerminated() {
			return nil, cerr.
				NewError("query terminated").
				WithTraceId(cctx.GetTraceId(ctx))
		}
		var docPointer map[string]any
		curErr := cursor.Decode(&docPointer)
		if curErr != nil {
			continue
		}

		item, err := c.Overrides.ConvertToPublic(docPointer)
		if err != nil {
			return items, nil
		}
		items = append(items, item)
	}

	if items != nil {
		c.Logger.Trace(ctx, "Retrieved %d from %s", len(items), c.CollectionName)
	}
	return items, nil
}

// GetOneRandom is gets a random item from items that match to a given filter.
// This method shall be called by a func (c *IdentifiableMongoDbPersistence) getOneRandom method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- filter any (optional) a filter BSON object
//	Returns: item any, err error random item and error, if theq are occured
func (c *MongoDbPersistence[T]) GetOneRandom(ctx context.Context,
	filter any) (item T, err error) {

	docCount, err := c.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return item, err
	}

	var options mongoopt.FindOptions
	rand.Seed(time.Now().UnixNano())
	var itemNum int64 = rand.Int63n(docCount)
	var itemLim int64 = 1

	if itemNum < 0 {
		itemNum = 0
	}
	options.Skip = &itemNum
	options.Limit = &itemLim

	cursor, err := c.Collection.Find(ctx, filter, &options)
	if err != nil {
		return item, err
	}
	defer cursor.Close(ctx)

	var docPointer map[string]any
	cursor.Next(ctx)
	err = cursor.Decode(&docPointer)
	if err != nil {
		return item, err
	}

	return c.Overrides.ConvertToPublic(docPointer)
}

// Create was creates a data item.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- item any an item to be created.
//	Returns: result any, err error created item and error, if they are occurred
func (c *MongoDbPersistence[T]) Create(ctx context.Context, item T) (result T, err error) {
	newItem, err := c.Overrides.ConvertFromPublic(item)
	if err != nil {
		return result, err
	}
	insRes, err := c.Collection.InsertOne(ctx, newItem)
	if err != nil {
		return result, err
	}

	result, err = c.Overrides.ConvertToPublic(newItem)
	if err != nil {
		return result, err
	}
	c.Logger.Trace(ctx, "Created in %s with id = %s", c.Collection, insRes.InsertedID)
	return result, nil
}

// DeleteByFilter is deletes data items that match to a given filter.
// This method shall be called by a func (c *IdentifiableMongoDbPersistence) deleteByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- filter any (optional) a filter BSON object.
//	Returns: error or nil for success.
func (c *MongoDbPersistence[T]) DeleteByFilter(ctx context.Context, filter any) error {
	res, err := c.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	c.Logger.Trace(ctx, "Deleted %d items from %s", res.DeletedCount, c.Collection)
	return nil
}

// GetCountByFilter is gets a count of data items retrieved by a given filter.
// This method shall be called by a func (c *IdentifiableMongoDbPersistence) GetCountByFilter method from child type that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to Trace execution through call chain.
//		- filter any
//	Returns: count int, err error a data count or error, if they are occurred
func (c *MongoDbPersistence[T]) GetCountByFilter(ctx context.Context, filter any) (count int64, err error) {

	// Configure options
	var options mongoopt.CountOptions
	count, err = c.Collection.CountDocuments(ctx, filter, &options)
	if err != nil {
		return 0, err
	}
	c.Logger.Trace(ctx, "Find %d items in %s", count, c.CollectionName)
	return count, nil
}
