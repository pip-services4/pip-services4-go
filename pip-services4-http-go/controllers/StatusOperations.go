package controllers

import (
	"context"
	"net/http"
	"time"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cinfo "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// StatusOperations helper class for status service
type StatusOperations struct {
	*RestOperations
	startTime   time.Time
	references2 crefer.IReferences
	contextInfo *cinfo.ContextInfo
}

// NewStatusOperations creates new instance of StatusOperations
func NewStatusOperations() *StatusOperations {
	c := StatusOperations{}
	c.RestOperations = NewRestOperations()
	c.startTime = time.Now()
	c.DependencyResolver.Put(
		context.Background(),
		"context-info",
		crefer.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"),
	)
	return &c
}

// SetReferences  sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences references to locate the component dependencies.
func (c *StatusOperations) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.references2 = references
	c.RestOperations.SetReferences(ctx, references)

	depRes := c.DependencyResolver.GetOneOptional("context-info")
	if depRes != nil {
		if ctxInfo, ok := depRes.(*cinfo.ContextInfo); ok {
			c.contextInfo = ctxInfo
		}
	}
}

// GetStatusOperation return function for get status
func (c *StatusOperations) GetStatusOperation() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		c.Status(res, req)
	}
}

// Status method handles status requests
//
//	Parameters:
//		- req *http.Request  an HTTP request
//		- res  http.ResponseWriter  an HTTP response
func (c *StatusOperations) Status(res http.ResponseWriter, req *http.Request) {

	id := ""
	if c.contextInfo != nil {
		id = c.contextInfo.ContextId
	}

	name := "Unknown"
	if c.contextInfo != nil {
		name = c.contextInfo.Name
	}

	description := ""
	if c.contextInfo != nil {
		description = c.contextInfo.Description
	}

	uptime := time.Since(c.startTime)

	properties := make(map[string]string)
	if c.contextInfo != nil {
		properties = c.contextInfo.Properties
	}

	var components []string
	if c.references2 != nil {
		for _, locator := range c.references2.GetAllLocators() {
			components = append(components, cconv.StringConverter.ToString(locator))
		}
	}

	status := make(map[string]any)

	status["id"] = id
	status["name"] = name
	status["description"] = description
	status["start_time"] = cconv.StringConverter.ToString(c.startTime)
	status["current_time"] = cconv.StringConverter.ToString(time.Now())
	status["uptime"] = uptime
	status["properties"] = properties
	status["components"] = components

	c.SendResult(res, req, status, nil)
}
