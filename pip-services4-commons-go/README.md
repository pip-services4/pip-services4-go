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

## Develop

For development you shall install the following prerequisites:
* Golang v1.20+
* Visual Studio Code or another IDE of your choice
* Docker
* Git

Run automated tests:
```bash
go test -v ./test/...
``

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
