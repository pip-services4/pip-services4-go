package auth

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// CredentialResolver helper class to retrieve component credentials.
// If credentials are configured to be retrieved from ICredentialStore,
// it automatically locates ICredentialStore in component references and retrieve
// credentials from there using store_key parameter.
//
//	Configuration parameters:
//		credential:
//			- store_key: (optional) a key to retrieve the credentials from ICredentialStore
//				- ... other credential parameters
//		credentials: alternative to credential
//			- [credential params 1]: first credential parameters
//				- ... credential parameters for key 1
//				- ...
//			- [credential params N]: Nth credential parameters
//				- ... credential parameters for key N
//	References:
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
//
// see CredentialParams
// see ICredentialStore
//
//	Example:
//		config := NewConfigParamsFromTuples(
//			"credential.user", "jdoe",
//			"credential.pass",  "pass123"
//		);
//
//		credentialResolver := NewCredentialResolver();
//		credentialResolver.Configure(context.Background(), config);
//		credentialResolver.SetReferences(context.Background(), references);
//
//		cred, err := credentialResolver.Lookup(context.Background(), "123");
//		// Now use credential...
type CredentialResolver struct {
	credentials []*CredentialParams
	references  refer.IReferences
}

// NewEmptyCredentialResolver creates a new instance of credentials resolver.
//
//	Returns: *CredentialResolver
func NewEmptyCredentialResolver() *CredentialResolver {
	return &CredentialResolver{
		credentials: make([]*CredentialParams, 0),
		references:  nil,
	}
}

// NewCredentialResolver creates a new instance of credentials resolver.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams component configuration parameters
//		- references refer.IReferences component references
//	Returns: *CredentialResolver
func NewCredentialResolver(ctx context.Context, config *config.ConfigParams, references refer.IReferences) *CredentialResolver {
	c := &CredentialResolver{
		credentials: make([]*CredentialParams, 0),
		references:  references,
	}

	if config != nil {
		c.Configure(ctx, config)
	}

	return c
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- config *config.ConfigParams configuration parameters to be set.
func (c *CredentialResolver) Configure(ctx context.Context, config *config.ConfigParams) {
	credentials := NewManyCredentialParamsFromConfig(config)
	c.credentials = append(c.credentials, credentials...)
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- references refer.IReferences references to locate the component dependencies.
func (c *CredentialResolver) SetReferences(ctx context.Context, references refer.IReferences) {
	c.references = references
}

// GetAll gets all credentials configured in component configuration.
// Redirect to CredentialStores is not done at this point.
// If you need fully fleshed credential use lookup method instead.
//
//	Returns: []*CredentialParams a list with credential parameters
func (c *CredentialResolver) GetAll() []*CredentialParams {
	return c.credentials
}

// Add a new credential to component credentials
//
//	Parameters:
//		- credential *CredentialParams new credential parameters to be added
func (c *CredentialResolver) Add(credential *CredentialParams) {
	c.credentials = append(c.credentials, credential)
}

func (c *CredentialResolver) lookupInStores(ctx context.Context,
	credential *CredentialParams) (result *CredentialParams, err error) {

	if !credential.UseCredentialStore() {
		return credential, nil
	}

	key := credential.StoreKey()
	storeDescriptor := refer.NewDescriptor("*", "credential_store", "*", "*", "*")
	if c.references == nil {
		return nil, refer.NewReferenceError(ctx, storeDescriptor)
	}

	components := c.references.GetOptional(storeDescriptor)
	if len(components) == 0 {
		err := refer.NewReferenceError(ctx, storeDescriptor)
		return nil, err
	}

	for _, component := range components {
		if store, ok := component.(ICredentialStore); ok && store != nil {
			credential, err = store.Lookup(ctx, key)
			if credential != nil || err != nil {
				return credential, err
			}
		}
	}

	return nil, errors.NewConfigError(
		cctx.GetTraceId(ctx), "MISSING_CREDENTIALS", "missing credential param")
}

// Lookup component credential parameters.
// If credentials are configured to be retrieved from Credential store it finds a
// ICredentialStore and lookups credentials there.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: *CredentialParams? error
func (c *CredentialResolver) Lookup(ctx context.Context) (*CredentialParams, error) {
	if len(c.credentials) == 0 {
		return nil, nil
	}

	lookupCredentials := make([]*CredentialParams, 0)

	for _, credential := range c.credentials {
		if !credential.UseCredentialStore() {
			return credential, nil
		}

		lookupCredentials = append(lookupCredentials, credential)
	}

	for _, credential := range lookupCredentials {
		_c, err := c.lookupInStores(ctx, credential)
		if _c != nil || err != nil {
			return _c, err
		}
	}

	return nil, errors.NewConfigError(
		cctx.GetTraceId(ctx), "MISSING_CREDENTIALS", "missing credential param")
}
