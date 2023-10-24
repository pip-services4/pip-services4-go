# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Observability Components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The Observability module contains observability component definitions that can be used to build applications and services.

The module contains the following packages:
- **Count** - performance counters
- **Log** - basic logging components that provide console and composite logging, as well as an interface for developing custom loggers
- **Trace** - tracing components

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/component_configuration/) 
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-observability-go)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-observability-go@latest
```
Example how to use Logging and Performance counters.
Here we are going to use CompositeLogger and CompositeCounters components.
They will pass through calls to loggers and counters that are set in references.

```go
import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
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

func (c *MyComponent) MyMethod(ctx context.Context,  param1 any) {
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
