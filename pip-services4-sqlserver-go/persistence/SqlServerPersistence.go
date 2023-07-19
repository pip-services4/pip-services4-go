package persistence

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/connect"
)

type ISqlServerPersistenceOverrides[T any] interface {
	DefineSchema()
	ConvertFromPublic(item T) (map[string]any, error)
	ConvertToPublic(item *sql.Rows) (T, error)
	ConvertFromPublicPartial(item map[string]any) (map[string]any, error)
}

// SqlServerPersistence Abstract persistence component that stores data in SqlServer using plain driver.
//
// This is the most basic persistence component that is only
// able to store data items of any type. Specific CRUD operations
// over the data items must be implemented in child classes by
// accessing c._db or c.collection properties.
//
//	Configuration parameters
//		- collection:                  (optional) SqlServer collection name
//		- schema:                  	   (optional) SqlServer schema, default "public"
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
//	import (
//		"context"
//		"fmt"
//
//		cconf "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/config"
//		cpersist "github.com/pip-services4/pip-services4-go/pip-services4-data-go/persistence"
//		persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/persistence"
//		"github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/test/fixtures"
//	)
//
//	type SqlServerServerPersistence struct {
//		*persist.SqlServerPersistence[fixtures.Dummy]
//	}
//
//	func NewSqlServerServerPersistence() *SqlServerServerPersistence {
//		c := &SqlServerServerPersistence{}
//		c.SqlServerPersistence = persist.InheritSqlServerPersistence[fixtures.Dummy](c, "mydata")
//		return c
//	}
//
//	func (c *SqlServerServerPersistence) GetOneByName(ctx context.Context, name string) (item fixtures.Dummy, err error) {
//
//		query := "SELECT * FROM " + c.QuotedTableName() + " WHERE [id]=@p1"
//
//		rows, err := c.Client.QueryContext(ctx, query, name)
//		if err != nil {
//			return item, err
//		}
//		defer rows.Close()
//
//		if !rows.Next() {
//			return item, rows.Err()
//		}
//
//		if err == nil {
//			return c.Overrides.ConvertToPublic(rows)
//		}
//		return item, err
//	}
//
//	func (c *SqlServerServerPersistence) Set(ctx context.Context, item fixtures.Dummy) (result fixtures.Dummy, err error) {
//		objMap, convErr := c.Overrides.ConvertFromPublic(item)
//		if convErr != nil {
//			return result, convErr
//		}
//
//		GenerateObjectMapIdIfNotExists(objMap)
//
//		columns, values := c.GenerateColumnsAndValues(objMap)
//
//		paramsStr := c.GenerateParameters(len(values))
//		columnsStr := c.GenerateColumns(columns)
//		setParams := c.GenerateSetParameters(columns)
//		id := cpersist.GetObjectId(objMap)
//
//		query := "INSERT INTO " + c.QuotedTableName() + " (" + columnsStr + ") OUTPUT INSERTED.* VALUES (" + paramsStr + ")"
//
//		rows, err := c.Client.QueryContext(ctx, query, values...)
//		if err != nil {
//			return result, err
//		}
//		defer rows.Close()
//
//		if rows.Next() {
//			result, convErr = c.Overrides.ConvertToPublic(rows)
//			if convErr != nil {
//				return result, convErr
//			}
//			return result, nil
//		}
//
//		values = append(values, id)
//		query = "UPDATE " + c.QuotedTableName() + " SET " + setParams + " OUTPUT INSERTED.* WHERE [id]=@p" + strconv.FormatInt(int64(len(values)), 10)
//		rows, err = c.Client.QueryContext(ctx, query, values...)
//		if err != nil {
//			return result, err
//		}
//		defer rows.Close()
//
//		if !rows.Next() {
//			return result, rows.Err()
//		}
//
//		result, convErr = c.Overrides.ConvertToPublic(rows)
//		if convErr != nil {
//			return result, convErr
//		}
//
//		return result, rows.Err()
//
//	}
//
//	func main() {
//		persistence := NewSqlServerServerPersistence()
//		persistence.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//			"host", "localhost",
//			"port", 27017,
//		))
//
//		err := persistence.Open(context.Background())
//
//		res, err := persistence.Set(context.Background(), MyData{Id: "1", Name: "ABC"})
//		item, err := persistence.GetOneByName(context.Background(), "ABC")
//		fmt.Println(item) // Result: { Id: "1", Name: "ABC" }
//	}
type SqlServerPersistence[T any] struct {
	Overrides ISqlServerPersistenceOverrides[T]
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
	//The SqlServer connection component.
	Connection *conn.SqlServerConnection
	//The SqlServer connection pool object.
	Client *sql.DB
	//The SqlServer database name.
	DatabaseName string
	//The SqlServer database schema name. If not set use "public" by default
	SchemaName string
	//The SqlServer table object.
	TableName   string
	MaxPageSize int

	// Defines channel which closed before closing persistence and signals about terminating
	// all going processes
	//	!IMPORTANT if you do not Close existing query response the persistence can not be closed
	//	see IsTerminated method
	isTerminated chan struct{}
}

// InheritSqlServerPersistence creates a new instance of the persistence component.
//
//	Parameters:
//		- overrides References to override virtual methods
//		- tableName    (optional) a table name.
func InheritSqlServerPersistence[T any](overrides ISqlServerPersistenceOverrides[T], tableName string) *SqlServerPersistence[T] {
	c := &SqlServerPersistence[T]{
		Overrides: overrides,
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"collection", nil,
			"dependencies.connection", "*:connection:sqlserver:*:1.0",
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
		JsonConvertor:    cconv.NewDefaultCustomTypeJsonConvertor[T](),
		JsonMapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[map[string]any](),
		isTerminated:     make(chan struct{}),
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
func (c *SqlServerPersistence[T]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config

	c.DependencyResolver.Configure(ctx, config)

	c.TableName = config.GetAsStringWithDefault("collection", c.TableName)
	c.TableName = config.GetAsStringWithDefault("table", c.TableName)
	c.MaxPageSize = config.GetAsIntegerWithDefault("options.max_page_size", c.MaxPageSize)
	c.SchemaName = config.GetAsStringWithDefault("schema", c.SchemaName)
}

// SetReferences to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references references to locate the component dependencies.
func (c *SqlServerPersistence[T]) SetReferences(ctx context.Context, references cref.IReferences) {

	c.references = references
	c.Logger.SetReferences(ctx, references)

	// Get connection
	c.DependencyResolver.SetReferences(ctx, references)
	result := c.DependencyResolver.GetOneOptional("connection")

	if dep, ok := result.(*conn.SqlServerConnection); ok {
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
func (c *SqlServerPersistence[T]) UnsetReferences() {
	c.Connection = nil
}

func (c *SqlServerPersistence[T]) createConnection(ctx context.Context) *conn.SqlServerConnection {
	connection := conn.NewSqlServerConnection()
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
func (c *SqlServerPersistence[T]) EnsureIndex(name string, keys map[string]string, options map[string]string) {
	builder := "CREATE"
	if options == nil {
		options = make(map[string]string, 0)
	}

	if options["unique"] != "" {
		builder += " UNIQUE"
	}

	indexName := c.QuoteIdentifier(name)

	if c.SchemaName != "" {
		indexName = c.QuoteIdentifier(c.SchemaName) + "." + indexName
	}

	builder += " INDEX " + indexName + " ON " + c.QuotedTableName()

	if options["type"] != "" {
		builder += " " + options["type"]
	}

	fields := ""
	for key, _ := range keys {
		if fields != "" {
			fields += ", "
		}
		fields += c.QuoteIdentifier(key)
		asc := keys[key]
		if asc != "1" {
			fields += " DESC"
		}
	}

	builder += "(" + fields + ")"

	c.EnsureSchema(builder)
}

// DefineSchema a database schema for this persistence, have to call in child class
// Override in child classes
func (c *SqlServerPersistence[T]) DefineSchema() {
	c.ClearSchema()
}

// EnsureSchema adds a statement to schema definition
//
//		Parameters:
//	  - schemaStatement a statement to be added to the schema
func (c *SqlServerPersistence[T]) EnsureSchema(schemaStatement string) {
	c.schemaStatements = append(c.schemaStatements, schemaStatement)
}

// ClearSchema clears all auto-created objects
func (c *SqlServerPersistence[T]) ClearSchema() {
	c.schemaStatements = []string{}
}

// ConvertToPublic converts object value from internal to func (c * SqlServerPersistence) format.
//
//	Parameters:
//		- value an object in internal format to convert.
//	Returns: converted object in func (c * SqlServerPersistence) format.
func (c *SqlServerPersistence[T]) ConvertToPublic(rows *sql.Rows) (T, error) {
	var defaultValue T
	columns, err := rows.Columns()
	if err != nil {
		return defaultValue, err
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]any' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]any, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// result map
	mapItem := make(map[string]string, len(columns))

	// get RawBytes from data
	err = rows.Scan(scanArgs...)
	if err != nil {
		return defaultValue, err
	}

	for i := 0; i < len(columns); i++ {
		// Here we can check if the value is nil (NULL value)
		mapItem[columns[i]] = string(values[i])
	}

	if err = rows.Err(); err != nil {
		return defaultValue, err
	}

	jsonBuf, toJsonErr := cconv.JsonConverter.ToJson(mapItem)
	if toJsonErr != nil {
		return defaultValue, toJsonErr
	}

	item, fromJsonErr := c.JsonConvertor.FromJson(jsonBuf)

	return item, fromJsonErr

}

// ConvertFromPublic сonvert object value from func (c * SqlServerPersistence) to internal format.
//
//	Parameters:
//		- value an object in func (c * SqlServerPersistence) format to convert.
//	Returns: converted object in internal format.
func (c *SqlServerPersistence[T]) ConvertFromPublic(value T) (map[string]any, error) {
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
func (c *SqlServerPersistence[T]) ConvertFromPublicPartial(value map[string]any) (map[string]any, error) {
	buf, toJsonErr := cconv.JsonConverter.ToJson(value)
	if toJsonErr != nil {
		return nil, toJsonErr
	}

	item, fromJsonErr := c.JsonMapConvertor.FromJson(buf)
	return item, fromJsonErr
}

func (c *SqlServerPersistence[T]) QuoteIdentifier(value string) string {
	if value == "" {
		return value
	}
	if value[0] == '[' {
		return value
	}
	return "[" + strings.ReplaceAll(value, ".", "].[") + "]"
}

// QuotedTableName return quoted SchemaName with TableName ("schema"."table")
func (c *SqlServerPersistence[T]) QuotedTableName() string {
	if c.TableName == "" {
		return ""
	}

	builder := c.QuoteIdentifier(c.TableName)
	if c.SchemaName != "" {
		builder = c.QuoteIdentifier(c.SchemaName) + "." + builder
	}

	return builder
}

// IsOpen checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *SqlServerPersistence[T]) IsOpen() bool {
	return c.opened
}

// IsTerminated checks if the wee need to terminate process before close component.
//
//	Returns: true if you need terminate your processes.
func (c *SqlServerPersistence[T]) IsTerminated() bool {
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
func (c *SqlServerPersistence[T]) Open(ctx context.Context) (err error) {
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
		err = cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "SqlServer connection is missing")
	}

	if err == nil && !c.Connection.IsOpen() {
		err = cerr.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "SqlServer connection is not opened")
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
		err = cerr.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to sqlserver failed").WithCause(err)
	} else {
		c.opened = true
		c.Logger.Debug(ctx, "Connected to sqlserver database %s, collection %s", c.DatabaseName, c.QuotedTableName())
	}

	return err
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *SqlServerPersistence[T]) Close(ctx context.Context) (err error) {
	if !c.opened {
		return nil
	}

	if c.Connection == nil {
		return cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "SqlServer connection is missing")
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
func (c *SqlServerPersistence[T]) Clear(ctx context.Context) error {
	// Return error if collection is not set
	if c.TableName == "" {
		return errors.New("Table name is not defined")
	}

	rows, err := c.Client.QueryContext(ctx, "DELETE FROM "+c.QuotedTableName())
	if err != nil {
		return cerr.
			NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to sqlserver failed").
			WithCause(err)
	}
	rows.Close()
	return nil
}

func (c *SqlServerPersistence[T]) CreateSchema(ctx context.Context) (err error) {
	if len(c.schemaStatements) == 0 {
		return nil
	}

	// Check if table exist to determine weither to auto create objects
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
	return nil
}

func (c *SqlServerPersistence[T]) checkTableExists(ctx context.Context) (bool, error) {
	// Check if table exist to determine either to auto create objects
	query := "SELECT OBJECT_ID('" + c.TableName + "', 'U') as oid" // TODO check
	result, err := c.Client.QueryContext(ctx, query)
	if err != nil {
		return false, err
	}

	defer result.Close()

	values := make([]sql.RawBytes, 1)
	// rows.Scan wants '[]any' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]any, 1)
	scanArgs[0] = &values[0]

	// If table already exists then exit
	if result.Next() {

		// get RawBytes from data
		err = result.Scan(scanArgs...)
		if err != nil {
			return false, err
		}

		for _, col := range values {
			if col != nil {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
	return false, nil
}

// GenerateColumns generates a list of column names to use in SQL statements like: "column1,column2,column3"
//
//	Parameters:
//		- columns an array with column values
//	Returns: a generated list of column names
func (c *SqlServerPersistence[T]) GenerateColumns(columns []string) string {
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

// GenerateParameters generates a list of value parameters to use in SQL statements like: "?,?,?"
//
//	Parameters:
//		- values an array with column values or a key-value map
//	Returns: a generated list of value parameters
func (c *SqlServerPersistence[T]) GenerateParameters(valuesCount int) string {
	if valuesCount <= 0 {
		return ""
	}

	builder := strings.Builder{}
	for index := 1; index <= valuesCount; index++ {
		if builder.String() != "" {
			builder.WriteString(",")
		}
		builder.WriteString("@p" + strconv.Itoa(index))
	}

	return builder.String()
}

// GenerateSetParameters generates a list of column sets to use in UPDATE statements like: column1=?,column2=?
//
//	Parameters:
//		- values an array with column values or a key-value map
//	Returns: a generated list of column sets
func (c *SqlServerPersistence[T]) GenerateSetParameters(columns []string) string {

	if len(columns) == 0 {
		return ""
	}
	setParamsBuf := strings.Builder{}
	for i := range columns {
		if setParamsBuf.String() != "" {
			setParamsBuf.WriteString(",")
		}
		setParamsBuf.WriteString(c.QuoteIdentifier(columns[i]) + "=@p" + strconv.Itoa(i+1))
	}
	return setParamsBuf.String()
}

// GenerateColumnsAndValues generates a list of column parameters
//
//	Parameters:
//		- values an array with column values or a key-value map
//	Returns: a generated list of column values
func (c *SqlServerPersistence[T]) GenerateColumnsAndValues(objMap map[string]any) ([]string, []any) {
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
// This method shall be called by a func (c * SqlServerPersistence) getPageByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object
//		- paging            (optional) paging parameters
//		- sort              (optional) sorting JSON object
//		- select            (optional) projection JSON object
//	Returns: receives a data page or error.
func (c *SqlServerPersistence[T]) GetPageByFilter(ctx context.Context,
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
	} else {
		query += " ORDER BY 1"
	}

	if skip < 0 {
		skip = 0
	}

	query += " OFFSET " + strconv.FormatInt(skip, 10) + " ROWS FETCH NEXT " + strconv.FormatInt(take, 10) + " ROWS ONLY"

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

		return *cquery.NewDataPage[T](items, int(count)), nil
	}

	return *cquery.NewDataPage[T](items, cquery.EmptyTotalValue), rows.Err()
}

// GetCountByFilter gets a number of data items retrieved by a given filter.
// This method shall be called by a func (c * SqlServerPersistence) getCountByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object
//	Returns: data page or error.
func (c *SqlServerPersistence[T]) GetCountByFilter(ctx context.Context,
	filter string) (int64, error) {

	query := "SELECT COUNT(*) AS count FROM " + c.QuotedTableName()
	if len(filter) > 0 {
		query += " WHERE " + filter
	}

	rows, err := c.Client.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	values := make([]sql.RawBytes, 1)
	scanArgs := make([]any, 1)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return 0, err
		}

		count = cconv.LongConverter.ToLong(string(values[0]))
	}

	if count != 0 {
		c.Logger.Trace(ctx, "Counted %d items in %s", count, c.TableName)
	}

	return count, rows.Err()
}

// GetListByFilter gets a list of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c * SqlServerPersistence) getListByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter           (optional) a filter JSON object
//		- paging           (optional) paging parameters
//		- sort             (optional) sorting JSON object
//		- select           (optional) projection JSON object
//	Returns: data list or error.
func (c *SqlServerPersistence[T]) GetListByFilter(ctx context.Context,
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
// This method shall be called by a func (c * SqlServerPersistence) getOneRandom method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object
//	Returns: random item or error.
func (c *SqlServerPersistence[T]) GetOneRandom(ctx context.Context, filter string) (item T, err error) {
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
	query += " ORDER BY (SELECT NULL) OFFSET " + strconv.FormatInt(pos, 10) + " ROWS FETCH NEXT 1 ROWS ONLY"

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
func (c *SqlServerPersistence[T]) Create(ctx context.Context, item T) (result T, err error) {
	objMap, convErr := c.Overrides.ConvertFromPublic(item)
	if convErr != nil {
		return result, convErr
	}

	columns, values := c.GenerateColumnsAndValues(objMap)

	columnsStr := c.GenerateColumns(columns)
	paramsStr := c.GenerateParameters(len(values))

	query := "INSERT INTO " + c.QuotedTableName() + " (" + columnsStr + ") OUTPUT INSERTED.* VALUES (" + paramsStr + ")"

	rows, err := c.Client.QueryContext(ctx, query, values...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	id := GetObjectId[any](item)
	c.Logger.Trace(ctx, "Created in %s with id = %s", c.TableName, id)
	return item, nil
}

// DeleteByFilter deletes data items that match to a given filter.
// This method shall be called by a func (c * SqlServerPersistence) deleteByFilter method from child class that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- filter            (optional) a filter JSON object.
//	Returns: error or nil for success.
func (c *SqlServerPersistence[T]) DeleteByFilter(ctx context.Context, filter string) error {
	query := "DELETE FROM " + c.QuotedTableName()
	if len(filter) > 0 {
		query += " WHERE " + filter
	}

	result, err := c.Client.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	c.Logger.Trace(ctx, "Deleted %d items from %s", count, c.TableName)
	return nil
}

func (c *SqlServerPersistence[T]) cloneItem(item any) T {
	if cloneableItem, ok := item.(cdata.ICloneable[T]); ok {
		return cloneableItem.Clone()
	}

	strObject, _ := c.JsonConvertor.ToJson(item.(T))
	newItem, _ := c.JsonConvertor.FromJson(strObject)
	return newItem
}
