# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Azure specific components for Golang

This module is a part of the [Pip.Services](http://pip.services.org) polyglot microservices toolkit.

Contains packages used to create containers and services that do operations via the Azure Function protocol.

The module contains the following packages:
- **Clients** - client components for working with Azure cloud Functions.
- **Connect** - components for installation and connection settings.
- **Containers** - contains classes that act as containers to instantiate and run components.
- **Services** - contains interfaces and classes used to create services that do operations via the Azure Function protocol.

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/conceptual/configuration/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-azure-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-azure-go@latest
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

The library is created and maintained by **Sergey Seroukhov** and **Danil Prisiazhnyi**.
