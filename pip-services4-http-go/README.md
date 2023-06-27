# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> HTTP/REST Communication Components for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The rpc module provides the synchronous communication using local calls or the HTTP(S) protocol. It contains both server and client side implementations.

The module contains the following packages:
- **Auth** - authentication and authorization components
- **Build** - HTTP service factory
- **Clients** - mechanisms for retrieving connection settings from the microservice’s configuration and providing clients and services with these settings
- **Connect** - helper module to retrieve connections for HTTP-based services and clients
- **Services** - basic implementation of services for connecting via the HTTP/REST protocol and using the Commandable pattern over HTTP

<a name="links"></a> Quick links:

* [Your first microservice in Node.js](http://docs.pipservices.org/toolkit/getting_started/your_first_microservice/) 
* [Data Microservice. Step 5](http://docs.pipservices.org/toolkit/tutorials/data_microservice/step5/)
* [Microservice Facade](http://docs.pipservices.org/toolkit/tutorials/microservice_facade/) 
* [Client Library. Step 3](http://docs.pipservices.org/toolkit/tutorials/client_library/step2/)
* [Client Library. Step 4](http://docs.pipservices.org/toolkit/tutorials/client_library/step3/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-http-go)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-http-go@latest
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

The documentation is written by **Danyil Tretiakov**, **Levichev Dmitry**, **Danil Prisiazhnyi**
