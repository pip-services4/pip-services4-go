# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> AWS specific components for Golang


This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

This module contains components for supporting work with the AWS cloud platform.

The module contains the following packages:
- [**Build**](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/build) - factories for constructing module components
- [**Clients**](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/clients) - client components for working with Lambda AWS
- [**Connect**](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/connect) - components of installation and connection settings
- [**Container**](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/containers) - components for creating containers for Lambda server-side AWS functions
- [**Count**](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/count) - components of working with counters (metrics) with saving data in the CloudWatch AWS service
- [**Log**](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/log) - logging components with saving data in the CloudWatch AWS service

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/)
* [aws-doc-sdk-examples](https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/lambda_functions/blank-go)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-aws-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)

## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-aws-go@latest
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