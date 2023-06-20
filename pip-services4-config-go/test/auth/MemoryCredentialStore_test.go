package test_auth

import (
	"context"
	"testing"

	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	"github.com/stretchr/testify/assert"
)

func TestLookupAndStore(t *testing.T) {
	config := cconfig.NewConfigParamsFromTuples(
		"key1.user", "user1",
		"key1.pass", "pass1",
		"key2.user", "user2",
		"key2.pass", "pass2",
	)

	credentialStore := auth.NewEmptyMemoryCredentialStore()
	credentialStore.ReadCredentials(config)

	cred1, _ := credentialStore.Lookup(context.Background(), "key1")
	cred2, _ := credentialStore.Lookup(context.Background(), "key2")

	assert.Equal(t, "user1", cred1.Username())
	assert.Equal(t, "pass1", cred1.Password())
	assert.Equal(t, "user2", cred2.Username())
	assert.Equal(t, "pass2", cred2.Password())

	credConfig := auth.NewCredentialParamsFromTuples(
		"user", "user3",
		"pass", "pass3",
		"access_id", "12345",
	)

	credentialStore.Store(context.Background(), "key3", credConfig)

	cred3, _ := credentialStore.Lookup(context.Background(), "key3")

	assert.Equal(t, "user3", cred3.Username())
	assert.Equal(t, "pass3", cred3.Password())
	assert.Equal(t, "12345", cred3.AccessId())
}
