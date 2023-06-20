# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> IoC container for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit. It provides an inversion-of-control (IoC) container to facilitate the development of services and applications composed of loosely coupled components.

The module containes a basic in-memory container that can be embedded inside a service or application, or can be run by itself.
The second container type can run as a system level process and can be configured via command line arguments.
Also it can be used to create docker containers.

The containers can read configuration from JSON or YAML files use it as a recipe for instantiating and configuring components.
Component factories are used to create components based on their locators (descriptor) defined in the container configuration.
The factories shall be registered in containers or dynamically in the container configuration file.

The module contains the following packages:

- **Container** - Component container and container as a system process
- **Build** - Container default factory
- **Config** - Container configuration
- **Refer** - Container references

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/concepts/configuration/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-container-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)

## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-container-go@latest
```

Create a factory to create components based on their locators (descriptors).

```go
import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
)

var MyComponentDescriptor = refer.NewDescriptor("myservice", "mycomponent", "default", "*", "1.0")

func NewMyFactory() *build.Factory {
	f := build.NewFactory()

	f.RegisterType(MyComponentDescriptor, NewMyComponent)
}
```

Then create a process container and register the factory there. You can also register factories defined in other modules if you plan to include external components into your container.

```go
import (
	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	rpcbuild "github.com/pip-services3-gox/pip-services3-rpc-gox/build"
)

type MyProcess struct {
	*container.ProcessContainer
}

func NewMyProcess() *MyProcess {
	c := &MyProcess{
		ProcessContainer: container.NewProcessContainer("myservice", "My service running as a process"),
	}

	c.AddFactory(rpcbuild.NewDefaultRpcFactory())
	c.AddFactory(NewMyFactory())
}
```

Define YAML configuration file with components and their descriptors. The configuration file is pre-processed using [Handlebars templating](https://github.com/pip-services3-gox/pip-services3-expressions-gox) that allows to inject configuration parameters or dynamically include/exclude components using conditional blocks. The values for the templating engine are defined via process command line arguments or via environment variables. Support for environment variables works well in docker or other containers like AWS Lambda functions.


```yml
---
# Context information
- descriptor: "pip-services:context-info:default:default:1.0"
  name: myservice
  description: My service running in a process container

# Console logger
- descriptor: "pip-services:logger:console:default:1.0"
  level: {{LOG_LEVEL}}{{^LOG_LEVEL}}info{{/LOG_LEVEL}}

# Performance counters that posts values to log
- descriptor: "pip-services:counters:log:default:1.0"
  
# My component
- descriptor: "myservice:mycomponent:default:default:1.0"
  param1: XYZ
  param2: 987
  
{{#if HTTP_ENABLED}}
# HTTP endpoint version 1.0
- descriptor: "pip-services:endpoint:http:default:1.0"
  connection:
    protocol: "http"
    host: "0.0.0.0"
    port: {{HTTP_PORT}}{{^HTTP_PORT}}8080{{/HTTP_PORT}}

 # Default Status
- descriptor: "pip-services:status-service:http:default:1.0"

# Default Heartbeat
- descriptor: "pip-services:heartbeat-service:http:default:1.0"
{{/if}}
```

To instantiate and run the container we need a simple process launcher.

```go
package main

import (
	"context"
	"os"
)

func main() {
	proc := NewMyProcess()
	proc.SetConfigPath("./config/config.yml")
	proc.Run(context.Background(), os.Args)
}

```

And, finally, you can run your service launcher as

```bash
go run main.go
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
- **Volodymyr Tkachenko**
- **Sergey Seroukhov**
- **Mark Zontak**

The documentation is written by:
- **Levichev Dmitry**
