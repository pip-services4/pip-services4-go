package clients

import (
	"context"
	"net/http"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// Abstract client that calls commandable Azure Functions.
//
// Commandable services are generated automatically for ICommandable objects.
// Each command is exposed as action determined by "cmd" parameter.
//
//	Configuration parameters
//		- connections:
//		    - uri:                         (optional) full connection string or use protocol, app_name and function_name to build
//		    - protocol:                    (optional) connection protocol
//		    - app_name:                    (optional) Azure Function application name
//		    - function_name:               (optional) Azure Function name
//		- options:
//		     - retries:               number of retries (default: 3)
//		     - connect_timeout:       connection timeout in milliseconds (default: 10 sec)
//		     - timeout:               invocation timeout in milliseconds (default: 10 sec)
//		- credentials:
//		    - auth_code:                   Azure Function auth code if use custom authorization provide empty string
//
//	References
//		- *:logger:*:*:1.0				(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0			(optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0			(optional) IDiscovery services to resolve connection
//		- *:credential-store:*:*:1.0	(optional) Credential stores to resolve credentials
//
// see AzureFunction
//
//	Exammple:
//		type MyCommandableAzureClient struct {
//			*clients.CommandableAzureFunctionClient
//		}
//
//		func NewMyCommandableAzureClient() *MyCommandableAzureClient {
//			return &MyCommandableAzureClient{
//				CommandableAzureFunctionClient: azureclient.NewCommandableAzureFunctionClient(),
//			}
//		}
//
//		func (c *MyCommandableAzureClient) GetData(ctx context.Context, id string) MyData {
//			response, err := c.CallCommand(ctx, "dummies.get_dummies", cdata.NewAnyValueMapFromTuples("id", id))
//			if err != nil {
//				return MyData{}, err
//			}
//
//			return rpcclient.HandleHttpResponse[MyData](response, cctx.GetTraceId(ctx))
//		}
//
//		...
//		client := NewMyCommandableAzureClient()
//		client.Configure(config.NewConfigParamsFromTuples(
//			"connection.uri", "http://myapp.azurewebsites.net/api/myfunction",
//			"connection.protocol", "http",
//			"connection.app_name", "myapp",
//			"connection.function_name", "myfunction"
//			"credential.auth_code", "XXXX"
//		))
//		result := client.GetData("123", "1")
//		...
type CommandableAzureFunctionClient struct {
	*AzureFunctionClient
	name string
}

// Creates a new instance of this client.
// Parameters:
//   - name	a service name.
func NewCommandableAzureFunctionClient(name string) *CommandableAzureFunctionClient {
	return &CommandableAzureFunctionClient{name: name, AzureFunctionClient: NewAzureFunctionClient()}
}

// Calls a remote action in Azure Function.
// The name of the action is added as "cmd" parameter
// to the action parameters.
// Parameters:
//   - ctx context.Context execution context to trace execution through call chain.
//   - cmd	an action name
//   - params	command parameters.
//
// Returns action result.
func (c *CommandableAzureFunctionClient) CallCommand(ctx context.Context, cmd string, params *cdata.AnyValueMap) (*http.Response, error) {
	timing := c.Instrument(ctx, c.name+"."+cmd)
	r, err := c.Call(ctx, cmd, params)
	timing.EndTiming(ctx, err)
	return r, err
}
