package connect

import (
	"context"
	"sync"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Helper class to retrieve AWS connection and credential parameters,
// validate them and compose a AwsConnectionParams value.
//
// Configuration parameters
//
//   - connections:
//   - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//   - region:                      (optional) AWS region
//   - partition:                   (optional) AWS partition
//   - service:                     (optional) AWS service
//   - resource_type:               (optional) AWS resource type
//   - resource:                    (optional) AWS resource id
//   - arn:                         (optional) AWS resource ARN
//   - credentials:
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//   - access_id:                   AWS access/client id
//   - access_key:                  AWS access/client id
//
// References
//
//   - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connections
//
//   - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
//     See ConnectionParams (in the Pip.Services components package)
//     See IDiscovery (in the Pip.Services components package)
//
//     Example:
//
//     config := NewConfigParamsFromTuples(
//     "connection.region", "us-east1",
//     "connection.service", "s3",
//     "connection.bucket", "mybucket",
//     "credential.access_id", "XXXXXXXXXX",
//     "credential.access_key", "XXXXXXXXXX"
//     );
//
//     connectionResolver := NewAwsConnectionResolver()
//     connectionResolver.Configure(context.Background(), config)
//     connectionResolver.SetReferences(context.Background(), references)
//
//     err, connection :=connectionResolver.Resolve("123")
//     // Now use connection...
type AwsConnectionResolver struct {

	//The connection resolver.
	connectionResolver *ccon.ConnectionResolver

	//The credential resolver.
	credentialResolver * // `cauth` is a package that provides components for authentication and
	// authorization. It includes interfaces and implementations for managing user
	// identities, roles, and permissions, as well as for authenticating and
	// authorizing requests. In this code, `cauth.CredentialResolver` is used to
	// resolve AWS credentials.
	cauth.CredentialResolver
}

func NewAwsConnectionResolver() *AwsConnectionResolver {
	return &AwsConnectionResolver{
		connectionResolver: ccon.NewEmptyConnectionResolver(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
	}
}

// Configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *AwsConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)
}

// Sets references to dependent components.
//   - references 	references to locate the component dependencies.
func (c *AwsConnectionResolver) SetReferences(ctx context.Context, references cref.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.credentialResolver.SetReferences(ctx, references)
}

// Resolves connection and credental parameters and generates a single
// AWSConnectionParams value.
//   - ctx context.Context     execution context to trace execution through call chain.
//
// Returns: AwsConnectionParams connection aprmeters and error.
// See IDiscovery (in the Pip.Services components package)
func (c *AwsConnectionResolver) Resolve(ctx context.Context) (connection *AwsConnectionParams, err error) {
	connection = NewEmptyAwsConnectionParams()
	//var credential *cauth.CredentialParams
	var globalErr error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		data, connErr := c.connectionResolver.Resolve(ctx)
		if connErr == nil && data != nil {
			connection.Append(data.Value())
		}
		globalErr = connErr
	}()
	wg.Wait()

	if globalErr != nil {
		return nil, globalErr
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, credErr := c.credentialResolver.Lookup(ctx)
		if credErr == nil && data != nil {
			connection.Append(data.Value())
		}
		globalErr = credErr
	}()
	wg.Wait()

	if globalErr != nil {
		return nil, globalErr
	}
	// Force ARN parsing
	connection.SetArn(connection.GetArn())
	// Perform validation
	validErr := connection.Validate(ctx)

	if validErr != nil {
		return nil, validErr
	}
	return connection, nil
}
