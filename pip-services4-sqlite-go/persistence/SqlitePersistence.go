package persistence

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/connect"
)

type ISqlitePersistenceOverrides[T any] interface {
	DefineSchema()
	ConvertFromPublic(item T) (map[string]any, error)
	ConvertToPublic(item *sql.Rows) (T, error)
	ConvertFromPublicPartial(item map[string]any) (map[string]any, error)
}

// SqlitePersistence Abstract persistence component that stores data in SQLite using plain driver.
//
// This is the most basic persistence component that is only
// able to store data items of any type. Specific CRUD operations
// over the data items must be implemented in child classes by
// accessing c._db or c._collection properties.
//
//	Configuration parameters
//		- collection:                  (optional) SQLite collection name
//		- schema:                  	   (optional) SQLite schema, default "public"
//		- connection(s):
//			- discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//			- host:                      host name or IP address
//			- port:                      port number (default: 27017)
//			- uri:                       resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                 (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                  (optional) user name
//			- password:                  (optional) user password
//		- options:
//			- connect_timeout:      (optional) number of milliseconds to wait before timing out when connecting a new client (default: 0)
//			- idle_timeout:         (optional) number of milliseconds a client must sit idle in the pool and not be checked out (default: 10000)
//			- max_pool_size:        (optional) maximum number of clients the pool should contain (default: 10)
//
//	References:
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
//
// Example:
//
//	type MySqlitePersistence struct {
//		*persist.SqlitePersistence[MyData]
//	}
//
//	func NewMySqlitePersistence() *MySqlitePersistence {
//		c := &MySqlitePersistence{}
//		c.SqlitePersistence = persist.InheritSqlitePersistence[MyData](c, "mydata")
//		return c
//	}
//
//	func (c *MySqlitePersistence) DefineSchema() {
//		c.ClearSchema()
//		c.SqlitePersistence.DefineSchema()
//		// Row name must be in double quotes for properly case!!!
//		c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (\"id\" TEXT PRIMARY KEY, \"key\" TEXT, \"content\" TEXT)")
//		c.EnsureIndex(c.SqlitePersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
//	}
//
//	func (c *MySqlitePersistence) GetOneById(ctx context.Context, name string) (item MyData, err error) {
//		query := "SELECT * FROM " + c.QuotedTableName() + " WHERE \"name\"=$1"
//
//		qResult, err := c.Client.QueryContext(ctx, query, name)
//		if err != nil {
//			return item, err
//		}
//		defer qResult.Close()
//
//		if !qResult.Next() {
//			return item, qResult.Err()
//		}
//
//		result, err := c.Overrides.ConvertToPublic(qResult)
//
//		if err != nil {
//			c.Logger.Trace(ctx, "Nothing found from %s with name = %s", c.TableName, name)
//			return item, err
//		}
//		c.Logger.Trace(ctx, "Retrieved from %s with name = %s", c.TableName, name)
//		return result, err
//
//	}
//
//	func (c *MySqlitePersistence) Set(ctx context.Context, item MyData) (result MyData, err error) {
//		objMap, convErr := c.Overrides.ConvertFromPublic(item)
//		if convErr != nil {
//			return result, convErr
//		}
//
//		columns, values := c.GenerateColumnsAndValues(objMap)
//
//		paramsStr := c.GenerateParameters(len(values))
//		columnsStr := c.GenerateColumns(columns)
//		setParams := c.GenerateSetParameters(columns)
//
//		id := objMap["id"]
//
//		query := "INSERT INTO " + c.QuotedTableName() + " (" + columnsStr + ")" +
//			" VALUES (" + paramsStr + ")" +
//			" ON CONFLICT (\"id\") DO UPDATE SET " + setParams + " RETURNING *"
//
//		qResult, err := c.Client.QueryContext(ctx, query, values...)
//		if err != nil {
//			return result, err
//		}
//		defer qResult.Close()
//
//		if !qResult.Next() {
//			return result, qResult.Err()
//		}
//
//		result, convErr = c.Overrides.ConvertToPublic(qResult)
//		if convErr != nil {
//			return result, convErr
//		} else {
//			c.Logger.Trace(ctx, "Set in %s with id = %s", c.TableName, id)
//			return result, nil
//		}
//	}
type SqlitePersistence[T any] struct {
	Overrides ISqlitePersistenceOverrides[T]
	// Defines general JSON convertors
	JsonConvertor    cconv.IJSONEngine[T]
	JsonMapConvertor cconv.IJSONEngine[map[string]any]

	defaultConfig *cconf.ConfigParams

	config           *cconf.ConfigParams
	references       cref.IReferences
	opened           bool
	localConnection  bool
	schemaStatements []string

	//The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	//The logger.
	Logger *clog.CompositeLogger
	//The SQLite connection component.
	Connection *conn.SqliteConnection
	//The SQLite database.
	Client *sql.DB
	//The SQLite database name.
	DatabaseName string
	//The SQLite database schema name. If not set use "public" by default
	SchemaName string
	//The SQLite table object.
	TableName   string
	MaxPageSize int

	// Defines channel which closed before closing persistence and signals about terminating
	// all going processes
	//	!IMPORTANT if you do not Close existing query response the persistence can not be closed
	//	see IsTerminated method
	isTerminated chan struct{}
}

// InheritSqlitePersistence creates a new instance of the persistence component.
//
//	Parameters:
//		- ctx context.Context
//		- overrides References to override virtual methods
//		- tableName    (optional) a table name.
func InheritSqlitePersistence[T any](overrides ISqlitePersistenceOverrides[T], tableName string) *SqlitePersistence[T] {
	c := &SqlitePersistence[T]{
		Overrides: overrides,
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"collection", nil,
			"dependencies.connection", "*:connection:sqlite:*:1.0",
			"options.max_pool_size", 2,
			"options.keep_alive", 1,
			"options.connect_timeout", 5000,
			"options.auto_reconnect", true,
			"options.max_page_size", 100,
			"options.debug", true,
		),
		schemaStatements: make([]string, 0),
		Logger:           clog.NewCompositeLogger(),
		MaxPageSize:      100,
		TableName:        tableName,
		isTerminated:     make(chan struct{}),
		JsonConvertor:    cconv.NewDefaultCustomTypeJsonConvertor[T](),
		JsonMapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[map[string]any](),
	}

	c.DependencyResolver = cref.NewDependencyResolver()
	c.DependencyResolver.Configure(context.Background(), c.defaultConfig)

	return c
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *SqlitePersistence[T]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config

	c.DependencyResolver.Configure(ctx, config)

	c.TableName = config.GetAsStringWithDefault("collection", c.TableName)
	c.TableName = config.GetAsStringWithDefault("table", c.TableName)
	c.MaxPageSize = config.GetAsIntegerWithDefault("options.max_page_size", c.MaxPageSize)
}

// SetReferences to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references references to locate the component dependencies.
func (c *SqlitePersistence[T]) SetReferences(ctx context.Context, references cref.IReferences) {

	c.references = references
	c.Logger.SetReferences(ctx, references)

	// Get connection
	c.DependencyResolver.SetReferences(ctx, references)
	result := c.DependencyResolver.GetOneOptional("connection")

	if dep, ok := result.(*conn.SqliteConnection); ok {
		c.Connection = dep
	}
	// Or create a local one
	if c.Connection == nil {
		c.Connection = c.createConnection(ctx)
		c.localConnection = true
	} else {
		c.localConnection = false
	}
}

// UnsetReferences (clears) previously set references to dependent components.
func (c *SqlitePersistence[T]) UnsetReferences() {
	c.Connection = nil
}

func (c *SqlitePersistence[T]) createConnection(ctx context.Context) *conn.SqliteConnection {
	connection := conn.NewSqliteConnection()
	if c.config != nil {
		connection.Configure(ctx, c.config)
	}
	if c.references != nil {
		connection.SetReferences(ctx, c.references)
	}
	return connection
}

// EnsureIndex adds index definition to create it on opening
//
//	Parameters:
//		- keys index keys (fields)
//		- options index options
func (c *SqlitePersistence[T]) EnsureIndex(name string, keys map[string]string, options map[string]string) {
	builder := "CREATE"
	if options == nil {
		options = make(map[string]string, 0)
	}

	if options["unique"] != "" {
		builder += " UNIQUE"
	}

	indexName := c.QuoteIdentifier(name)

	builder += " INDEX IF NOT EXISTS " + indexName + " ON " + c.QuotedTableName()

	if options["type"] != "" {
		builder += " " + options["type"]
	}

	fields := ""
	for key := range keys {
		if fields != "" {
			fields += ", "
		}
		fields += key
		asc := keys[key]
		if asc != "1" {
			fields += " DESC"
		}
	}

	builder += "(" + fields + ")"

	c.EnsureSchema(builder)
}

// DefineSchema a database schema for this persistence, have to call in child class
func (c *SqlitePersistence[T]) DefineSchema() {
	// Override in child classes

	if len(c.SchemaName) > 0 {
		c.EnsureSchema("CREATE SCHEMA IF NOT EXISTS " + c.QuoteIdentifier(c.SchemaName))
	}
}

// EnsureSchema adds a statement to schema definition
//
//		Parameters:
//	  - schemaStatement a statement to be added to the schema
func (c *SqlitePersistence[T]) EnsureSchema(schemaStatement string) {
	c.schemaStatements = append(c.schemaStatements, schemaStatement)
}

// ClearSchema clears all auto-created objects
func (c *SqlitePersistence[T]) ClearSchema() {
	c.schemaStatements = []string{}
}

// ConvertToPublic converts object value from internal to func (c * SqlitePersistence) format.
//
//	Parameters:
//		- value an object in internal format to convert.
//	Returns: converted object in func (c * SqlitePersistence) format.
func (c *SqlitePersistence[T]) ConvertToPublic(rows *sql.Rows) (T, error) {
	var defaultValue T
	columns, err := rows.Columns()
	if err != nil || columns == nil || len(columns) == 0 {
		return defaultValue, err
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}

	err = rows.Scan(pointers...)
	if err != nil {
		return defaultValue, err
	}

	buf := make(map[string]interface{}, 0)

	for index, column := range columns {
		buf[column] = values[index]
	}

	jsonBuf, toJsonErr := cconv.JsonConverter.ToJson(buf)
	if toJsonErr != nil {
		return defaultValue, toJsonErr
	}

	item, fromJsonErr := c.JsonConvertor.FromJson(jsonBuf)

	return item, fromJsonErr

}

// ConvertFromPublic сonvert object value from func (c * SqlitePersistence) to internal format.
//
//	Parameters:
//		- value an object in func (c * SqlitePersistence) format to convert.
//	Returns: converted object in internal format.
func (c *SqlitePersistence[T]) ConvertFromPublic(value T) (map[string]any, error) {
	buf, toJsonErr := cconv.JsonConverter.ToJson(value)
	if toJsonErr != nil {
		return nil, toJsonErr
	}

	item, fromJsonErr := c.JsonMapConvertor.FromJson(buf)

	return item, fromJsonErr
}

// ConvertFromPublicPartial converts the given object from the public partial format.
//
//	Parameters:
//		- value the object to convert from the public partial format.
//	Returns: the initial object.
func (c *SqlitePersistence[T]) ConvertFromPublicPartial(value map[string]any) (map[string]any, error) {
	buf, toJsonErr := cconv.JsonConverter.ToJson(value)
	if toJsonErr != nil {
		return nil, toJsonErr
	}

	item, fromJsonErr := c.JsonMapConvertor.FromJson(buf)
	return item, fromJsonErr
}

func (c *SqlitePersistence[T]) QuoteIdentifier(value string) string {
	if value == "" {
		return value
	}
	if value[0] == '\'' {
		return value
	}
	return "\"" + value + "\""
}

// QuotedTableName return quoted SchemaName with TableName ("schema"."table")
func (c *SqlitePersistence[T]) QuotedTableName() string {
	if len(c.SchemaName) > 0 {
		return c.QuoteIdentifier(c.SchemaName) + "." + c.QuoteIdentifier(c.TableName)
	}
	return c.QuoteIdentifier(c.TableName)
}

// IsOpen checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *SqlitePersistence[T]) IsOpen() bool {
	return c.opened
}

// IsTerminated checks if the wee need to terminate process before close component.
//
//	Returns: true if you need terminate your processes.
func (c *SqlitePersistence[T]) IsTerminated() bool {
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

// Open the component.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *SqlitePersistence[T]) Open(ctx context.Context) (err error) {
	if c.opened {
		return nil
	}

	c.isTerminated = make(chan struct{})

	if c.Connection == nil {
		c.Connection = c.createConnection(ctx)
		c.localConnection = true
	}

	if c.localConnection {
		err = c.Connection.Open(ctx)
	}

	if err == nil && c.Connection == nil {
		err = cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "SQLite connection is missing")
	}

	if err == nil && !c.Connection.IsOpen() {
		err = cerr.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "SQLite connection is not opened")
	}

	c.opened = false

	if err != nil {
		return err
	}
	c.Client = c.Connection.GetConnection()
	c.DatabaseName = c.Connection.GetDatabaseName()

	// Define database schema
	c.Overrides.DefineSchema()

	// Recreate objects
	err = c.CreateSchema(ctx)
	if err != nil {
		c.Client = nil
		err = cerr.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to sqlite failed"+err.Error()).WithCause(err)
	} else {
		c.opened = true
		c.Logger.Debug(ctx, "Connected to sqlite database %s, table %s", c.DatabaseName, c.QuotedTableName())
	}

	return err
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *SqlitePersistence[T]) Close(ctx context.Context) (err error) {
	if !c.opened {
		return nil
	}

	if c.Connection == nil {
		return cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Sqlite connection is missing")
	}

	close(c.isTerminated)
	if c.localConnection {
		err = c.Connection.Close(ctx)
	}
	if err != nil {
		return err
	}
	c.opened = false
	c.Client = nil
	c.Connection = nil
	c.isTerminated = nil
	return nil
}

// Clear component state.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil no errors occured.
func (c *SqlitePersistence[T]) Clear(ctx context.Context) error {
	// Return error if collection is not set
	if c.TableName == "" {
		return errors.New("TABLE NAME IS NOT DEFINED")
	}

	_, err := c.Client.ExecContext(ctx, "DELETE FROM "+c.QuotedTableName())
	if err != nil {
		return cerr.
			NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to sqlite failed: "+err.Error()).
			WithCause(err)
	}
	return nil
}

func (c *SqlitePersistence[T]) CreateSchema(ctx context.Context) (err error) {
	if c.schemaStatements == nil || len(c.schemaStatements) == 0 {
		return nil
	}

	exists, err := c.checkTableExists(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	c.Logger.Debug(ctx, "Table "+c.QuotedTableName()+" does not exist. Creating database objects...")

	for _, dml := range c.schemaStatements {
		_, err := c.Client.ExecContext(ctx, dml)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to autocreate database object")
			return err
		}
	}
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for _, dml := range c.schemaStatements {
	// 		_, err := c.Client.ExecContext(ctx, dml)
	// 		if err != nil {
	// 			c.Logger.Error(ctx, err, "Failed to autocreate database object")
	// 		}
	// 	}
	// }()
	// wg.Wait()
	return nil
}

func (c *SqlitePersistence[T]) checkTableExists(ctx context.Context) (bool, error) {
	// Check if table exist to determine either to auto create objects
	query := "SELECT * FROM '" + c.TableName + "' LIMIT 1"
	_, qErr := c.Client.ExecContext(ctx, query)
	if qErr != nil {
		if !strings.Contains(qErr.Error(), "no such table") {
			return false, qErr
		}
		return false, nil
	}

	return true, nil
}

// GenerateColumns generates a list of column names to use in SQL statements like: "column1,column2,column3"
//
//	Parameters:
//		- columns an array with column values
//	Returns: a generated list of column names
func (c *SqlitePersistence[T]) GenerateColumns(columns []string) string {
	if len(columns) == 0 {
		return ""
	}

	builder := strings.Builder{}
	for _, item := range columns {
		if builder.String() != "" {
			builder.WriteString(",")
		}
		builder.WriteString(c.QuoteIdentifier(item))
	}
	return builder.String()
}

// GenerateParameters generates a list of value parameters to use in SQL statements like: "$1,$2,$3"
//
//	Parameters:
//		- values an array with column values or a key-value map
//	Returns: a generated list of value parameters
func (c *SqlitePersistence[T]) GenerateParameters(valuesCount int) string {
	if valuesCount <= 0 {
		return ""
	}

	builder := strings.Builder{}
	for index := 1; index <= valuesCount; index++ {
		if builder.String() != "" {
			builder.WriteString(",")
		}
		builder.WriteString("$")
		builder.WriteString(strconv.FormatInt((int64)(index), 10))
	}

	return builder.String()
}

// GenerateSetParameters generates a list of column sets to use in UPDATE statements like: column1=$1,column2=$2
//
//	Parameters:
//		- values an array with column values or a key-value map
//	Returns: a generated list of column sets
func (c *SqlitePersistence[T]) GenerateSetParameters(columns []string) string {

	if len(columns) == 0 {
		return ""
	}
	setParamsBuf := strings.Builder{}
	index := 1
	for i := range columns {
		if setParamsBuf.String() != "" {
			setParamsBuf.WriteString(",")
		}
		setParamsBuf.WriteString(c.QuoteIdentifier(columns[i]) + "=$" + strconv.FormatInt((int64)(index), 10))
		index++
	}
	return setParamsBuf.String()
}

// GenerateColumnsAndValues generates a list of column parameters
//
//	Parameters:
//		- values an array with column values or a key-value map
//	Returns: a generated list of column values
func (c *SqlitePersistence[T]) GenerateColumnsAndValues(objMap map[string]any) ([]string, []any) {
	if len(objMap) == 0 {
		return nil, nil
	}

	ln := len(objMap)
	columns := make([]string, 0, ln)
	values := make([]any, 0, ln)
	for _col, _val := range objMap {
		columns = append(columns, _col)
		values = append(values, _val)
	}
	return columns, values
}

// GetPageByFilter gets a page of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c * SqlitePersistence) getPageByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object
//		- paging            (optional) paging parameters
//		- sort              (optional) sorting JSON object
//		- select            (optional) projection JSON object
//	Returns: receives a data page or error.
func (c *SqlitePersistence[T]) GetPageByFilter(ctx context.Context,
	filter string, paging cquery.PagingParams, sort string, selection string) (page cquery.DataPage[T], err error) {

	query := "SELECT * FROM " + c.QuotedTableName()
	if len(selection) > 0 {
		query = "SELECT " + selection + " FROM " + c.QuotedTableName()
	}

	// Adjust max item count based on configuration paging
	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c.MaxPageSize))
	pagingEnabled := paging.Total

	if len(filter) > 0 {
		query += " WHERE " + filter
	}
	if len(sort) > 0 {
		query += " ORDER BY " + sort
	}

	query += " LIMIT " + strconv.FormatInt(take, 10)

	if skip >= 0 {
		query += " OFFSET " + strconv.FormatInt(skip, 10)
	}

	rows, err := c.Client.QueryContext(ctx, query)
	if err != nil {
		return *cquery.NewEmptyDataPage[T](), err
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		if c.IsTerminated() {
			rows.Close()
			return *cquery.NewEmptyDataPage[T](), cerr.
				NewError("query terminated").
				WithTraceId(cctx.GetTraceId(ctx))
		}
		item, convErr := c.Overrides.ConvertToPublic(rows)
		if convErr != nil {
			return page, convErr
		}
		items = append(items, item)
	}

	if items != nil {
		c.Logger.Trace(ctx, "Retrieved %d from %s", len(items), c.TableName)
	}

	if pagingEnabled {
		count, err := c.GetCountByFilter(ctx, filter)
		if err != nil {
			return *cquery.NewEmptyDataPage[T](), err
		}

		return *cquery.NewDataPage(items, int(count)), nil
	}

	return *cquery.NewDataPage(items, cquery.EmptyTotalValue), rows.Err()
}

// GetCountByFilter gets a number of data items retrieved by a given filter.
// This method shall be called by a func (c * SqlitePersistence) getCountByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object
//	Returns: data page or error.
func (c *SqlitePersistence[T]) GetCountByFilter(ctx context.Context,
	filter string) (int64, error) {

	query := "SELECT COUNT(*) AS count FROM " + c.QuotedTableName()
	if len(filter) > 0 {
		query += " WHERE " + filter
	}

	queryRes, err := c.Client.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer queryRes.Close()

	var count int64

	if queryRes.Next() {
		var cnt interface{}
		err := queryRes.Scan(&cnt)
		if err != nil {
			cnt = 0
		}
		count = cconv.LongConverter.ToLong(cnt)
	}

	if count != 0 {
		c.Logger.Trace(ctx, "Counted %d items in %s", count, c.TableName)
	}

	return count, queryRes.Err()
}

// GetListByFilter gets a list of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c * SqlitePersistence) getListByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter           (optional) a filter JSON object
//		- paging           (optional) paging parameters
//		- sort             (optional) sorting JSON object
//		- select           (optional) projection JSON object
//	Returns: data list or error.
func (c *SqlitePersistence[T]) GetListByFilter(ctx context.Context,
	filter string, sort string, selection string) (items []T, err error) {

	query := "SELECT * FROM " + c.QuotedTableName()

	if len(selection) > 0 {
		query = "SELECT " + selection + " FROM " + c.QuotedTableName()
	}

	if len(filter) > 0 {
		query += " WHERE " + filter
	}

	if len(sort) > 0 {
		query += " ORDER BY " + sort
	}

	rows, err := c.Client.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items = make([]T, 0, 1)
	for rows.Next() {
		if c.IsTerminated() {
			rows.Close()
			return nil, cerr.
				NewError("query terminated").
				WithTraceId(cctx.GetTraceId(ctx))
		}
		item, convErr := c.Overrides.ConvertToPublic(rows)
		if convErr != nil {
			return items, convErr
		}
		items = append(items, item)
	}

	if items != nil {
		c.Logger.Trace(ctx, "Retrieved %d from %s", len(items), c.TableName)
	}

	return items, rows.Err()
}

// GetOneRandom gets a random item from items that match to a given filter.
// This method shall be called by a func (c * SqlitePersistence) getOneRandom method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object
//	Returns: random item or error.
func (c *SqlitePersistence[T]) GetOneRandom(ctx context.Context, filter string) (item T, err error) {
	count, err := c.GetCountByFilter(ctx, filter)
	if err != nil {
		return item, err
	}
	if count == 0 {
		c.Logger.Trace(ctx, "Can't retriev random item from %s. Table is empty.", c.TableName)
		return item, nil
	}
	if c.IsTerminated() {
		return item, cerr.
			NewError("query terminated").
			WithTraceId(cctx.GetTraceId(ctx))
	}

	rand.Seed(time.Now().UnixNano())
	pos := rand.Int63n(int64(count))

	// build query
	query := "SELECT * FROM " + c.QuotedTableName()
	if len(filter) > 0 {
		query += " WHERE " + filter
	}
	query += " LIMIT 1" + " OFFSET " + strconv.FormatInt(pos, 10)

	rows, err := c.Client.QueryContext(ctx, query)
	if err != nil {
		return item, err
	}
	defer rows.Close()

	if !rows.Next() {
		c.Logger.Trace(ctx, "Random item wasn't found from %s", c.TableName)
		return item, rows.Err()
	}

	item, convErr := c.Overrides.ConvertToPublic(rows)
	if convErr != nil {
		return item, convErr
	}
	c.Logger.Trace(ctx, "Retrieved random item from %s", c.TableName)
	return item, nil

}

// Create creates a data item.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- item              an item to be created.
//	Returns: (optional) callback function that receives created item or error.
func (c *SqlitePersistence[T]) Create(ctx context.Context, item T) (result T, err error) {
	objMap, convErr := c.Overrides.ConvertFromPublic(item)
	if convErr != nil {
		return result, convErr
	}
	columns, values := c.GenerateColumnsAndValues(objMap)

	columnsStr := c.GenerateColumns(columns)
	paramsStr := c.GenerateParameters(len(values))

	query := "INSERT INTO " + c.QuotedTableName() +
		" (" + columnsStr + ") VALUES (" + paramsStr + ") RETURNING *"

	rows, err := c.Client.QueryContext(ctx, query, values...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	if !rows.Next() {
		return result, rows.Err()
	}

	result, convErr = c.Overrides.ConvertToPublic(rows)
	if convErr != nil {
		return result, convErr
	}
	id := GetObjectId[any](result)
	c.Logger.Trace(ctx, "Created in %s with id = %s", c.TableName, id)
	return result, nil
}

// DeleteByFilter deletes data items that match to a given filter.
// This method shall be called by a func (c * SqlitePersistence) deleteByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object.
//	Returns: error or nil for success.
func (c *SqlitePersistence[T]) DeleteByFilter(ctx context.Context, filter string) error {
	query := "DELETE FROM " + c.QuotedTableName()
	if len(filter) > 0 {
		query += " WHERE " + filter
	}

	qResult, qErr := c.Client.ExecContext(ctx, query)
	if qErr != nil {
		return qErr
	}

	count, err := qResult.RowsAffected()
	if err != nil {
		return err
	}

	c.Logger.Trace(ctx, "Deleted %d items from %s", count, c.TableName)
	return nil
}

func (c *SqlitePersistence[T]) cloneItem(item any) T {
	if cloneableItem, ok := item.(cdata.ICloneable[T]); ok {
		return cloneableItem.Clone()
	}

	strObject, _ := c.JsonConvertor.ToJson(item.(T))
	newItem, _ := c.JsonConvertor.FromJson(strObject)
	return newItem
}
