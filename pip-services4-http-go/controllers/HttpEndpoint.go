package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"
	"goji.io/pat"
	"goji.io/pattern"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"goji.io"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// HttpEndpoint used for creating HTTP endpoints. An endpoint is a URL,
// at which a given service can be accessed by a client.
//
//	Configuration parameters:
//		Parameters to pass to the configure method for component configuration:
//
//		- cors_headers - a comma-separated list of allowed CORS headers
//		- cors_origins - a comma-separated list of allowed CORS origins
//		- connection(s) - the connection resolver"s connections:
//			- "connection.discovery_key" - the key to use for connection resolving in a discovery service;
//			- "connection.protocol" - the connection"s protocol;
//			- "connection.host" - the target host;
//			- "connection.port" - the target port;
//			- "connection.uri" - the target URI.
//		- credential - the HTTPS credentials:
//			- "credential.ssl_key_file" - the SSL func (c *HttpEndpoint )key in PEM
//			- "credential.ssl_crt_file" - the SSL certificate in PEM
//			- "credential.ssl_ca_file" - the certificate authorities (root cerfiticates) in PEM
//
//	References:
//		A logger, counters, and a connection resolver can be referenced by passing the
//		following references to the object"s setReferences method:
//		- logger: "*:logger:*:*:1.0";
//		- counters: "*:counters:*:*:1.0";
//		- discovery: "*:discovery:*:*:1.0" (for the connection resolver).
//
//	Examples:
//		endpoint := NewHttpEndpoint();
//		endpoint.Configure(context.Background(), config);
//		endpoint.SetReferences(context.Background(), references);
//		...
//		endpoint.Open(context.Background())
type HttpEndpoint struct {
	defaultConfig          *cconf.ConfigParams
	server                 *http.Server
	mux                    *goji.Mux
	connectionResolver     *connect.HttpConnectionResolver
	logger                 *clog.CompositeLogger
	counters               *ccount.CompositeCounters
	maintenanceEnabled     bool
	fileMaxSize            int64
	protocolUpgradeEnabled bool
	uri                    string
	registrations          []IRegisterable
	allowedHeaders         []string
	allowedOrigins         []string
}

const (
	DefaultConnectionTimeout = "60000"
	DefaultFileMaxSize       = 200 * 1024 * 1024
	DefaultRequestMaxSize    = 1024 * 1024
)

// NewHttpEndpoint creates new HttpEndpoint
func NewHttpEndpoint() *HttpEndpoint {
	c := HttpEndpoint{}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", "3000",

		"credential.ssl_key_file", nil,
		"credential.ssl_crt_file", nil,
		"credential.ssl_ca_file", nil,

		"options.maintenance_enabled", false,
		"options.request_max_size", DefaultRequestMaxSize,
		"options.file_max_size", DefaultFileMaxSize,
		"options.connect_timeout", DefaultConnectionTimeout,
		"options.debug", "true",
	)
	c.connectionResolver = connect.NewHttpConnectionResolver()
	c.logger = clog.NewCompositeLogger()
	c.counters = ccount.NewCompositeCounters()
	c.maintenanceEnabled = false
	c.fileMaxSize = DefaultFileMaxSize
	c.protocolUpgradeEnabled = false
	c.registrations = make([]IRegisterable, 0)
	c.allowedHeaders = []string{
		//"Accept",
		//"Content-Type",
		//"Content-Length",
		//"Accept-Encoding",
		//"X-CSRF-Token",
		//"Authorization",
		"trace_id",
		//"access_token",
	}
	c.allowedOrigins = make([]string, 0)
	return &c
}

// Configure method are configures this HttpEndpoint using the given configuration parameters.
//
//	Configuration parameters:
//		- connection(s) - the connection resolver"s connections;
//			- "connection.discovery_key" - the key to use for connection resolving in a discovery service;
//			- "connection.protocol" - the connection"s protocol;
//			- "connection.host" - the target host;
//			- "connection.port" - the target port;
//			- "connection.uri" - the target URI.
//			- "credential.ssl_key_file" - SSL func (c *HttpEndpoint )key in PEM
//			- "credential.ssl_crt_file" - SSL certificate in PEM
//			- "credential.ssl_ca_file" - Certificate authority (root certificate) in PEM
//	Parameters:
//		- ctx context.Context
//		- config    configuration parameters, containing a "connection(s)" section.
func (c *HttpEndpoint) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.connectionResolver.Configure(ctx, config)

	c.maintenanceEnabled = config.GetAsBooleanWithDefault("options.maintenance_enabled", c.maintenanceEnabled)
	c.fileMaxSize = config.GetAsLongWithDefault("options.file_max_size", c.fileMaxSize)
	c.protocolUpgradeEnabled = config.GetAsBooleanWithDefault("options.protocol_upgrade_enabled", c.protocolUpgradeEnabled)

	headers := strings.Split(config.GetAsStringWithDefault("cors_headers", ""), ",")
	if len(headers) > 0 {
		for _, header := range headers {
			c.AddCorsHeader(strings.TrimSpace(header), "")
		}
	}

	origins := strings.Split(config.GetAsStringWithDefault("cors_origins", ""), ",")
	if len(origins) > 0 {
		for _, origin := range origins {
			c.AddCorsHeader("", strings.TrimSpace(origin))
		}
	}
}

// SetReferences method are sets references to this endpoint"s logger, counters, and connection resolver.
//
//	References:
//		- logger: "*:logger:*:*:1.0"
//		- counters: "*:counters:*:*:1.0"
//		- discovery: "*:discovery:*:*:1.0" (for the connection resolver)
//	Parameters:
//		- references an IReferences object, containing references to a logger, counters,
//			and a connection resolver.
func (c *HttpEndpoint) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.logger.SetReferences(ctx, references)
	c.counters.SetReferences(ctx, references)
	c.connectionResolver.SetReferences(ctx, references)
}

// IsOpen method is  whether this endpoint is open with an actively listening REST server.
func (c *HttpEndpoint) IsOpen() bool {
	return c.server != nil
}

// Open a connection using the parameters resolved by the referenced connection
// resolver and creates a REST server (service) using the set options and parameters.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error an error if one is raised.
func (c *HttpEndpoint) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}
	connection, credential, err := c.connectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}

	c.uri = connection.Uri()
	url := connection.Host() + ":" + strconv.Itoa(connection.Port())

	c.mux = goji.NewMux()
	c.server = &http.Server{Addr: url, Handler: c.mux}
	// Provide container context to http handler
	if ctx != nil {
		c.server.BaseContext = func(listener net.Listener) context.Context {
			return ctx
		}
	}

	c.mux.Use(cors.New(cors.Options{
		AllowedOrigins: c.allowedOrigins,
		AllowedMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
			"PATCH",
		},
		AllowedHeaders: c.allowedHeaders,
	}).Handler)

	c.mux.Use(c.noCache)
	c.mux.Use(c.doMaintenance)

	c.performRegistrations()

	if connection.Protocol() == "https" {
		sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")

		go func() {
			defer cctx.DefaultErrorHandlerWithShutdown(ctx)

			servErr := c.server.ListenAndServeTLS(sslKeyFile, sslCrtFile)
			if servErr != nil && !errors.Is(servErr, http.ErrServerClosed) {
				cctx.SendShutdownSignalWithErr(ctx, servErr)
			}
		}()

	} else {

		go func() {
			defer cctx.DefaultErrorHandlerWithShutdown(ctx)

			servErr := c.server.ListenAndServe()
			if servErr != nil && !errors.Is(servErr, http.ErrServerClosed) {
				cctx.SendShutdownSignalWithErr(ctx, servErr)
			}
		}()
	}

	regErr := c.connectionResolver.Register(ctx)
	if regErr != nil {
		c.logger.Error(ctx, regErr, "ERROR_REG_SRV", "Can't register REST service at %s", c.uri)
	}
	c.logger.Debug(ctx, "Opened REST service at %s", c.uri)
	return regErr
}

// noCache prevents IE from caching REST requests
func (c *HttpEndpoint) noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// doMaintenance returns maintenance error code
func (c *HttpEndpoint) doMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make this more sophisticated
		if c.maintenanceEnabled {
			w.Header().Add("Retry-After", "3600")
			jsonStr, _ := cconv.JsonConverter.ToJson(503)
			_, _ = w.Write([]byte(jsonStr))
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Close method are closes this endpoint and the REST server (service) that was opened earlier.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error an error if one is raised.
func (c *HttpEndpoint) Close(ctx context.Context) error {
	if c.server != nil {
		// Attempt a graceful shutdown
		_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		clErr := c.server.Shutdown(_ctx)
		if clErr != nil {
			c.logger.Warn(ctx, "Failed while closing REST service: %s", clErr.Error())
			return clErr
		}
		c.logger.Debug(ctx, "Closed REST service at %s", c.uri)
		c.server = nil
		c.uri = ""
	}
	return nil
}

// Register a registrable object for dynamic endpoint discovery.
//
//	Parameters:
//		- registration IRegisterable implements of IRegisterable interface.
//	See IRegisterable
func (c *HttpEndpoint) Register(registration IRegisterable) {
	c.registrations = append(c.registrations, registration)
}

// Unregister a registerable object, so that it is no longer used in dynamic
// endpoint discovery.
//
//	Parameters:
//		- registration  IRegisterable  the registration to remove.
//	See IRegisterable
func (c *HttpEndpoint) Unregister(registration IRegisterable) {
	for i := range c.registrations {
		if c.registrations[i] == registration {
			if i == len(c.registrations)-1 {
				c.registrations = c.registrations[:i]
			} else {
				c.registrations = append(c.registrations[:i], c.registrations[i+1:]...)
			}
		} else {
			i++
		}
	}
}

func (c *HttpEndpoint) performRegistrations() {
	for _, registration := range c.registrations {
		registration.Register()
	}
}

func (c *HttpEndpoint) fixRoute(route string) string {
	if len(route) > 0 && !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	return route
}

// GetTraceId method returns traceId from request
//
//	Parameters:
//		- req *http.Request  request
//	Returns: string trace_id or empty string
func (c *HttpEndpoint) GetTraceId(req *http.Request) string {
	traceId := req.URL.Query().Get("trace_id")
	if traceId == "" {
		traceId = req.Header.Get("trace_id")
	}
	return traceId
}

// RegisterRoute method are registers an action in this objects REST server (service)
// by the given method and route.
//
//	Parameters:
//		- method   string     the HTTP method of the route.
//		- route    string     the route to register in this object"s REST server (service).
//		- schema   *cvalid.Schema     the schema to use for parameter validation.
//		- action   http.HandlerFunc     the action to perform at the given route.
func (c *HttpEndpoint) RegisterRoute(method string, route string, schema *cvalid.Schema,
	action http.HandlerFunc) {

	method = strings.ToLower(method)
	if method == "del" {
		method = "delete"
	}
	route = c.fixRoute(route)
	actionCurl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					msg := cconv.StringConverter.ToString(r)
					err = errors.New(msg)
				}
				ctx := utils.ContextHelper.NewContextWithTraceId(r.Context(), c.GetTraceId(r))
				c.logger.Error(ctx, err, "http handler panics with error")
			}
		}()
		//  Perform validation
		if schema != nil {
			var params = make(map[string]any, 0)
			for k, v := range r.URL.Query() {
				params[k] = v[0]
			}

			if reqVars, ok := r.Context().Value(pattern.AllVariables).(map[pattern.Variable]any); ok {
				for k, v := range reqVars {
					params[string(k)] = v
				}
			}

			// Make copy of request
			bodyBuf, bodyErr := io.ReadAll(r.Body)
			if bodyErr != nil {
				HttpResponseSender.SendError(w, r, bodyErr)
				return
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			var body any
			_ = json.Unmarshal(bodyBuf, &body)
			params["body"] = body

			traceId := c.GetTraceId(r)
			err := schema.ValidateAndReturnError(traceId, params, false)
			if err != nil {
				HttpResponseSender.SendError(w, r, err)
				return
			}
		}
		action(w, r)
	})
	c.mux.HandleFunc(pat.NewWithMethods(route, strings.ToUpper(method)), actionCurl)
}

// RegisterRouteWithAuth method are registers an action with authorization in this objects REST server (service)
// by the given method and route.
// Parameters:
//   - method    string    the HTTP method of the route.
//   - route     string    the route to register in this object"s REST server (service).
//   - schema    *cvalid.Schema    the schema to use for parameter validation.
//   - authorize     the authorization interceptor
//   - action        the action to perform at the given route.
func (c *HttpEndpoint) RegisterRouteWithAuth(method string, route string, schema *cvalid.Schema,
	authorize func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc),
	action http.HandlerFunc) {

	if authorize != nil {
		nextAction := action
		action = func(w http.ResponseWriter, r *http.Request) {
			authorize(w, r, nextAction)
		}
	}

	c.RegisterRoute(method, route, schema, action)
}

// RegisterInterceptor method are registers a middleware action for the given route.
// Parameters:
//   - route         the route to register in this object"s REST server (service).
//   - action        the middleware action to perform at the given route.
func (c *HttpEndpoint) RegisterInterceptor(route string, action func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {

	route = c.fixRoute(route)
	interceptorFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched, _ := regexp.MatchString(route, r.URL.Path)
			if route != "" && !matched {
				next.ServeHTTP(w, r)
			} else {
				action(w, r, next.ServeHTTP)
			}
		})
	}
	c.mux.Use(interceptorFunc)
}

// AddCorsHeader method adds allowed header, ignore if it already exists
// must be called before to opening endpoint
func (c *HttpEndpoint) AddCorsHeader(header string, origin string) {

	if len(header) > 0 {
		contain := false
		for _, allowedHeader := range c.allowedHeaders {
			if allowedHeader == header {
				contain = true
				break
			}
		}
		if !contain {
			c.allowedHeaders = append(c.allowedHeaders, header)
		}
	}
	if len(origin) > 0 {
		contain := false
		for _, allowedOrigin := range c.allowedOrigins {
			if allowedOrigin == origin {
				contain = true
				break
			}
		}
		if !contain {
			c.allowedOrigins = append(c.allowedOrigins, origin)
		}
	}
}
