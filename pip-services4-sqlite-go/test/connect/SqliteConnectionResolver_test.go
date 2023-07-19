package test_connect

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestSqliteConnectionResolverConnectionConfigWithParams(t *testing.T) {

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.database", "../../data/test.db",
	)

	resolver := conn.NewSqliteConnectionResolver()
	resolver.Configure(context.Background(), dbConfig)

	config, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)

	assert.NotNil(t, config)
	assert.Equal(t, "../../data/test.db", config)
}
