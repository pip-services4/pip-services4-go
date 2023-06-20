package test_auth

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	"github.com/stretchr/testify/assert"
)

func TestCredentialResolverConfigure(t *testing.T) {
	restConfig := config.NewConfigParamsFromTuples(
		"credential.username", "Negrienko",
		"credential.password", "qwerty",
		"credential.access_key", "key",
		"credential.store_key", "store key",
	)
	credentialResolver := auth.NewCredentialResolver(context.Background(), restConfig, nil)
	credentials := credentialResolver.GetAll()
	assert.Len(t, credentials, 1)

	credential := credentials[0]
	assert.Equal(t, "Negrienko", credential.Username())
	assert.Equal(t, "qwerty", credential.Password())
	assert.Equal(t, "key", credential.AccessKey())
	assert.Equal(t, "store key", credential.StoreKey())
}

func TestCredentialResolverLookup(t *testing.T) {
	credentialResolver := auth.NewEmptyCredentialResolver()

	credential, err := credentialResolver.Lookup(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, credential)

	restConfigWithoutStoreKey := config.NewConfigParamsFromTuples(
		"credential.username", "Negrienko",
		"credential.password", "qwerty",
		"credential.access_key", "key",
	)
	credentialResolver = auth.NewCredentialResolver(context.Background(), restConfigWithoutStoreKey, nil)

	credential, err = credentialResolver.Lookup(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, "Negrienko", credential.Username())
	assert.Equal(t, "qwerty", credential.Password())
	assert.Equal(t, "key", credential.AccessKey())
	assert.Equal(t, "", credential.StoreKey())

	restConfig := config.NewConfigParamsFromTuples(
		"credential.username", "Negrienko",
		"credential.password", "qwerty",
		"credential.access_key", "key",
		"credential.store_key", "store key",
	)
	credentialResolver = auth.NewCredentialResolver(context.Background(), restConfig, nil)

	credential, err = credentialResolver.Lookup(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, credential)
}
