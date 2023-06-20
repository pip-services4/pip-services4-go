# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Data Handling Components for Golang

This module is a part of the [Pip.Services](http://pip.services.org) polyglot microservices toolkit.

It dynamic and static objects and data handling components.

The module contains the following packages:
- **Data** - data patterns
- **Keys**- object key (id) generators
- **Process**- data processing components
- **Query**- data query objects
- **Random** - random data generators
- **Validate** - validation patterns

<a name="links"></a> Quick links:

* [Memory persistence](http://docs.pipservices.org/conceptual/persistences/memory_persistence/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-data-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-data-go@latest
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
