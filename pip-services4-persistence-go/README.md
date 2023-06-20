# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Data Persistence Components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit. It contains generic interfaces for data access components as well as abstract implementations for in-memory and file persistence.

The persistence components come in two kinds. The first kind is a basic persistence that can work with any object types and provides only minimal set of operations. 
The second kind is so called "identifieable" persistence with works with "identifable" data objects, i.e. objects that have unique ID field. The identifiable persistence provides a full set or CRUD operations that covers most common cases.

The module contains the following packages:
- **Read** - generic data reading interfaces.
- **Write** - generic data writing interfaces.
- **Persistence** - in-memory and file persistence components, as well as JSON persister class.

<a name="links"></a> Quick links:

* [Memory persistence](http://docs.pipservices.org/conceptual/persistences/memory_persistence/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-persistence-go@latest
```

As an example, lets implement persistence for the following data object.

```go
type MyData struct {
	id      string
	key     string
	content string
}

func (c *MyData) GetId() string {
	return c.id
}
```

Our persistence component shall implement the following interface with a basic set of CRUD operations.

```go
type IMyPersistence interface {
	GetPageByFilter(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[MyData], err error)
	GetOneById(ctx context.Context, id string) (item MyData, err error)
	GetOneByKey(ctx context.Context, key string) (item MyData, err error)
	Create(ctx context.Context, item MyData) (result MyData, err error)
	Update(ctx context.Context, item MyData) (result MyData, err error)
	DeleteById(ctx context.Context, id string) (item MyData, err error)
}
```

To implement in-memory persistence component you shall inherit `IdentifiableMemoryPersistence`. 
Most CRUD operations will come from the base class. You only need to override `GetPageByFilter` method with a custom filter function.
And implement a `GetOneByKey` custom persistence method that doesn't exist in the base class.

```go
type MyMemoryPersistence struct {
	*persistence.IdentifiableMemoryPersistence[MyData, string]
}

func NewMyMemoryPersistence() *MyMemoryPersistence {
	return &MyMemoryPersistence{
		IdentifiableMemoryPersistence: persistence.NewIdentifiableMemoryPersistence[MyData, string](),
	}
}

func (c *MyMemoryPersistence) composeFilter(filter cdata.FilterParams) func(item MyData) bool {
	if &filter == nil {
		filter = *cdata.NewEmptyFilterParams()
	}

	id, idOk := filter.GetAsNullableString("id")
	ids, idsOk := filter.GetAsNullableArray("ids")
	key, keyOK := filter.GetAsNullableString("key")

	return func(item MyData) bool {
		if idOk && item.id != id {
			return false
		}
		if idsOk && ids.Contains(item.id) {
			return false
		}
		if keyOK && item.key != key {
			return false
		}

		return true
	}
}

func (c *MyMemoryPersistence) GetPageByFilter(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[main.MyData], err error) {
	return c.IdentifiableMemoryPersistence.GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
}

func (c *MyMemoryPersistence) GetOneByKey(ctx context.Context, key string) (item MyData, err error) {
	var resItem *MyData

	for _, item := range c.Items {
		if item.key == key {
			resItem = &item
			break
		}
	}

	if resItem != nil {
		c.Logger.Trace(ctx, "Found object by key=%s", key)
	} else {
		c.Logger.Trace(ctx, "Cannot find by key=%s", key)
	}

	return *resItem, nil
}
```

It is easy to create file persistence by adding a persister object to the implemented in-memory persistence component.

```go
type MyFilePersistence struct {
	*MyMemoryPersistence
	persister *persistence.JsonFilePersister[MyData]
}

func NewMyFilePersistence(path string) *MyFilePersistence {
	c := &MyFilePersistence{}
	c.persister = persistence.NewJsonFilePersister[MyData](path)
	c.Loader = c.persister
	c.Saver = c.persister
	return c
}

func (c *MyFilePersistence) Configure(ctx context.Context, config *config.ConfigParams) {
	c.IdentifiableMemoryPersistence.Configure(ctx, config)
	c.persister.Configure(ctx, config)
}
```

Configuration for your microservice that includes memory and file persistence may look the following way.

```yml
...
{{#if MEMORY_ENABLED}}
- descriptor: "myservice:persistence:memory:default:1.0"
{{/if}}

{{#if FILE_ENABLED}}
- descriptor: "myservice:persistence:file:default:1.0"
  path: {{FILE_PATH}}{{#unless FILE_PATH}}"../data/data.json"{{/unless}}
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
