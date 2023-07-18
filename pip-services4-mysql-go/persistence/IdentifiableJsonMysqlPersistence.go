package persistence

import (
	"context"
	"database/sql"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// IdentifiableJsonMySqlPersistence is an abstract persistence component that stores data in MySQL in JSON or JSONB fields
// and implements a number of CRUD operations over data items with unique ids.
// The data items must implement IIdentifiable interface.
//
// The JSON table has only two fields: id and data.
//
// In basic scenarios child classes shall only override getPageByFilter,
// getListByFilter or deleteByFilter operations with specific filter function.
// All other operations can be used out of the box.
//
// In complex scenarios child classes can implement additional operations by
// accessing c._collection and c._model properties.
//
//	Configuration parameters
//
//		- collection:                  (optional) MySQL collection name
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
//	References
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
//
// Example:
//
//	type MyMySqlPersistence struct {
//		*persist.IdentifiableJsonMySqlPersistence[MyData, string]
//	}
//
//	func NewMyMySqlPersistence() *MyMySqlPersistence {
//		c := &MyMySqlPersistence{}
//		c.IdentifiableJsonMySqlPersistence = persist.InheritIdentifiableJsonMySqlPersistence[MyData, string](c, "mydata")
//		return c
//	}
//
//	func (c *MyMySqlPersistence) DefineSchema() {
//		c.ClearSchema()
//		c.EnsureTable("", "")
//		c.EnsureSchema("ALTER TABLE `" + c.TableName + "` ADD `data_key` VARCHAR(50) AS (JSON_UNQUOTE(`data`->\"$.key\"))")
//		c.EnsureIndex(c.TableName+"_json_key", map[string]string{"data_key": "1"}, map[string]string{"unique": "true"})
//	}
//
//	func (c *MyMySqlPersistence) GetPageByFilter(ctx context.Context,
//		filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[MyData], err error) {
//
//		key, ok := filter.GetAsNullableString("Key")
//		filterObj := ""
//		if ok && key != "" {
//			filterObj += "data->'$.key'='" + key + "'"
//		}
//
//		return c.IdentifiableJsonMySqlPersistence.GetPageByFilter(ctx,
//			filterObj, paging,
//			"", "",
//		)
//	}
//
//	func main() {
//		persistence := NewMyMySqlPersistence()
//		persistence.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//			"host", "localhost",
//			"port", 27017,
//		))
//		err := persitence.Open(context.Background())
//
//		item, err := persistence.Create(context.Background(), Mydata{Id: "1", Name: "ABC"})
//		page, err := persistence.GetPageByFilter(context.Background(), *NewFilterParamsFromTuples("name", "ABC"), nil)
//
//		fmt.Println(page.data) // Result: { id: "1", name: "ABC" }
//
//		res, err := persistence.DeleteById(context.Background(), "1")
//	}
type IdentifiableJsonMySqlPersistence[T any, K any] struct {
	*IdentifiableMySqlPersistence[T, K]
}

// InheritIdentifiableJsonMySqlPersistence creates a new instance of the persistence component.
//
//	Parameters:
//		- overrides References to override virtual methods
//		- tableName    (optional) a table name.
func InheritIdentifiableJsonMySqlPersistence[T any, K any](overrides IMySqlPersistenceOverrides[T], tableName string) *IdentifiableJsonMySqlPersistence[T, K] {
	c := &IdentifiableJsonMySqlPersistence[T, K]{}
	c.IdentifiableMySqlPersistence = InheritIdentifiableMySqlPersistence[T, K](overrides, tableName)
	return c
}

// EnsureTable Adds DML statement to automatically create JSON(B) table
//
//		Parameters:
//	  - idType type of the id column (default: VARCHAR(32))
//	  - dataType type of the data column (default: JSON)
func (c *IdentifiableJsonMySqlPersistence[T, K]) EnsureTable(idType string, dataType string) {
	if idType == "" {
		idType = "VARCHAR(32)"
	}
	if dataType == "" {
		dataType = "JSON"
	}

	if c.SchemaName != "" {
		query := "CREATE SCHEMA IF NOT EXISTS " + c.QuoteIdentifier(c.SchemaName)
		c.EnsureSchema(query)
	}
	query := "CREATE TABLE IF NOT EXISTS " + c.QuotedTableName() + " (`id` " + idType + " PRIMARY KEY, `data` " + dataType + ")"
	c.EnsureSchema(query)
}

// ConvertToPublic converts object value from internal to public format.
//
//	Parameters:
//		- value an object in internal format to convert.
//	Returns: converted object in public format.
func (c *IdentifiableJsonMySqlPersistence[T, K]) ConvertToPublic(rows *sql.Rows) (T, error) {
	var defaultValue T
	columns, err := rows.Columns()
	if err != nil {
		return defaultValue, err
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
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

	item, fromJsonErr := c.JsonConvertor.FromJson(mapItem["data"])

	return item, fromJsonErr
}

// ConvertFromPublic convert object value from public to internal format.
//
//		Parameters:
//	   - value     an object in public format to convert.
//
// Returns converted object in internal format.
func (c *IdentifiableJsonMySqlPersistence[T, K]) ConvertFromPublic(value T) (map[string]any, error) {
	id := GetObjectId[K](value)

	data, convErr := c.JsonConvertor.ToJson(value)
	if convErr != nil {
		return nil, convErr
	}

	result := map[string]any{
		"id":   id,
		"data": data,
	}
	return result, nil
}

// ConvertFromPublicPartial convert object value from public to internal format.
//
//	Parameters:
//		- value     an object in public format to convert.
//	Returns: converted object in internal format.
func (c *IdentifiableJsonMySqlPersistence[T, K]) ConvertFromPublicPartial(value map[string]any) (map[string]any, error) {
	buf, toJsonErr := cconv.JsonConverter.ToJson(value)
	if toJsonErr != nil {
		return nil, toJsonErr
	}
	item, fromJsonErr := c.IdentifiableMySqlPersistence.JsonConvertor.FromJson(buf)
	if toJsonErr != nil {
		return nil, fromJsonErr
	}
	return c.ConvertFromPublic(item)
}

// UpdatePartially updates only few selected fields in a data item.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//		- id                an id of data item to be updated.
//		- data              a map with fields to be updated.
//
// Returns: receives updated item or error.
func (c *IdentifiableJsonMySqlPersistence[T, K]) UpdatePartially(ctx context.Context,
	id K, data cdata.AnyValueMap) (result T, err error) {
	buf, toJsonErr := cconv.JsonConverter.ToJson(data.Value())
	if toJsonErr != nil {
		return result, toJsonErr
	}

	query := "UPDATE " + c.QuotedTableName() + " SET `data`=JSON_MERGE_PATCH(data,?) WHERE id=?"
	values := []any{buf, id}

	_, err = c.IdentifiableMySqlPersistence.Client.ExecContext(ctx, query, values...)
	if err != nil {
		return result, err
	}

	// Getting result
	query = "SELECT * FROM " + c.QuotedTableName() + " WHERE id=?"
	rows, err := c.Client.QueryContext(ctx, query, []any{id}...)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	if !rows.Next() {
		return result, rows.Err()
	}

	if err == nil {
		result, convErr := c.IdentifiableMySqlPersistence.Overrides.ConvertToPublic(rows)
		if convErr != nil {
			return result, convErr
		}
		c.IdentifiableMySqlPersistence.Logger.Trace(ctx, "Updated partially in %s with id = %s", c.IdentifiableMySqlPersistence.TableName, id)
		return result, nil
	}
	return result, rows.Err()
}
