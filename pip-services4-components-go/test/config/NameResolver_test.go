package test_config

import (
	"testing"

	conf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

func TestResolveName(t *testing.T) {
	var config = conf.NewConfigParamsFromTuples("id", "ABC")
	var name = conf.NameResolver.Resolve(config)
	assert.Equal(t, "ABC", name)

	config = conf.NewConfigParamsFromTuples("name", "ABC")
	name = conf.NameResolver.Resolve(config)
	assert.Equal(t, "ABC", name)
}

func TestResolveEmptyName(t *testing.T) {
	var config = conf.NewConfigParamsFromTuples()
	var name = conf.NameResolver.Resolve(config)
	assert.Equal(t, "", name)
}
