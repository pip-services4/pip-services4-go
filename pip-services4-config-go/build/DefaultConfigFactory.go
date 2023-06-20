package build

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/config"
)

// MemoryCredentialStoreDescriptor Creates ICredentialStore components by their descriptors.
var MemoryCredentialStoreDescriptor = refer.NewDescriptor("pip-services", "credential-store", "memory", "*", "1.0")
var MemoryConfigReaderDescriptor = refer.NewDescriptor("pip-services", "config-reader", "memory", "*", "1.0")
var JsonConfigReaderDescriptor = refer.NewDescriptor("pip-services", "config-reader", "json", "*", "1.0")
var YamlConfigReaderDescriptor = refer.NewDescriptor("pip-services", "config-reader", "yaml", "*", "1.0")
var MemoryDiscoveryDescriptor = refer.NewDescriptor("pip-services", "discovery", "memory", "*", "1.0")

// NewDefaultConfigFactory create a new instance of the factory.
//
//	Returns: *build.Factory
func NewDefaultConfigFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(MemoryCredentialStoreDescriptor, auth.NewEmptyMemoryCredentialStore)
	factory.RegisterType(MemoryConfigReaderDescriptor, config.NewEmptyMemoryConfigReader)
	factory.RegisterType(JsonConfigReaderDescriptor, config.NewJsonConfigReader)
	factory.RegisterType(YamlConfigReaderDescriptor, config.NewEmptyYamlConfigReader)
	factory.RegisterType(MemoryDiscoveryDescriptor, config.NewEmptyMemoryConfigReader)

	return factory
}
