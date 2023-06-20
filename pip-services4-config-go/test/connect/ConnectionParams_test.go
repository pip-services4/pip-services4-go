package test_connect

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestGetAndSetDiscoveryKey(t *testing.T) {
	connection := connect.NewEmptyConnectionParams()
	assert.Equal(t, "", connection.DiscoveryKey())

	connection.SetDiscoveryKey("Discovery key")
	assert.Equal(t, "Discovery key", connection.DiscoveryKey())
	assert.True(t, connection.UseDiscovery())
}

func TestGetAndSetProtocol(t *testing.T) {
	connection := connect.NewEmptyConnectionParams()
	assert.Equal(t, "", connection.Protocol())
	assert.Equal(t, "https", connection.ProtocolWithDefault("https"))

	connection.SetProtocol("smtp")
	assert.Equal(t, "smtp", connection.Protocol())
}

func TestGetAndSetHost(t *testing.T) {
	connection := connect.NewEmptyConnectionParams()
	assert.Equal(t, "", connection.Host())

	connection.SetHost("localhost")
	assert.Equal(t, "localhost", connection.Host())
}

func TestGetAndSetPort(t *testing.T) {
	connection := connect.NewEmptyConnectionParams()
	assert.Equal(t, 0, connection.Port())
	assert.Equal(t, 8080, connection.PortWithDefault(8080))

	connection.SetPort(80)
	assert.Equal(t, 80, connection.Port())
}

func TestGetAndSetUri(t *testing.T) {
	connection := connect.NewEmptyConnectionParams()
	assert.Equal(t, "", connection.Uri())

	connection.SetUri("https://pipgoals:3000")
	assert.Equal(t, "https://pipgoals:3000", connection.Uri())
}
