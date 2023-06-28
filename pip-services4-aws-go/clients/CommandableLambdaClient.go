package clients

import (
	"context"

	"github.com/aws/aws-sdk-go/service/lambda"
)

// Abstract client that calls commandable AWS Lambda Functions.
//
// Commandable services are generated automatically for ICommandable objects.
// Each command is exposed as action determined by "cmd" parameter.
//
// # Configuration parameters
//
// - connections:
//   - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//   - region:                      (optional) AWS region
//
// - credentials:
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//   - access_id:                   AWS access/client id
//   - access_key:                  AWS access/client id
//
// - options:
//   - connect_timeout:             (optional) connection timeout in milliseconds (default: 10 sec)
//
// # References
//
// - \*:logger:\*:\*:1.0            (optional) ILogger components to pass log messages
// - \*:counters:\*:\*:1.0          (optional) ICounters components to pass collected measurements
// - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
// - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
// # See LambdaFunction
//
// Example:
//
//	type MyLambdaClient struct {
//	    *CommandableLambdaClient
//	}
//
//	...
//
//	func (c* MyLambdaClient) GetData(ctx context.Context, id string)(result DataPage[MyData], err error) {
//
//	    valVal, err := c.callCommand(ctx,
//	          "get_data",
//	          map[string]any{ "id": id })
//
//	    if calErr != nil {
//		    return nil, calErr
//		}
//
//		defer timing.EndTiming(ctx, err)
//
//		return awsclient.HandleLambdaResponse[cdata.DataPage[MyData]](calValue)
//	}
//
//	...
//
// client := NewMyLambdaClient();
// client.Configure(context.Background(), NewConfigParamsFromTuples(
//
//	"connection.region", "us-east-1",
//	"connection.access_id", "XXXXXXXXXXX",
//	"connection.access_key", "XXXXXXXXXXX",
//	"connection.arn", "YYYYYYYYYYYYY"
//
// ));
//
// res, err := client.GetData(context.Background(), "123", "1")
// ...
type CommandableLambdaClient struct {
	*LambdaClient
	name string
}

// Creates a new instance of this client.
//   - name a service name.
func NewCommandableLambdaClient(name string) *CommandableLambdaClient {
	c := &CommandableLambdaClient{
		LambdaClient: NewLambdaClient(),
	}
	c.name = name
	return c
}

// Calls a remote action in AWS Lambda function.
// The name of the action is added as "cmd" parameter
// to the action parameters.
//   - ctx context.Context execution context to trace execution through call chain.
//   - cmd               an action name
//   - params            command parameters.
//   - Return           *lambda.InvokeOutput result or error.
func (c *CommandableLambdaClient) CallCommand(ctx context.Context, cmd string, params map[string]any) (result *lambda.InvokeOutput, err error) {
	timing := c.Instrument(ctx, c.name+"."+cmd)
	callRes, callErr := c.Call(ctx, cmd, params)
	timing.EndTiming(ctx, callErr)
	return callRes, callErr
}
