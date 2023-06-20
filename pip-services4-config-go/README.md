# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Config Components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The Config module contains configuration component definitions that can be used to build applications and services.

The module contains the following packages:
- **Auth** - authentication credential stores
- **Config** - configuration readers and managers, whose main task is to deliver configuration parameters to the application from wherever they are being stored
- **Connect** - connection discovery and configuration services

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/component_configuration/) 
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-config-go)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-config-go@latest
```
Then you are ready to start using the Pip.Services patterns to augment your backend code.

For instance, here is how you can implement a component, that receives configuration, get assigned references,
can be opened and closed using the patterns from this module.

```go
package main

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services3-components-gox/config"
	"github.com/pip-services4/pip-services4-go/pip-services3-components-gox/refer"
	"github.com/pip-services4/pip-services4-go/pip-services3-config-gox/auth"
	"github.com/pip-services4/pip-services4-go/pip-services3-config-gox/connect"
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

	err := myComponent.Open(context.Background())
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
