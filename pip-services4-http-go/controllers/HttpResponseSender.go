package services

import (
	"io"
	"net/http"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
)

// HttpResponseSender helper class that handles HTTP-based responses.
var HttpResponseSender = _THttpResponseSender{}

type _THttpResponseSender struct {
}

// SendError sends error serialized as ErrorDescription object
// and appropriate HTTP status code.
// If status code is not defined, it uses 500 status code.
//
//	Parameters:
//		- req  *http.Request     a HTTP request object.
//		- res  http.ResponseWriter     a HTTP response object.
//		- err  error     an error object to be sent.
func (c *_THttpResponseSender) SendError(res http.ResponseWriter, req *http.Request, err error) {

	appErr := cerr.ApplicationError{
		Status: 500,
	}
	appErr = *appErr.Wrap(err)
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(appErr.Status)
	jsonObjStr, jsonErr := cconv.JsonConverter.ToJson(appErr)
	if jsonErr == nil {
		_, _ = io.WriteString(res, jsonObjStr)
	}
}

// SendResult sends result as JSON object.
// That function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 200 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
//
//	Parameters:
//		- req  *http.Request     a HTTP request object.
//		- res  http.ResponseWriter     a HTTP response object.
//		- result any  result object to be send
//		- err  error     an error object to be sent.
func (c *_THttpResponseSender) SendResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	if err != nil {
		HttpResponseSender.SendError(res, req, err)
		return
	}
	if result == nil {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(204)
	} else {
		res.Header().Add("Content-Type", "application/json")
		jsonObjStr, jsonErr := cconv.JsonConverter.ToJson(result)
		if jsonErr == nil {
			_, _ = io.WriteString(res, jsonObjStr)
		}
	}
}

// SendEmptyResult are sends an empty result with 204 status code.
// If error occur it sends ErrorDescription with appropriate status code.
//
//	Parameters:
//		- req  *http.Request     a HTTP request object.
//		- res  http.ResponseWriter     a HTTP response object.
func (c *_THttpResponseSender) SendEmptyResult(res http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		HttpResponseSender.SendError(res, req, err)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(204)
}

// SendCreatedResult are sends newly created object as JSON.
// That function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 201 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
//
//	Parameters:
//		- req  *http.Request     a HTTP request object.
//		- res  http.ResponseWriter     a HTTP response object.
func (c *_THttpResponseSender) SendCreatedResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	if err != nil {
		HttpResponseSender.SendError(res, req, err)
		return
	}
	if result == nil {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(204)
	} else {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(201)
		jsonObjStr, jsonErr := cconv.JsonConverter.ToJson(result)
		if jsonErr == nil {
			_, _ = io.WriteString(res, jsonObjStr)
		}
	}
}

// SendDeletedResult are sends deleted object as JSON.
// That function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 200 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
//
//	Parameters:
//		- req  *http.Request     a HTTP request object.
//		- res  http.ResponseWriter     a HTTP response object.
func (c *_THttpResponseSender) SendDeletedResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	if err != nil {
		HttpResponseSender.SendError(res, req, err)
		return
	}
	if result == nil {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(204)
	} else {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(200)
		jsonObjStr, jsonErr := cconv.JsonConverter.ToJson(result)
		if jsonErr == nil {
			_, _ = io.WriteString(res, jsonObjStr)
		}
	}
}
