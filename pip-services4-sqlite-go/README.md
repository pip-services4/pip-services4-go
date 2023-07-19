# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> SQLite components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit. It provides a set of components to implement SQLite persistence.

Client was based on [SQLite go driver](https://github.com/mattn/go-sqlite3)
[Official docs](https://pkg.go.dev/github.com/mattn/go-sqlite3) for SQLite Go driver

The module contains the following packages:
- **Build** -  Factory to create SQLite persistence components.
- **Connect** - Connection component to configure SQLite connection to database.
- **Persistence** - abstract persistence components to perform basic CRUD operations.

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/conceptual/configuration/component_configuration/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)

## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go@latest
```


As an example, lets create persistence for the following data object.

```go
type MyData struct {
	Id      string `bson:"_id" json:"id"`
	Key     string `bson:"key" json:"key"`
	Content string `bson:"content" json:"content"`
}

func (d *MyData) SetId(id string) {
	d.Id = id
}

func (d MyData) GetId() string {
	return d.Id
}

func (d MyData) Clone() MyData {
	return MyData{
		Id:      d.Id,
		Key:     d.Key,
		Content: d.Content,
	}
}
```

The persistence component shall implement the following interface with a basic set of CRUD operations.

```go
type IMyPersistence interface {
	GetPageByFilter(ctx context.Context, filter data.FilterParams, paging data.PagingParams) (page data.DataPage[MyData], err error)
	GetOneById(ctx context.Context, id string) (item MyData, err error)
	GetOneByKey(ctx context.Context, key string) (item MyData, err error)
	Create(ctx context.Context, item MyData) (result MyData, err error)
	Update(ctx context.Context, item MyData) (result MyData, err error)
	DeleteById(ctx context.Context, id string) (item MyData, err error)
}
```

To implement postgresql persistence component you shall inherit `IdentifiableSqlitePersistence`. 
Most CRUD operations will come from the base class. You only need to override `GetPageByFilter` method with a custom filter function.
And implement a `GetOneByKey` custom persistence method that doesn't exist in the base class.

```go
type MySqlitePersistence struct {
	*persistence.IdentifiableSqlitePersistence[MyData, string]
}

func NewMySqlitePersistence() *MySqlitePersistence {
	c := &MySqlitePersistence{}
	c.IdentifiableSqlitePersistence = persistence.InheritIdentifiableSqlitePersistence[MyData, string](c, "myobjects")
	return c
}

func (c *MySqlitePersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableSqlitePersistence.DefineSchema()
	// Row name must be in double quotes for properly case!!!
	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (\"id\" TEXT PRIMARY KEY, \"key\" TEXT, \"content\" TEXT)")
	c.EnsureIndex(c.IdentifiableSqlitePersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *MySqlitePersistence) composeFilter(filter data.FilterParams) string {
	criteria := make([]string, 0)

	if key, ok := filter.GetAsNullableString("key"); ok && key != "" {
		criteria = append(criteria, "key='"+key+"'")
	}

	if id, ok := filter.GetAsNullableString("id"); ok && id != "" {
		criteria = append(criteria, "id='"+id+"'")
	}

	if tempIds, ok := filter.GetAsNullableString("ids"); ok && tempIds != "" {
		ids := strings.Split(tempIds, ",")
		criteria = append(criteria, "id IN ('"+strings.Join(ids, "','")+"')")
	}

	if len(criteria) > 0 {
		return strings.Join(criteria, " AND ")
	} else {
		return ""
	}
}

func (c *MySqlitePersistence) GetPageByFilter(ctx context.Context,
	filter data.FilterParams, paging data.PagingParams) (page data.DataPage[MyData], err error) {

	return c.IdentifiableSqlitePersistence.GetPageByFilter(ctx,
		c.composeFilter(filter), paging,
		"", "",
	)
}

func (c *MySqlitePersistence) GetOneById(ctx context.Context, key string) (item MyData, err error) {
	query := "SELECT * FROM " + c.QuotedTableName() + " WHERE \"key\"=$1"

	qResult, err := c.Client.QueryContext(ctx, query, key)
	if err != nil {
		return item, err
	}
	defer qResult.Close()

	if !qResult.Next() {
		return item, qResult.Err()
	}

	result, err := c.Overrides.ConvertToPublic(qResult)

	if err == nil {
		c.Logger.Trace(ctx, "Retrieved from %s with key = %s", c.TableName, key)
		return result, err
	}
	c.Logger.Trace(ctx, "Nothing found from %s with key = %s", c.TableName, key)
	return item, err
}
```

Alternatively you can store data in non-relational format using `IdentificableJsonSqlitePersistence`.
It stores data in tables with two columns - `id` with unique object id and `data` with object data serialized as JSON.
To access data fields you shall use `JSON_EXTRACT(data, '$.field')` expression.

```go
import "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/persistence"

type MySqlitePersistence struct {
	*persistence.IdentifableJsonSqlitePersistence[MyData, string]
}

func NewMySqlitePersistence() *MySqlitePersistence {
	c := &MySqlitePersistence{}
	c.IdentifableJsonSqlitePersistence = persistence.InheritIdentifiableJsonSqlitePersistence[]()[MyData, string](c, "myobjects")
	return c
}

func (c *MySqlitePersistence) DefineSchema() {
	c.EnsureTable("VARCHAR(32)", "JSON")
	c.EnsureIndex(c.TableName+"_json_key", map[string]string{"JSON_EXTRACT(data, '$.key')": "1"}, map[string]string{"unique": "true"})
}

func (c *MySqlitePersistence) composeFilter(filter data.FilterParams) string {
	criteria := make([]string, 0)

	if key, ok := filter.GetAsNullableString("key"); ok && key != "" {
		criteria = append(criteria, "JSON_EXTRACT(data, '$.key')='" + key + "'")
	}

	if id, ok := filter.GetAsNullableString("id"); ok && id != "" {
		criteria = append(criteria, "JSON_EXTRACT(data, '$.id')='" + id + "'")
	}

	if tempIds, ok := filter.GetAsNullableString("ids"); ok && tempIds != "" {
		ids := strings.Split(tempIds, ",")
		criteria = append(criteria, "JSON_EXTRACT(data, '$.id') IN ('"+strings.Join(ids, "','")+"')")
	}

	if len(criteria) > 0 {
		return strings.Join(criteria, " AND ")
	} else {
		return ""
	}
}

func (c *MySqlitePersistence) GetPageByFilter(ctx context.Context,
	filter data.FilterParams, paging data.PagingParams) (page data.DataPage[MyData], err error) {

	return c.IdentifiableSqlitePersistence.GetPageByFilter(ctx,
		c.composeFilter(filter), paging,
		"", "",
	)
}

func (c *MySqlitePersistence) GetOneById(ctx context.Context, key string) (item MyData, err error) {
	query := "SELECT * FROM " + c.QuotedTableName() + " WHERE JSON_EXTRACT(data, '$.key')=$1"

	qResult, err := c.Client.QueryContext(ctx, query, key)
	if err != nil {
		return item, err
	}
	defer qResult.Close()

	if !qResult.Next() {
		return item, qResult.Err()
	}

	result, err := c.Overrides.ConvertToPublic(qResult)

	if err == nil {
		c.Logger.Trace(ctx, "Retrieved from %s with key = %s", c.TableName, key)
		return result, err
	}
	c.Logger.Trace(ctx, "Nothing found from %s with key = %s", c.TableName, key)
	return item, err
}
```

Configuration for your microservice that includes SQLite persistence may look the following way.

```yaml
...
{{#if SQLITE_ENABLED}}
- descriptor: pip-services:connection:postgres:con1:1.0
  connection:
    database: {{SQLITE_DB}}{{#unless SQLITE_DB}}./data/app.db{{/unless}}
    
- descriptor: myservice:persistence:postgres:default:1.0
  dependencies:
    connection: pip-services:connection:postgres:con1:1.0
  table: {{SQLITE_TABLE}}{{#unless SQLITE_TABLE}}myobjects{{/unless}}
{{/if}}
...
```

## Develop

For development you shall install the following prerequisites:
* Golang v1.20+
* Visual Studio Code or another IDE of your choice
* Docker
* Git

Run automated tests:
```bash
go test -v ./test/...
```

Generate API documentation:
```bash
./docgen.ps1
```

Before committing changes run dockerized test as:
```bash
./test.ps1
./clear.ps1
```

## Contacts

The Golang version of Pip.Services is created and maintained by:
- **Levichev Dmitry**
- **Sergey Seroukhov**

The documentation is written by:
- **Levichev Dmitry**
