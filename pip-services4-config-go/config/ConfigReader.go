package config

import (
	"context"

	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache"
)

// ConfigReader abstract config reader that supports configuration parameterization.
//
//	Configuration parameters:
//		parameters this entire section is used as template parameters
type ConfigReader struct {
	parameters *cconfig.ConfigParams
}

// SectionNameParameters is a name of ConfigReader section
const SectionNameParameters = "parameters"

// NewConfigReader creates a new instance of the config reader.
//
//	Returns: *ConfigReader
func NewConfigReader() *ConfigReader {
	return &ConfigReader{
		parameters: cconfig.NewEmptyConfigParams(),
	}
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *ConfigReader) Configure(ctx context.Context, config *cconfig.ConfigParams) {
	parameters := config.GetSection(SectionNameParameters)
	if parameters.Len() > 0 {
		c.parameters = parameters
	}
}

// Parameterize configuration template given as string with dynamic parameters.
// The method uses Mustache template engine implemented in expressions module
//
//	Parameters:
//		- config string a string with configuration template to be parameterized
//		- parameters *config.ConfigParams dynamic parameters to inject into the template
//	Returns: string, error a parameterized configuration string abd error.
func (c *ConfigReader) Parameterize(config string, parameters *cconfig.ConfigParams) (string, error) {
	if parameters == nil {
		parameters = cconfig.NewEmptyConfigParams()
	}

	parameters = c.parameters.Override(parameters)

	value := parameters.Value()

	mustacheTemplate, err := mustache.NewMustacheTemplateFromString(config)
	if err != nil {
		return "", err
	}

	result, err := mustacheTemplate.EvaluateWithVariables(value)
	return result, err
}

// AddChangeListener - Adds a listener that will be notified when configuration is changed
func (c *ConfigReader) AddChangeListener(ctx context.Context, listener cexec.INotifiable) {
	// Do nothing...
}

// RemoveChangeListener - Remove a previously added change listener.
func (c *ConfigReader) RemoveChangeListener(ctx context.Context, listener cexec.INotifiable) {
	// Do nothing...
}
