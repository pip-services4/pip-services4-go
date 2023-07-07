package clients

import (
	"context"
	"net/http"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// Abstract client that calls commandable Google Cloud Functions.
//
// Commandable services are generated automatically for ICommandable objects.
// Each command is exposed as action determined by "cmd" parameter.
//
//	Configuration parameters
//		- connections:
//			- uri:           full connection uri with specific app and function name
//			- protocol:      connection protocol
//			- project_id:    is your Google Cloud Platform project ID
//			- region:        is the region where your function is deployed
//			- function:      is the name of the HTTP function you deployed
//			- org_id:        organization name
//		- options:
//			- retries:               number of retries (default: 3)
//			- connect_timeout:       connection timeout in milliseconds (default: 10 sec)
//			- timeout:               invocation timeout in milliseconds (default: 10 sec)
//		- credentials:
//			- account: the service account name
//			- auth_token:    Google-generated ID token or null if using custom auth (IAM)
//
//	References
//		- *:logger:*:*:1.0				(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0			(optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0			(optional) IDiscovery services to resolve connection
//		- *:credential-store:*:*:1.0	(optional) Credential stores to resolve credentials
//
// see CloudFunction
//
//	Exammple:
//		type MyCommandableGoogleClient struct {
//			clients.CommandableCloudFunctionClient
//		}
//
//		func NewMyCommandableGoogleClient() *MyCommandableGoogleClient {
//			return &MyCommandableGoogleClient{
//				CommandableCloudFunctionClient: *gcpclient.NewCommandableCloudFunctionClient(),
//			}
//		}
//
//		func (c *MyCommandableGoogleClient) GetData(ctx context.Context, id string) MyData {
//			response, err := c.CallCommand(ctx, "dummies.get_dummies", cdata.NewAnyValueMapFromTuples("id", id))
//			if err != nil {
//				return MyData{}, err
//			}
//
//			return rpcclient.HandleHttpResponse[MyData](response, cctx.GetTraceId(ctx))
//		}
//
//		...
//		client := NewMyCommandableGoogleClient()
//		client.Configure(config.NewConfigParamsFromTuples(
//			"connection.uri", "http://region-id.cloudfunctions.net/myfunction",
//			"connection.protocol", "http",
//			"connection.region", "region",
//			"connection.function", "myfunction",
//			"connection.project_id", "id",
//			"credential.auth_token", "XXX",
//		))
//		result := client.GetData("123", "1")
//		...
type CommandableCloudFunctionClient struct {
	*CloudFunctionClient
	name string
}

// Creates a new instance of this client.
// Parameters:
//   - name	a service name.
func NewCommandableCloudFunctionClient(name string) *CommandableCloudFunctionClient {
	return &CommandableCloudFunctionClient{name: name, CloudFunctionClient: NewCloudFunctionClient()}
}

// Calls a remote action in Google Function.
// The name of the action is added as "cmd" parameter
// to the action parameters.
// Parameters:
//   - ctx context.Context a context to trace execution through call chain.
//   - cmd	an action name
//   - params	command parameters.
//
// Returns action result.
func (c *CommandableCloudFunctionClient) CallCommand(ctx context.Context, cmd string, params *cdata.AnyValueMap) (*http.Response, error) {
	timing := c.Instrument(ctx, c.name+"."+cmd)
	r, err := c.Call(ctx, cmd, params)
	timing.EndTiming(ctx, err)
	return r, err
}
