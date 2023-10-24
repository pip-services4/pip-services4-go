# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Business Logic Components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The Logic module contains standard component definitions to handle complex business transactions.

The module contains the following packages:
- **Cache** - distributed cache
- **Lock** -  distributed lock components
- **State** -  distributed state management components

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/component_configuration/) 
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-logic-go)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-logic-go@latest
```

Example how to use caching and locking.
Here we assume that references are passed externally.

```go
package main

import (
	"context"

	"github.com/pip-services4-go/pip-services4-commons-go/refer"
	"github.com/pip-services4-go/pip-services4-components-go/cache"
	"github.com/pip-services4-go/pip-services4-components-go/lock"
)

func main() {
	// Use the component
	myComponent := NewMyComponent()

	myComponent.SetReferences(context.Background(), refer.NewReferencesFromTuples(context.Background(),
		refer.NewDescriptor("pip-services", "cache", "memory", "default", "1.0"), cache.NewMemoryCache[any](),
		refer.NewDescriptor("pip-services", "lock", "memory", "default", "1.0"), lock.NewMemoryLock(),
	))

	result, err := myComponent.MyMethod(context.Background(), "123", "my_param")
}

type MyComponent struct {
	cache cache.ICache[any]
	lock  lock.ILock
}

func NewMyComponent() *MyComponent {
	return &MyComponent{}
}

func (c *MyComponent) SetReferences(ctx context.Context, references refer.IReferences) {
	res, errDescr := references.GetOneRequired(refer.NewDescriptor("*", "cache", "*", "*", "1.0"))
	if errDescr != nil {
		panic(errDescr)
	}
	c.cache = res.(cache.ICache[any])

	res, errDescr = references.GetOneRequired(refer.NewDescriptor("*", "lock", "*", "*", "1.0"))
	if errDescr != nil {
		panic(errDescr)
	}
	c.lock = res.(lock.ILock)
}

func (c *MyComponent) MyMethod(ctx context.Context, param1 any) (any, error) {
	// First check cache for result
	result, err := c.cache.Retrieve(ctx, "mykey")
	if result != nil || err != nil {
		return result, err
	}

	// Lock..
	err = c.lock.AcquireLock(ctx, "mykey", 1000, 1000)
	if err != nil {
		return result, err
	}

	// Do processing
	// ...

	// Store result to cache async
	_, err = c.cache.Store(ctx, "mykey", result, 3600000)
	if err != nil {
		return result, err
	}

	// Release lock async
	err = c.lock.ReleaseLock(ctx, "mykey")
	if err != nil {
		return result, err
	}
	return result, nil
}

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

The library is created and maintained by **Sergey Seroukhov**.

The documentation is written by **Danyil Tretiakov** and **Levichev Dmitry**.
