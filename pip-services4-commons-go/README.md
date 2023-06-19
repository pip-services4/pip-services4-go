# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Portable Abstractions and Patterns for Golang

This module is a part of the [Pip.Services](http://pip.services.org) polyglot microservices toolkit.
It provides a set of basic patterns used in microservices or backend services.
Also the module implemenets a reasonably thin abstraction layer over most fundamental functions across
all languages supported by the toolkit to facilitate symmetric implementation.

The module contains the following packages:

- **Convert** - portable value converters
- **Data** - data patterns
- **Errors**- application errors
- **Reflect** - portable reflection utilities

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/conceptual/configuration/component_configuration/)
* [Locator Pattern](http://docs.pipservices.org/conceptual/component/component_references/)
* [Component Lifecycle](http://docs.pipservices.org/conceptual/component/component_lifecycle/)
* [Data Patterns](http://docs.pipservices.org/conceptual/persistences/memory_persistence/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-commons-go)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-commons-go@latest
```
Then you are ready to start using the Pip.Services patterns to augment your backend code.

For instance, here is how you can implement a component, that receives configuration, get assigned references,
can be opened and closed using the patterns from this module.

```go

import (
	"context"
	"fmt"

	"github.com/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4-go/pip-services4-components-go/refer"
)

type MyComponentA struct {
	param1           string
	param2           int
	anotherComponent MyComponentB
	opened           bool
}

func NewMyComponentA() *MyComponentA {
	return &MyComponentA{
		param1: "ABC",
		param2: 123,
		opened: false,
	}
}

type MyComponentB struct{
    // ...
}

func (c *MyComponentA) Configure(ctx context.Context, config *config.ConfigParams) {
	c.param1 = config.GetAsStringWithDefault("param1", c.param1)
	c.param2 = config.GetAsIntegerWithDefault("param2", c.param2)
}

func (c *MyComponentA) SetReferences(ctx context.Context, references refer.IReferences) {
	res, err := references.GetOneRequired(refer.NewDescriptor("myservice", "mycomponent-b", "*", "*", "1.0"))
	if err != nil {
		panic(err)
	}

	c.anotherComponent = res.(MyComponentB)
}

func (c *MyComponentA) IsOpen() bool {
	return c.opened
}

func (c *MyComponentA) Open(ctx context.Context) error {
	c.opened = true
	fmt.Println("MyComponentA has been opened.")
	return nil
}

func (c *MyComponentA) Close(ctx context.Context) error {
	c.opened = false
	fmt.Println("MyComponentA has been closed.")
	return nil
}

```

Then here is how the component can be used in the code

```go
package main

import (
	"context"
	"fmt"

	"github.com/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4-go/pip-services4-components-go/refer"
)

func main() {
	myComponentA := NewMyComponentA()

	// Configure the component
	myComponentA.Configure(context.Background(), config.NewConfigParamsFromTuples(
		"param1", "XYZ",
		"param2", 987,
	))

	// Set references to the component
	myComponentA.SetReferences(context.Background(),
		refer.NewReferencesFromTuples(context.Background(),
			refer.NewDescriptor("myservice", "mycomponent-b", "default", "default", "1.0"), &MyComponentB{},
		),
	)

	// Open the component
	err := myComponentA.Open(context.Background(), "123")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("MyComponentA has been opened.")
	}
}
```

## Develop

For development you shall install the following prerequisites:
* Golang v1.18+
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
