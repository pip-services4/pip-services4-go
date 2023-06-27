# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> GRPC Communication Components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The grpc module is used to organize synchronous data exchange using calls through the gRPC protocol. It has implementations of both the server and client parts.

The module contains the following packages:

- **Build** - factories for creating gRPC controllers
- **Clients** - basic client components that use the gRPC protocol and Commandable pattern through gRPC
- **Controllers** - basic controller implementations for connecting via the gRPC protocol and using the Commandable pattern via gRPC

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/component_configuration/) 
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-grpc-go)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-grpc-go@latest
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
