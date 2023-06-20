package config

import (
	"context"
	"path/filepath"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-config-go/config"
)

// ContainerConfigReader Helper class that reads container configuration from JSON or YAML file.
var ContainerConfigReader = &_TContainerConfigReader{}

type _TContainerConfigReader struct{}

// ReadFromFile reads container configuration from JSON or YAML file.
// The type of the file is determined by file extension.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- path string a path to component configuration file.
//		- parameters *config.ConfigParams values to parameters the configuration or null to skip parameterization.
//	Returns: ContainerConfig, error the read container configuration and error
func (c *_TContainerConfigReader) ReadFromFile(ctx context.Context,
	path string, parameters *config.ConfigParams) (ContainerConfig, error) {
	traceId := utils.ContextHelper.GetTraceId(ctx)
	if path == "" {
		return nil, errors.NewConfigError(traceId, "NO_PATH", "Missing config file path")
	}

	ext := filepath.Ext(path)

	if ext == ".json" {
		return c.ReadFromJsonFile(ctx, path, parameters)
	}

	if ext == ".yaml" || ext == ".yml" {
		return c.ReadFromYamlFile(ctx, path, parameters)
	}

	return c.ReadFromJsonFile(ctx, path, parameters)
}

// ReadFromJsonFile reads container configuration from JSON file.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- path string a path to component configuration file.
//		- parameters *config.ConfigParams values to parameters the configuration or null to skip parameterization.
//	Returns: ContainerConfig, error the read container configuration and error
func (c *_TContainerConfigReader) ReadFromJsonFile(ctx context.Context,
	path string, parameters *config.ConfigParams) (ContainerConfig, error) {

	config, err := cconfig.ReadJsonConfig(ctx, path, parameters)
	if err != nil {
		return nil, err
	}
	return ReadContainerConfigFromConfig(config)
}

// ReadFromYamlFile reads container configuration from YAML file.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- path string a path to component configuration file.
//		- parameters *config.ConfigParams values to parameters the configuration or null to skip parameterization.
//	Returns: ContainerConfig, error the read container configuration and error
func (c *_TContainerConfigReader) ReadFromYamlFile(ctx context.Context,
	path string, parameters *config.ConfigParams) (ContainerConfig, error) {

	config, err := cconfig.ReadYamlConfig(ctx, path, parameters)
	if err != nil {
		return nil, err
	}
	return ReadContainerConfigFromConfig(config)
}
