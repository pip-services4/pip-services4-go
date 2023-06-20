package config

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"sort"
)

// ContainerConfig Container configuration defined as a list of component configurations.
type ContainerConfig []*ComponentConfig

// NewContainerConfig creates a new instance of container configuration.
//
//	Parameters: components *ComponentConfig a list of component configurations.
//	Returns: ContainerConfig
func NewContainerConfig(components ...*ComponentConfig) ContainerConfig {
	return components
}

// NewContainerConfigFromValue creates a new ContainerConfig object filled with key-value
// pairs from specified object. The value is converted into
// ConfigParams object which is used to create the object.
//
//	see FromConfig
//	Parameters: value any an object with key-value pairs used to initialize a new ContainerConfig.
//	Returns: ContainerConfig a new ContainerConfig object.
func NewContainerConfigFromValue(value any) ContainerConfig {
	config := config.NewConfigParamsFromValue(value)
	result, _ := ReadContainerConfigFromConfig(config)
	return result
}

// ReadContainerConfigFromConfig creates a new ContainerConfig object based on configuration parameters.
// Each section in the configuration parameters is converted into a component configuration.
//
//	Parameters: config *config.ConfigParams
//	Returns: ContainerConfig, error a new ContainerConfig object and error
func ReadContainerConfigFromConfig(config *config.ConfigParams) (ContainerConfig, error) {
	if config == nil {
		return []*ComponentConfig{}, nil
	}

	names := config.GetSectionNames()
	// Sort so components should come in a right order
	sort.Strings(names)
	result := make([]*ComponentConfig, len(names))
	for i, v := range names {
		c := config.GetSection(v)
		componentConfig, err := ReadComponentConfigFromConfig(c)
		if err != nil {
			return nil, err
		}
		result[i] = componentConfig
	}

	return result, nil
}
