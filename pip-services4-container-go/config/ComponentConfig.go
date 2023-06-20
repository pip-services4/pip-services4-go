package config

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// ComponentConfig configuration of a component inside a container.
// The configuration includes type information or descriptor, and component configuration parameters.
type ComponentConfig struct {
	Descriptor *refer.Descriptor
	Type       *reflect.TypeDescriptor
	Config     *config.ConfigParams
}

// NewComponentConfigFromDescriptor creates a new instance of the component configuration.
//
//	Parameters:
//		- descriptor *refer.Descriptor a components descriptor (locator).
//		- config *config.ConfigParams component configuration parameters.
//	Returns: *ComponentConfig
func NewComponentConfigFromDescriptor(descriptor *refer.Descriptor,
	config *config.ConfigParams) *ComponentConfig {
	return &ComponentConfig{
		Descriptor: descriptor,
		Config:     config,
	}
}

// NewComponentConfigFromType creates a new instance of the component configuration.
//
//	Parameters:
//		- typ *reflect.TypeDescriptor a components type descriptor.
//		- config *config.ConfigParams component configuration parameters.
//	Returns: *ComponentConfig
func NewComponentConfigFromType(typ *reflect.TypeDescriptor,
	config *config.ConfigParams) *ComponentConfig {
	return &ComponentConfig{
		Type:   typ,
		Config: config,
	}
}

// ReadComponentConfigFromConfig creates a new instance of ComponentConfig
// based on section from container configuration.
//
//	Parameters: config *config.ConfigParams component parameters from container configuration
//	Returns: *ComponentConfig, error a newly created ComponentConfig and
//		ConfigError when neither component descriptor or type is found.
func ReadComponentConfigFromConfig(config *config.ConfigParams) (result *ComponentConfig, err error) {
	descriptor, err1 := refer.ParseDescriptorFromString(config.GetAsString("descriptor"))
	if err1 != nil {
		return nil, err1
	}

	typ, err2 := reflect.ParseTypeDescriptorFromString(config.GetAsString("type"))
	if err2 != nil {
		return nil, err2
	}

	if descriptor == nil && typ == nil {
		err = errors.NewConfigError(
			"",
			"BAD_CONFIG",
			"Component configuration must have descriptor or type",
		)
		return nil, err
	}

	return &ComponentConfig{
		Descriptor: descriptor,
		Type:       typ,
		Config:     config,
	}, nil
}
