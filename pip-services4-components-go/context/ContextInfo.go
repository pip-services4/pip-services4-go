package context

import (
	"context"
	"os"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// ContextInfo Context information component that provides detail information
// about execution context: container or/and process.
// Most often ContextInfo is used by logging and performance counters to identify
// source of the collected logs and metrics.
//
//	Configuration parameters:
//		- name: the context (container or process) name
//		- description: human-readable description of the context
//		- properties: entire section of additional descriptive properties
//		- ...
//
//	Example:
//		contextInfo := NewContextInfo();
//		contextInfo.Configure(context.Background(), NewConfigParamsFromTuples(
//			ContextInfoParameterName, "MyMicroservice",
//			ContextInfoParameterDescription, "My first microservice"
//		));
//
//		context.Name;     	// Result: "MyMicroservice"
//		context.ContextId;	// Possible result: "mylaptop"
//		context.StartTime;	// Possible result: 2018-01-01:22:12:23.45Z
type ContextInfo struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ContextId   string            `json:"context_id"`
	StartTime   time.Time         `json:"start_time"`
	Properties  map[string]string `json:"properties"`
}

const (
	ContextInfoNameUnknown              = "unknown"
	ContextInfoParameterName            = "name"
	ContextInfoParameterInfoName        = "info.name"
	ContextInfoParameterDescription     = "description"
	ContextInfoParameterInfoDescription = "info.description"
	ContextInfoSectionNameProperties    = "properties"
)

// NewContextInfo creates a new instance of this context info.
//
//	Returns: *ContextInfo
func NewContextInfo() *ContextInfo {
	c := &ContextInfo{
		Name:       ContextInfoNameUnknown,
		StartTime:  time.Now(),
		Properties: map[string]string{},
	}
	c.ContextId, _ = os.Hostname()
	return c
}

// Uptime calculates the context uptime as from the start time.
//
//	Returns: int64 number of milliseconds from the context start time.
func (c *ContextInfo) Uptime() int64 {
	return time.Now().Unix() - c.StartTime.Unix()
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- config *config.ConfigParams configuration parameters to be set.
func (c *ContextInfo) Configure(ctx context.Context, cfg *config.ConfigParams) {
	c.Name = cfg.GetAsStringWithDefault(ContextInfoParameterName, c.Name)
	c.Name = cfg.GetAsStringWithDefault(ContextInfoParameterInfoName, c.Name)

	c.Description = cfg.GetAsStringWithDefault(ContextInfoParameterDescription, c.Description)
	c.Description = cfg.GetAsStringWithDefault(ContextInfoParameterInfoDescription, c.Description)

	if p, ok := cfg.GetSection(ContextInfoSectionNameProperties).InnerValue().(map[string]string); ok {
		c.Properties = p
	}
}

// NewContextInfoFromConfig creates a new instance of this context info.
//
//	Parameters:
//		- cfg *config.ConfigParams a context configuration parameters.
//	Returns: *ContextInfo
func NewContextInfoFromConfig(ctx context.Context, cfg *config.ConfigParams) *ContextInfo {
	result := NewContextInfo()
	result.Configure(ctx, cfg)
	return result
}
