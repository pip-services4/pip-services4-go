package auth

import (
	"net/http"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type BasicAuthManager struct {
}

func (c *BasicAuthManager) Anybody() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		next.ServeHTTP(res, req)
	}
}

func (c *BasicAuthManager) Signed() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
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
