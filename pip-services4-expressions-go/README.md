# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Tokenizers, parsers and expression calculators Golang

This module is a part of the [Pip.Services](http://pip.services.org) polyglot microservices toolkit.
It provides syntax and lexical analyzers and expression calculator optimized for repeated calculations.

The module contains the following packages:
- **Calculator** - Expression calculator
- **CSV** - CSV tokenizer
- **IO** - input/output utility classes to support lexical analysis
- **Mustache** - Mustache templating engine
- **Tokenizers** - lexical analyzers to break incoming character streams into tokens
- **Variants** - dynamic objects that can hold any values and operators for them

<a name="links"></a> Quick links:

* [API Reference](https://godoc.org/github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/)
* [Change Log](CHANGELOG.md)
* [Get Help](http://docs.pipservices.org/get_help/)
* [Contribute](http://docs.pipservices.org/contribute/)

## Use

Get the package from the Github repository:
```bash
go get -u github.com/pip-services4/pip-services4-go/pip-services4-expressions-go@latest
```

The example below shows how to use expression calculator to dynamically
calculate user-defined expressions.

```golang
import (
    "fmt"

    calc "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator"
    vars "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/variables"
    variants "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
)

...
calculator := calc.NewExpressionCalculator()

calculator.SetExpression("A + b / (3 - Max(-123, 1)*2)")

variables := vars.NewVariableCollection()
variables.Add(vars.NewVariable("A", variants.NewVariantFromInteger(1)))
variables.Add(vars.NewVariable("B", variants.NewVariantFromString("3")))

result, err := calculator.EvaluateWithVariables(variables)
if err != nil {
    fmt.Println("Failed to calculate the expression")
} else {
    fmt.Println("The result of the expression is " + result.AsString())
}
...
```

This is an example to process mustache templates.

```golang
import (
    "fmt"

    mustache "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache"
)

template := mustache.NewMustacheTemplate()
template.SetTemplate("Hello, {{{NAME}}}{{#ESCLAMATION}}!{{/ESCLAMATION}}{{#unless ESCLAMATION}}.{{/unless}}")
result, err := template.EvaluateWithVariables(map[string]string{ "NAME": "Mike", "ESCLAMATION": "true" })
if err != nil {
    fmt.Println("Failed to evaluate mustache template")
} else {
    fmt.Println("The result of template evaluation is '" + result + "'")
}
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
