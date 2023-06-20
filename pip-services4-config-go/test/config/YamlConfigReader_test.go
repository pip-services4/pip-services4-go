package test_config

import (
	"context"
	"testing"

	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/config"
	"github.com/stretchr/testify/assert"
)

func TestYamlConfigReader(t *testing.T) {
	parameters := cconfig.NewConfigParamsFromTuples(
		"param1", "Test Param 1",
		"param2", "Test Param 2",
	)
	config, err := config.ReadYamlConfig(context.Background(), "./config.yml", parameters)

	assert.Nil(t, err)
	assert.Equal(t, 9, config.Len())
	assert.Equal(t, 123, config.GetAsInteger("field1.field11"))
	assert.Equal(t, "ABC", config.GetAsString("field1.field12"))
	assert.Equal(t, 123, config.GetAsInteger("field2.0"))
	assert.Equal(t, "ABC", config.GetAsString("field2.1"))
	assert.Equal(t, 543, config.GetAsInteger("field2.2.field21"))
	assert.Equal(t, "XYZ", config.GetAsString("field2.2.field22"))
	assert.Equal(t, true, config.GetAsBoolean("field3"))
	assert.Equal(t, "Test Param 1", config.GetAsString("field4"))
	assert.Equal(t, "Test Param 2", config.GetAsString("field5"))
}

func TestYamlConfigReaderReadArray(t *testing.T) {
	config, err := config.ReadYamlConfig(context.Background(), "./dummy.yml", nil)

	assert.Nil(t, err)
	assert.True(t, config.Len() > 0)
}
