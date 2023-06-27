package controllers

import (
	"strings"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccomands "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type CommandableSwaggerDocument struct {
	content string

	Commands []ccomands.ICommand

	Version   string
	BaseRoute string

	InfoTitle          string
	InfoDescription    string
	InfoVersion        string
	InfoTermsOfService string

	InfoContactName  string
	InfoContactUrl   string
	InfoContactEmail string

	InfoLicenseName string
	InfoLicenseUrl  string

	objectType map[string]any
}

func NewCommandableSwaggerDocument(baseRoute string, config *cconf.ConfigParams, commands []ccomands.ICommand) *CommandableSwaggerDocument {
	c := &CommandableSwaggerDocument{
		content:     "",
		Version:     "3.0.2",
		InfoVersion: "1",
		BaseRoute:   baseRoute,
		Commands:    make([]ccomands.ICommand, 0),
		objectType:  map[string]any{"type": "object"},
	}

	if commands != nil {
		c.Commands = commands
	}

	if config == nil {
		config = cconf.NewEmptyConfigParams()
	}

	c.InfoTitle = config.GetAsStringWithDefault("name", "CommandableHttpController")
	c.InfoDescription = config.GetAsStringWithDefault("description", "Commandable microservice")
	return c
}

func (c *CommandableSwaggerDocument) ToString() string {
	var data = map[string]any{
		"openapi": c.Version,
		"info": map[string]any{
			"title":             c.InfoTitle,
			"description":       c.InfoDescription,
			"version":           c.InfoVersion,
			"termsOfController": c.InfoTermsOfService,
			"contact": map[string]any{
				"name":  c.InfoContactName,
				"url":   c.InfoContactUrl,
				"email": c.InfoContactEmail,
			},
			"license": map[string]any{
				"name": c.InfoLicenseName,
				"url":  c.InfoLicenseUrl,
			},
		},
		"paths": c.createPathsData(),
	}

	c.WriteData(0, data)

	//console.log(c.content);
	return c.content
}

func (c *CommandableSwaggerDocument) createPathsData() map[string]any {
	var data = make(map[string]any, 0)

	for index := 0; index < len(c.Commands); index++ {
		command := c.Commands[index]

		var path = c.BaseRoute + "/" + command.Name()
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		data[path] = map[string]any{

			"post": map[string]any{
				"tags":        []any{c.BaseRoute},
				"operationId": command.Name(),
				"requestBody": c.createRequestBodyData(command),
				"responses":   c.createResponsesData(),
			},
		}
	}

	return data
}

func (c *CommandableSwaggerDocument) createRequestBodyData(command ccomands.ICommand) map[string]any {
	var schemaData = c.createSchemaData(command)
	if schemaData == nil {
		return nil
	}

	return map[string]any{
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": schemaData,
			},
		},
	}
}

func (c *CommandableSwaggerDocument) createSchemaData(command ccomands.ICommand) map[string]any {
	var schema = command.(*ccomands.Command).GetSchema().(*cvalid.ObjectSchema)

	if schema == nil || schema.Properties() == nil {
		return nil
	}

	return c.createPropertyData(schema, true)
}

func (c *CommandableSwaggerDocument) createPropertyData(schema ISchemaWithProperties, includeRequired bool) map[string]any {
	properties := make(map[string]any, 0)
	required := make([]string, 0)

	for _, property := range schema.Properties() {
		if property.Type() != nil {
			propertyName := property.Name()
			propertyType := property.Type()

			if _propertyType, ok := propertyType.(ISchemaBaseWithValueType); ok {
				properties[propertyName] = map[string]any{
					"type":  "array",
					"items": c.createPropertyTypeData(_propertyType.ValueType()),
				}
			} else {
				properties[propertyName] = c.createPropertyTypeData(propertyType)
			}

			if includeRequired && property.Required() {
				required = append(required, property.Name())
			}
		} else {
			properties[property.Name()] = c.objectType
		}
	}

	var data = map[string]any{
		"properties": properties,
	}

	if len(required) > 0 {
		data["required"] = required
	}

	return data
}

type ISchemaBaseWithValueType interface {
	cvalid.ISchemaBase
	ValueType() any
}

type ISchemaWithProperties interface {
	cvalid.ISchema
	Properties() []*cvalid.PropertySchema
}

func (c *CommandableSwaggerDocument) createPropertyTypeData(propertyType any) map[string]any {
	if _propertyType, ok := propertyType.(ISchemaWithProperties); ok {
		objectMap := c.createPropertyData(_propertyType, false)

		for k, v := range c.objectType {
			objectMap[k] = v
		}
		return objectMap
	}

	var typeCode cconv.TypeCode

	if _typeCode, ok := propertyType.(cconv.TypeCode); ok {
		typeCode = _typeCode
	} else {
		typeCode = cconv.TypeConverter.ToTypeCode(propertyType)
	}

	if typeCode == cconv.Map || typeCode == cconv.Unknown {
		typeCode = cconv.Object
	}

	switch typeCode {
	case cconv.Integer:
		return map[string]any{
			"type":   "integer",
			"format": "int32",
		}
	case cconv.Long:
		return map[string]any{
			"type":   "number",
			"format": "int64",
		}
	case cconv.Float:
		return map[string]any{
			"type":   "number",
			"format": "float",
		}
	case cconv.Double:
		return map[string]any{
			"type":   "number",
			"format": "double",
		}
	case cconv.DateTime:
		return map[string]any{
			"type":   "string",
			"format": "date-time",
		}
	case cconv.Boolean:
		return map[string]any{
			"type": "boolean",
		}
	default:
		return map[string]any{
			"type": cconv.TypeConverter.ToString(typeCode),
		}
	}

}

func (c *CommandableSwaggerDocument) createResponsesData() map[string]any {
	return map[string]any{

		"200": map[string]any{
			"description": "Successful response",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": map[string]any{
						"type": "object",
					},
				},
			},
		},
	}
}

func (c *CommandableSwaggerDocument) WriteData(indent int, data map[string]any) {

	for key, value := range data {
		if val, ok := value.(string); ok {
			c.writeAsString(indent, key, val)
		} else {
			if arr, ok := value.([]any); ok {
				if len(arr) > 0 {
					c.WriteName(indent, key)
					for index := 0; index < len(arr); index++ {
						item := arr[index].(string)
						c.writeArrayItem(indent+1, item, false)
					}
				}
			} else {
				if m, ok := value.(map[string]any); ok {
					notEmpty := false
					for _, v := range m {
						if v != nil {
							notEmpty = true
							break
						}
					}
					if notEmpty {
						c.WriteName(indent, key)
						c.WriteData(indent+1, m)
					}
				} else {
					c.writeAsObject(indent, key, value)
				}
			}
		}
	}
}

func (c *CommandableSwaggerDocument) WriteName(indent int, name string) {
	var spaces = c.GetSpaces(indent)
	c.content += spaces + name + ":\n"
}

func (c *CommandableSwaggerDocument) writeArrayItem(indent int, name string, isObjectItem bool) {
	var spaces = c.GetSpaces(indent)
	c.content += spaces + "- "

	if isObjectItem {
		c.content += name + ":\n"
	} else {
		c.content += name + "\n"
	}
}

func (c *CommandableSwaggerDocument) writeAsObject(indent int, name string, value any) {
	if value == nil {
		return
	}

	var spaces = c.GetSpaces(indent)
	c.content += spaces + name + ": " + cconv.StringConverter.ToString(value) + "\n"
}

func (c *CommandableSwaggerDocument) writeAsString(indent int, name string, value string) {
	if value == "" {
		return
	}

	var spaces = c.GetSpaces(indent)
	c.content += spaces + name + ": '" + value + "'\n"
}

func (c *CommandableSwaggerDocument) GetSpaces(length int) string {
	return strings.Repeat(" ", length*2)
}
