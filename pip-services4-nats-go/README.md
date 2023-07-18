# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> NATS messaging components for Pip.Services in Go

This library is a part of [Pip.Services](https://github.com/pip-services/pip-services) project.
The Nats module contains a set of components for working with the message queue via NATS server [https://nats.io/](https://nats.io/).

The module contains the following packages:
- **Build** - Factory for constructing module components
- **Connect** - Components for creating and configuring a connection with Nats
- **Queues** - Message Queuing components that implement the standard [Messaging](https://github.com/pip-services4/pip-services4-go/pip-services4-messaging-go) module interface

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-nats-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)

## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-nats-go@latest
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
- **Sergey Seroukhov**

The documentation is written by:
- **Levichev Dmitry**
