# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> ElasticSearch components for PipServices in Go

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The Elasticsearch module contains logging components with data storage on the Elasticsearch server.

The module contains the following packages:
- **Build** - contains a factory for the construction of components
- **Log** - Logging components

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/)
* [Logging](http://docs.pipservices.org/getting_started/recipes/logging/)
* [Virtual memory configuration](https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-elasticsearch-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)

## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-elasticsearch-go@latest
```

Microservice components shall perform logging usual way using CompositeLogger component. The CompositeLogger will get ElasticSearchLogger from references and will redirect log messages there among other destinations.

```go
import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

type MyComponent struct {
	logger *log.CompositeLogger
}

func (c *MyComponent) Configure(ctx context.Context, config *config.ConfigParams) {
	c.logger.Configure(ctx, config)
}
func (c *MyComponent) SetReferences(ctx context.Context, references refer.IReferences) {
	c.logger.SetReferences(ctx, references)
}

func (c *MyComponent) MyMethod(ctx context.Context, param1 any) (any, error) {
	c.logger.Trace(ctx, "Executed method mycomponent.mymethod")
	// ....
}
```

Configuration for your microservice that includes ElasticSearch logger may look the following way.

```yml
...
{{#if ELASTICSEARCH_ENABLED}}
- descriptor: pip-services:logger:elasticsearch:default:1.0
  connection:
    uri: {{{ELASTICSEARCG_SERVICE_URI}}}
    host: {{{ELASTICSEARCH_SERVICE_HOST}}}{{#unless ELASTICSEARCH_SERVICE_HOST}}localhost{{/unless}}
    port: {{ELASTICSEARCG_SERVICE_PORT}}{{#unless ELASTICSEARCH_SERVICE_PORT}}9200{{/unless}}\ 
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

The library is created and maintained by **Sergey Seroukhov** and **Levichev Dmitry**.

The documentation is written by:
- **Levichev Dmitry**