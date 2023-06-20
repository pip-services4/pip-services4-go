package config

import (
	"context"

	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache"
)

// MemoryConfigReader is a config reader that stores configuration in memory.
// The reader supports parameterization using Mustache template engine implemented in expressions module.
// Configuration parameters: The configuration parameters are the configuration template
//
//	see IConfigReader
//	Example
//		config := NewConfigParamsFromTuples(
//			"connection.host", "{{SERVICE_HOST}}",
//			"connection.port", "{{SERVICE_PORT}}{{^SERVICE_PORT}}8080{{/SERVICE_PORT}}"
//		);
//		configReader := NewMemoryConfigReader();
//		configReader.Configure(context.Background(), config);
//		parameters := NewConfigParamsFromValue(process.env);
//		res, err := configReader.ReadConfig(context.Background(), "123", parameters);
//			Possible result: connection.host=10.1.1.100;connection.port=8080
type MemoryConfigReader struct {
	config *cconfig.ConfigParams
}

// NewEmptyMemoryConfigReader creates a new instance of config reader.
//
//	Returns: *MemoryConfigReader
func NewEmptyMemoryConfigReader() *MemoryConfigReader {
	return &MemoryConfigReader{
		config: cconfig.NewEmptyConfigParams(),
	}
}

// NewMemoryConfigReader creates a new instance of config reader.
//
//	Parameters: config *cconfig.ConfigParams component configuration parameters
//	Returns: *MemoryConfigReader
func NewMemoryConfigReader(config *cconfig.ConfigParams) *MemoryConfigReader {
	return &MemoryConfigReader{
		config: config,
	}
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *cconfig.ConfigParams configuration parameters to be set.
func (c *MemoryConfigReader) Configure(ctx context.Context, config *cconfig.ConfigParams) {
	c.config = config
}

// ReadConfig reads configuration and parameterize it with given values.
//
//	Parameters:
//		- ctx context.Context
//		-  transaction id to trace execution through call chain.
//		- parameters *cconfig.ConfigParams values to parameters
//			the configuration or null to skip parameterization.
//	Returns: *cconfig.ConfigParams, error configuration or error.
func (c *MemoryConfigReader) ReadConfig(ctx context.Context,
	parameters *cconfig.ConfigParams) (*cconfig.ConfigParams, error) {

	if parameters != nil {
		template := c.config.String()
		value := parameters.Value()

		mustacheTemplate, err := mustache.NewMustacheTemplateFromString(template)
		if err != nil {
			return nil, err
		}

		config, err := mustacheTemplate.EvaluateWithVariables(value)
		if err != nil {
			return nil, err
		}

		result := cconfig.NewConfigParamsFromString(config)
		return result, nil
	} else {
		result := cconfig.NewConfigParamsFromValue(c.config.Value())
		return result, nil
	}
}

// AddChangeListener - Adds a listener that will be notified when configuration is changed
func (c *MemoryConfigReader) AddChangeListener(ctx context.Context, listener cexec.INotifiable) {
	// Do nothing...
}

// RemoveChangeListener - Remove a previously added change listener.
func (c *MemoryConfigReader) RemoveChangeListener(ctx context.Context, listener cexec.INotifiable) {
	// Do nothing...
}
