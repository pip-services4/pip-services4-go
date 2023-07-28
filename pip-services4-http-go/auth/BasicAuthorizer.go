package auth

import (
	"net/http"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type BasicAuthorizer struct {
}

func (c *BasicAuthorizer) Anybody() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		next.ServeHTTP(res, req)
	}
}

func (c *BasicAuthorizer) Signed() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		_, ok := req.Context().Value(PipAuthUser).(cdata.AnyValueMap)
		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("",
					"NOT_SIGNED",
					"User must be signed in to perform this operation",
				).WithStatus(401),
			)
		} else {
			next.ServeHTTP(res, req)
		}
	}
}
