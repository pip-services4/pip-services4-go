package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"goji.io/pat"
)

// RestOperations helper class for REST operations
type RestOperations struct {
	Logger             *clog.CompositeLogger
	Counters           *ccount.CompositeCounters
	DependencyResolver *crefer.DependencyResolver
}

// NewRestOperations creates new instance of RestOperations
func NewRestOperations() *RestOperations {
	ro := RestOperations{}
	ro.Logger = clog.NewCompositeLogger()
	ro.Counters = ccount.NewCompositeCounters()
	ro.DependencyResolver = crefer.NewDependencyResolver()
	return &ro
}

// Configure method are configures this RestOperations using the given configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *cconf.ConfigParams confif parameters
func (c *RestOperations) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(ctx, config)
}

// SetReferences method are sets references to this RestOperations logger, counters, and connection resolver.
//
//	Parameters:
//		- ctx context.Context
//		- references    an IReferences object, containing references to a logger, counters,
//			and a dependency resolver.
func (c *RestOperations) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// GetTraceId method returns traceId from request
//
//	Parameters:
//		- req *http.Request  request
//	Returns: string trace_id or empty string
func (c *RestOperations) GetTraceId(req *http.Request) string {
	traceId := req.URL.Query().Get("trace_id")
	if traceId == "" {
		traceId = req.Header.Get("trace_id")
	}
	return traceId
}

// GetFilterParams method reruns filter params object from request
//
//	Parameters:
//		- req *http.Request  request
//	Returns: *cdata.FilterParams filter params object
func (c *RestOperations) GetFilterParams(req *http.Request) *cquery.FilterParams {

	params := req.URL.Query()
	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")
	filter := cquery.NewFilterParamsFromValue(
		params,
	)
	return filter
}

// GetPagingParams method reruns paging params object from request
//
//	Parameters:
//		- req *http.Request  request
//	Returns: *cdata.PagingParams pagings params object
func (c *RestOperations) GetPagingParams(req *http.Request) *cquery.PagingParams {

	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	paging := cquery.NewPagingParamsFromValue(
		paginParams,
	)
	return paging
}

// GetParam methods helps get all params from query
//
//		Parameters:
//	  - req  incoming request
//	  - name parameter name
//
// Returns: value or empty string if param not exists
func (c *RestOperations) GetParam(req *http.Request, name string) string {
	param := req.URL.Query().Get(name)
	if param == "" {
		param = pat.Param(req, name)
	}
	return param
}

// DecodeBody methods helps decode body
//
//	Parameters:
//		- req incoming request
//		- target pointer on target variable for decode
//
// Returns: error
func (c *RestOperations) DecodeBody(req *http.Request, target any) error {
	bodyBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, target)

	if err != nil {
		return err
	}

	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

func (c *RestOperations) SendResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	HttpResponseSender.SendResult(res, req, result, err)
}

func (c *RestOperations) SendEmptyResult(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendEmptyResult(res, req, err)
}

func (c *RestOperations) SendCreatedResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *RestOperations) SendDeletedResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	HttpResponseSender.SendDeletedResult(res, req, result, err)
}

func (c *RestOperations) SendError(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendError(res, req, err)
}

func (c *RestOperations) SendBadRequest(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewBadRequestError(traceId, "BAD_REQUEST", message)
	c.SendError(res, req, err)
}

func (c *RestOperations) SendUnauthorized(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewUnauthorizedError(traceId, "UNAUTHORIZED", message)
	c.SendError(res, req, err)
}

func (c *RestOperations) SendNotFound(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewNotFoundError(traceId, "NOT_FOUND", message)
	c.SendError(res, req, err)
}

func (c *RestOperations) SendConflict(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewConflictError(traceId, "CONFLICT", message)
	c.SendError(res, req, err)
}

func (c *RestOperations) SendSessionExpired(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewUnknownError(traceId, "SESSION_EXPIRED", message)
	err.Status = 440
	c.SendError(res, req, err)
}

func (c *RestOperations) SendInternalError(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewUnknownError(traceId, "INTERNAL", message)
	c.SendError(res, req, err)
}

func (c *RestOperations) SendServerUnavailable(res http.ResponseWriter, req *http.Request, message string) {
	traceId := c.GetTraceId(req)
	err := cerr.NewConflictError(traceId, "SERVER_UNAVAILABLE", message)
	err.Status = 503
	c.SendError(res, req, err)
}
