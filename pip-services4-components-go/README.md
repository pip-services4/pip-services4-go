# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Component definitions for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The Components module contains standard component definitions that can be used to build applications and services.

The module contains the following packages:

- **Auth** - authentication credential stores
- **Build** - factories
- **Cache** - distributed cache
- **Component** - the root package
- **Config** - configuration readers
- **Connect** - connection discovery services
- **Count** - performance counters
- **Info** - context info
- **Lock** - distributed locks
- **Log** - logging components
- **Test** - test components

<a name="links"></a> Quick links:

* [Logging](http://docs.pipservices.org/getting_started/recipes/logging/)
* [Configuration](http://docs.pipservices.org/concepts/configuration/) 
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-components-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-components-go@latest
```

Example how to use Logging and Performance counters. Here we are going to use CompositeLogger and CompositeCounters components. They will pass through calls to loggers and counters that are set in references.

```go
import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/count"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/log"
)

type MyComponent struct {
	logger   *log.CompositeLogger
	counters *count.CompositeCounters
}

func (c *MyComponent) Configure(ctx context.Context, config *config.ConfigParams) {
	c.logger.Configure(ctx, config)
}

func (c *MyComponent) SetReferences(ctx context.Context, references refer.IReferences) {
	c.logger.SetReferences(ctx, references)
	c.counters.SetReferences(ctx, references)
}

func (c *MyComponent) MyMethod(ctx context.Context, param1 any) {
	c.logger.Trace(ctx, "Executed method mycomponent.mymethod")
	c.counters.Increment(ctx, "mycomponent.mymethod.exec_count", 1)
	timing := c.counters.BeginTiming(ctx, "mycomponent.mymethod.exec_time")
	defer timing.EndTiming(ctx)
	// ....

	if err != nil {
		c.logger.Error(ctx, err, "Failed to execute mycomponent.mymethod")
		c.counters.Increment(ctx, "mycomponent.mymethod.error_count", 1)
	}
}
```

Example how to get connection parameters and credentials using resolvers. The resolvers support "discovery_key" and "store_key" configuration parameters to retrieve configuration from discovery services and credential stores respectively.

```go
package main

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/auth"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/connect"
)

func main() {
	// Using the component
	myComponent := NewMyComponent()

	myComponent.Configure(context.Background(), config.NewConfigParamsFromTuples(
		"connection.host", "localhost",
		"connection.port", 1234,
		"credential.username", "anonymous",
		"credential.password", "pass123",
	))

	err := myComponent.Open(context.Background(), "123")
}

type MyComponent struct {
	connectionResolver *connect.ConnectionResolver
	credentialResolver *auth.CredentialResolver
}

func NewMyComponent() *MyComponent {
	return &MyComponent{
		connectionResolver: connect.NewEmptyConnectionResolver(),
		credentialResolver: auth.NewEmptyCredentialResolver(),
	}
}

func (c *MyComponent) Configure(ctx context.Context, config *config.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)
}

func (c *MyComponent) SetReferences(ctx context.Context, references refer.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.credentialResolver.SetReferences(ctx, references)
}

// ...

func (c *MyComponent) IsOpen() bool {
	panic("not implemented") // TODO: Implement
}

func (c *MyComponent) Open(ctx context.Context) error {
	connection, err := c.connectionResolver.Resolve(ctx)
	credential, err := c.credentialResolver.Lookup(ctx)

	host := connection.Host()
	port := connection.Port()
	user := credential.Username()
	pass := credential.Password()
}

func (c *MyComponent) Close(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

```

Example how to use caching and locking. Here we assume that references are passed externally.

```go
package main

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/cache"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/lock"
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

If you need to create components using their locators (descriptors) implement component factories similar to the example below.

```go
package main

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
)

var MyComponentDescriptor = refer.NewDescriptor("myservice", "mycomponent", "default", "*", "1.0")

func NewMyFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(MyComponentDescriptor, NewMyComponent)

	return factory
}

func main() {
	// Using the factory
	myFactory := NewMyFactory()
	myComponent1, err = myFactory.Create(refer.NewDescriptor("myservice", "mycomponent", "default", "myComponent1", "1.0"))
	myComponent2, err := myFactory.Create(refer.NewDescriptor("myservice", "mycomponent", "default", "myComponent2", "1.0"))
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

The documentation is written by **Levichev Dmitry**.
