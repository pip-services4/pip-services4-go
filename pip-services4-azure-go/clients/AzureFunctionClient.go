package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	azureconn "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/connect"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	ckeys "github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
)

// Abstract client that calls Azure Functions.
//
// When making calls "cmd" parameter determines which what action shall be called, while
// other parameters are passed to the action itself.
//
//	Configuration parameters
//		- connections:
//		    - uri:                         (optional) full connection string or use protocol, app_name and function_name t
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
// see AzureFunction, CommandableAzureClient
//
//	Example:
//		type MyAzureFunctionClient struct {
//			*clients.AzureFunctionClient
//		}
//
//		func NewMyAzureFunctionClient() *MyAzureFunctionClient {
//			return &MyAzureFunctionClient{
//				AzureFunctionClient: azureclient.NewAzureFunctionClient(),
//			}
//		}
//
//		func (c *MyAzureFunctionClient) GetData(ctx context.Context, id string) MyData {
//			timing := c.Instrument(ctx, "myclient.get_data")
//
//			response, err := c.Call(ctx, "get_data", data.NewAnyValueMapFromTuples("id", dummyId))
//
//			defer timing.EndTiming(ctx, err)
//			return rpcclients.HandleHttpResponse[MyData](response, utils.ContextHelper.GetTraceId(ctx))
//		}
//
//		...
//
//		client := NewMyAzureFunctionClient()
//		client.Configure(config.NewConfigParamsFromTuples(
//			"connection.uri", "http://myapp.azurewebsites.net/api/myfunction",
//			"connection.protocol", "http",
//			"connection.app_name", "myapp",
//			"connection.function_name", "myfunction"
//			"credential.auth_code", "XXXX"
//		result := client.GetData("123", "1")
type AzureFunctionClient struct {
	// The HTTP client.
	Client *http.Client
	// The Azure Function connection parameters
	Connection *azureconn.AzureFunctionConnectionParams
	// The number of retries.
	Retries int
	// The default headers to be added to every request.
	Headers *cdata.StringValueMap
	// The connection timeout in milliseconds.
	ConnectTimeout int
	// The invocation timeout in milliseconds.
	Timeout int
	// The remote service uri which is calculated on open.
	Uri string
	// The connection resolver.
	ConnectionResolver *azureconn.AzureFunctionConnectionResolver
	// The dependency resolver.
	DependencyResolver *crefer.DependencyResolver

	// The logger.
	Logger *clog.CompositeLogger
	// The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

const (
	DefaultConnectTimeout = 10000
	DefaultTimeout        = 10000
	DefaultRetriesCount   = 3
)

// Creates new instance of AzureFunctionClient
func NewAzureFunctionClient() *AzureFunctionClient {
	c := AzureFunctionClient{}

	c.ConnectionResolver = azureconn.NewAzureFunctionConnectionResolver()
	c.DependencyResolver = crefer.NewDependencyResolver()
	c.Logger = clog.NewCompositeLogger()
	c.Counters = ccount.NewCompositeCounters()
	c.Tracer = ctrace.NewCompositeTracer()
	c.Headers = cdata.NewEmptyStringValueMap()

	return &c
}

// Configure object by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config: ConfigParams configuration parameters to be set.
func (c *AzureFunctionClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.DependencyResolver.Configure(ctx, config)

	c.Retries = config.GetAsIntegerWithDefault("options.retries", DefaultRetriesCount)
	c.ConnectTimeout = config.GetAsIntegerWithDefault("options.connectTimeout", DefaultConnectTimeout)
	c.Timeout = config.GetAsIntegerWithDefault("options.timeout", DefaultTimeout)
}

// SetReferences sets references to dependent components.
//
//	see IReferences
//	Parameters:
//		- ctx context.Context
//		- references IReferences references to locate the component dependencies.
func (c *AzureFunctionClient) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a services.InstrumentTiming object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- name string a method name.
//	Returns: services.InstrumentTiming object to end the time measurement.
func (c *AzureFunctionClient) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger.Trace(ctx, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")
	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")
	return rpctrace.NewInstrumentTiming(ctx, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// IsOpen Checks if the component is opened.
//
//	Returns: bool true if the component has been opened and false otherwise.
func (c *AzureFunctionClient) IsOpen() bool {
	return c.Client != nil
}

// Open opens the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Return: error
func (c *AzureFunctionClient) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}

	connection, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}

	c.Uri, _ = connection.FunctionUri()
	c.Client = &http.Client{
		// Timeout includes connection time, any redirects, and reading the response body
		Timeout: time.Duration(c.Timeout+c.ConnectTimeout) * time.Millisecond,
	}

	if c.Client == nil {
		return cerr.NewConnectionError(
			utils.ContextHelper.GetTraceId(ctx),
			"CANNOT_CONNECT",
			"Connection to Azure function service failed",
		).WithDetails("url", c.Uri)
	}

	c.Logger.Debug(ctx, "Azure function client connected to %s", c.Uri)

	return nil
}

// Closes component and frees used resources.
// Parameters:
//
//   - ctx context.Context	execution context to trace execution through call chain.
func (c *AzureFunctionClient) Close(ctx context.Context) error {
	if c.Client != nil {
		c.Logger.Debug(ctx, "Closed Azure function service at %s", c.Uri)
		c.Client = nil
		c.Uri = ""
	}
	return nil
}

// Performs Azure Function invocation.
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
//   - cmd	an action name to be called.
//   - args	action arguments
//
// Returns action result.
func (c *AzureFunctionClient) Call(ctx context.Context, cmd string,
	args *cdata.AnyValueMap) (*http.Response, error) {
	traceId := utils.ContextHelper.GetClient(ctx)
	if cmd == "" {
		cerr.NewUnknownError(traceId, "NO_COMMAND", "Cmd parameter is missing")
	}

	if !c.IsOpen() {
		return nil, nil
	}

	if traceId == "" {
		traceId = ckeys.IdGenerator.NextShort()
	}

	if args == nil {
		args = cdata.NewEmptyAnyValueMap()
	}
	args.Put("cmd", cmd)
	args.Put("trace_id", traceId)

	var jsonStr string
	if args != nil {
		jsonStr, _ = cconv.JsonConverter.ToJson(args.Value())
	}

	retries := c.Retries
	var response *http.Response

	for retries > 0 {
		req, err := c.prepareRequest(ctx, http.MethodPost, c.Uri, []byte(jsonStr))
		if err != nil {
			return nil, err
		}

		response, err = c.Client.Do(req)
		if err != nil {
			retries--
			if retries == 0 {
				return nil, cerr.NewUnknownError(
					traceId,
					"COMMUNICATION_ERROR",
					"Unknown communication problem on Azure client",
				).
					WithCause(err)
			}

			err = c.waitForRetry(ctx, traceId, retries)
			if err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	if response.StatusCode == 204 {
		_ = response.Body.Close()
		return nil, nil
	}

	if response.StatusCode >= 400 {
		defer response.Body.Close()
		return nil, c.handleResponseError(response, traceId)
	}

	return response, nil
}

// AddFilterParams method are adds filter parameters (with the same name as they defined)
// to invocation parameter map.
//
//	Parameters:
//		- params  *cdata.StringValueMap      invocation parameters.
//		- filter  *cdata.FilterParams     (optional) filter parameters
//	Returns: invocation parameters with added filter parameters.
func (c *AzureFunctionClient) AddFilterParams(params *cdata.StringValueMap, filter *cquery.FilterParams) *cdata.StringValueMap {

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	if filter != nil {
		for k, v := range filter.Value() {
			params.Put(k, v)
		}
	}
	return params
}

// AddPagingParams method are adds paging parameters (skip, take, total) to invocation parameter map.
// Parameters:
//   - params        invocation parameters.
//   - paging        (optional) paging parameters
//
// Return invocation parameters with added paging parameters.
func (c *AzureFunctionClient) AddPagingParams(params *cdata.StringValueMap, paging *cquery.PagingParams) *cdata.StringValueMap {
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if paging != nil {
		params.Put("total", paging.Total)
		if paging.Skip >= 0 {
			params.Put("skip", paging.Skip)
		}
		if paging.Take >= 0 {
			params.Put("take", paging.Take)
		}
	}

	return params
}

func (c *AzureFunctionClient) waitForRetry(ctx context.Context, traceId string, retries int) error {
	waitTime := c.Timeout * int(math.Pow(float64(c.Retries-retries), 2))

	select {
	case <-time.After(time.Duration(waitTime) * time.Millisecond):
		return nil
	case <-ctx.Done():
		return cerr.ApplicationErrorFactory.Create(
			&cerr.ErrorDescription{
				Type:     "Application",
				Category: "Application",
				Code:     "CONTEXT_CANCELLED",
				Message:  "request canceled by parent context",
				TraceId:  traceId,
			},
		)
	}
}

func (c *AzureFunctionClient) prepareRequest(ctx context.Context,
	method string, url string, body []byte) (*http.Request, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, cerr.NewUnknownError(
			utils.ContextHelper.GetTraceId(ctx),
			"UNSUPPORTED_METHOD",
			"Method is not supported by Azure client",
		).WithDetails("verb", method).WithCause(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers.Value() {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *AzureFunctionClient) handleResponseError(response *http.Response, traceId string) error {
	r, rErr := ioutil.ReadAll(response.Body)
	if rErr != nil {
		eDesct := cerr.ErrorDescription{
			Type:     "Application",
			Category: "Application",
			Status:   response.StatusCode,
			Code:     "",
			Message:  rErr.Error(),
			TraceId:  traceId,
		}
		return cerr.ApplicationErrorFactory.Create(&eDesct).WithCause(rErr)
	}

	appErr := cerr.ApplicationError{}
	_ = json.Unmarshal(r, &appErr)
	if appErr.Status == 0 && len(r) > 0 { // not standart Pip.Services error
		values := make(map[string]any)
		decodeErr := json.Unmarshal(r, &values)
		if decodeErr != nil { // not json response
			appErr.Message = (string)(r)
		}
		appErr.Details = values
	}
	appErr.Status = response.StatusCode
	return &appErr
}
