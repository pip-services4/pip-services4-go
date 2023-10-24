# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Swagger UI for Golang

This module is a part of the [Pip.Services](http://pipservices.org) polyglot microservices toolkit.

The swagger module provides a Swagger UI that can be added into microservices and seamlessly integrated with existing REST and Commandable HTTP services.

The module contains the following packages:
- **Build** - Swagger service factory
- **Services** - Swagger UI service

<a name="links"></a> Quick links:

* [Configuration](http://docs.pipservices.org/conceptual/configuration/component_configuration/)
* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)


## Use

Install the Go package as
```bash
go get github.com/pip-services4/pip-services4-go/pip-services4-swagger-go
```

Develop a RESTful service component. For example, it may look the following way.
In the `Register` method we load an Open API specification for the service.
If you are planning to use the REST service as a library, then embed the Open API specification
as a resource using a library like [Statik](https://github.com/rakyll/statik).
You can also enable swagger by default in the constractor by setting `SwaggerEnable` property.
```golang
type MyRestController struct {
	*cservices.RestController
}

func NewMyRestController() *MyRestController {
	c := MyRestController{}
	c.RestController = cservices.InheritRestController(&c)
	c.BaseRoute = "myservice"
  c.SwaggerEnable = true
	return &c
}

func (c *MyRestController) greeting(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	name := params.Get("name")
  result := "Hello, " + name + "!"
	c.SendResult(res, req, result, nil)
}

func (c *MyRestController) Register() {
	c.RegisterOpenApiSpecFromFile("./services/myservice.yml")

	c.RegisterRoute(
		"get", "/greeting",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("name", cconv.String).Schema,
		c.greeting,
	)
}
```

The Open API specification for the service shall be prepared either manually
or using [Swagger Editor](https://editor.swagger.io/)
```yaml
openapi: '3.0.2'
info:
  title: 'MyService'
  description: 'MyService REST API'
  version: '1'
paths:
  /myservice/greeting:
    get:
      tags:
        - myservice
      operationId: 'greeting'
      parameters:
      - name: trace_id
        in: query
        description: Trace ID
        required: false
        schema:
          type: string
      - name: name
        in: query
        description: Name of a person
        required: true
        schema:
          type: string
      responses:
        200:
          description: 'Successful response'
          content:
            application/json:
              schema:
                type: 'string'
```

Include Swagger service into `config.yml` file and enable swagger for your REST or Commandable HTTP services.
Also explicitely adding HttpEndpoint allows to share the same port betwee REST services and the Swagger service.
```yaml
---
...
# Shared HTTP Endpoint
- descriptor: "pip-services:endpoint:http:default:1.0"
  connection:
    protocol: http
    host: localhost
    port: 8080

# Swagger Service
- descriptor: "pip-services:swagger-controller:http:default:1.0"

# My RESTful Service
- descriptor: "myservice:service:rest:default:1.0"
  swagger:
    enable: true
```

Finally, remember to add factories to your container, to allow it creating required components.
```golang
...
import (
	cproc "github.com/pip-services4-go/pip-services4-container-go/container"
	hbuild "github.com/pip-services4/pip-services4-go/pip-services4-http-go/build"
	sbuild "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/build"
  ...
)

type MyProcess struct {
	*cproc.ProcessContainer
}

func NewMyProcess() *MyProcess {
	c := MyProcess{}
	c.ProcessContainer = cproc.NewProcessContainer("myservice", "MyService microservice")
	c.AddFactory(factory.NewMyServiceFactory())
	c.AddFactory(hbuild.NewDefaultHttpFactory())
	c.AddFactory(sbuild.NewDefaultSwaggerFactory())
	return &c
}
```

Launch the microservice and open the browser to open the Open API specification at
[http://localhost:8080/greeting/swagger](http://localhost:8080/greeting/swagger)

Then open the Swagger UI using the link [http://localhost:8080/swagger](http://localhost:8080/swagger).
The result shall look similar to the picture below.

<img src="swagger-ui.png"/>

## Develop

For development you shall install the following prerequisites:
* Golang 1.20+
* Visual Studio Code or another IDE of your choice
* Docker

Generate embedded resources:
```bash
./resgen.ps1
```

Run automated tests:
```bash
go test ./test/...
```

Generate API documentation:
```bash
./docgen.ps1
```

Before committing changes run dockerized build and test as:
```bash
./build.ps1
./test.ps1
./clear.ps1
```

## Contacts

The Golang version of Pip.Services is created and maintained by:
- **Sergey Seroukhov**
- **Danil Prisiazhnyi**
