package connect_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	gcpconn "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestEmptyConnection(t *testing.T) {
	connection := gcpconn.NewEmptyGcpConnectionParams()

	_, ok := connection.Uri()
	assert.False(t, ok)
	_, ok = connection.ProjectId()
	assert.False(t, ok)
	_, ok = connection.Function()
	assert.False(t, ok)
	_, ok = connection.Region()
	assert.False(t, ok)
	_, ok = connection.Protocol()
	assert.False(t, ok)
	_, ok = connection.AuthToken()
	assert.False(t, ok)
}

func TestComposeConfig(t *testing.T) {
	ctx := context.Background()

	config1 := cconf.NewConfigParamsFromTuples(
		"connection.uri", "http://east-my_test_project.cloudfunctions.net/myfunction",
		"credential.auth_token", "1234",
	)

	config2 := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.region", "east",
		"connection.function", "myfunction",
		"connection.project_id", "my_test_project",
		"credential.auth_token", "1234",
	)

	resolver := gcpconn.NewGcpConnectionResolver()

	resolver.Configure(ctx, config1)
	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)

	uri, _ := connection.Uri()
	region, _ := connection.Region()
	protocol, _ := connection.Protocol()
	function, _ := connection.Function()
	projectId, _ := connection.ProjectId()
	authToken, _ := connection.AuthToken()

	assert.Equal(t, "http://east-my_test_project.cloudfunctions.net/myfunction", uri)
	assert.Equal(t, "east", region)
	assert.Equal(t, "http", protocol)
	assert.Equal(t, "myfunction", function)
	assert.Equal(t, "my_test_project", projectId)
	assert.Equal(t, "1234", authToken)

	resolver = gcpconn.NewGcpConnectionResolver()
	resolver.Configure(ctx, config2)
	connection, err = resolver.Resolve(context.Background())
	assert.Nil(t, err)

	uri, _ = connection.Uri()
	region, _ = connection.Region()
	protocol, _ = connection.Protocol()
	function, _ = connection.Function()
	projectId, _ = connection.ProjectId()
	authToken, _ = connection.AuthToken()

	assert.Equal(t, "http://east-my_test_project.cloudfunctions.net/myfunction", uri)
	assert.Equal(t, "east", region)
	assert.Equal(t, "http", protocol)
	assert.Equal(t, "myfunction", function)
	assert.Equal(t, "my_test_project", projectId)
	assert.Equal(t, "1234", authToken)
}
