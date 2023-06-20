package test_config

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	conf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-container-go/config"

	"github.com/stretchr/testify/assert"
)

func TestComponentConfigType(t *testing.T) {
	typ := reflect.NewTypeDescriptor("new name", "")
	config := conf.NewConfigParamsFromTuples(
		"config.key", "key",
		"config.key2", "key2",
	)
	componentConfig := cconf.NewComponentConfigFromType(typ, config)

	assert.NotNil(t, componentConfig.Type)
	assert.Nil(t, componentConfig.Descriptor)
	assert.NotNil(t, componentConfig.Config)
}

func TestComponentConfigDescriptor(t *testing.T) {
	descriptor := refer.NewDescriptor("group", "type", "id", "default", "version")
	config := conf.NewConfigParamsFromTuples(
		"config.key", "key",
		"config.key2", "key2",
	)
	componentConfig := cconf.NewComponentConfigFromDescriptor(descriptor, config)

	assert.Nil(t, componentConfig.Type)
	assert.NotNil(t, componentConfig.Descriptor)
	assert.NotNil(t, componentConfig.Config)
}

func TestComponentConfigFromEmptyConfig(t *testing.T) {
	config := conf.NewEmptyConfigParams()
	componentConfig, err := cconf.ReadComponentConfigFromConfig(config)

	assert.NotNil(t, err)
	assert.Nil(t, componentConfig)
}

func TestComponentConfigFromWrongConfig(t *testing.T) {
	config := conf.NewConfigParamsFromTuples(
		"descriptor", "descriptor_name",
		"type", "type",
		"config.key", "key",
		"config.key2", "key2",
	)
	componentConfig, err := cconf.ReadComponentConfigFromConfig(config)

	assert.NotNil(t, err)
	assert.Nil(t, componentConfig)
}

func TestComponentConfigFromCorrectConfig(t *testing.T) {
	config := conf.NewConfigParamsFromTuples(
		"descriptor", "group:type:kind:name:version",
		"type", "type",
		"config.key", "key",
		"config.key2", "key2",
	)
	componentConfig, err := cconf.ReadComponentConfigFromConfig(config)

	assert.Nil(t, err)
	assert.NotNil(t, componentConfig)
	assert.NotNil(t, componentConfig.Descriptor)
	assert.NotNil(t, componentConfig.Type)

	assert.Equal(t, "group", componentConfig.Descriptor.Group())
	assert.Equal(t, "type", componentConfig.Descriptor.Type())
	assert.Equal(t, "kind", componentConfig.Descriptor.Kind())
	assert.Equal(t, "name", componentConfig.Descriptor.Name())
	assert.Equal(t, "version", componentConfig.Descriptor.Version())
}
