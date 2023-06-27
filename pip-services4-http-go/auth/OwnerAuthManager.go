package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type OwnerAuthManager struct {
}

func (c *OwnerAuthManager) Owner(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = string(PipAuthUserId)
	}
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
			userId := req.URL.Query().Get(idParam)
			if userId == "" {
				userId = mux.Vars(req)[idParam]
			}

			reqUserId, ok := req.Context().Value(PipAuthUserId).(string)
			if !ok || reqUserId != userId {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError(
						"", "FORBIDDEN",
						"Only data owner can perform this operation",
					).WithStatus(403),
				)
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}

func (c *OwnerAuthManager) OwnerOrAdmin(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = string(PipAuthUserId)
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		user, ok := req.Context().Value(PipAuthUser).(cdata.AnyValueMap)

		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("",
					"NOT_SIGNED",
					"User must be signed in to perform this operation",
				).WithStatus(401),
			)
		} else {

			userId := req.URL.Query().Get(idParam)
			if userId == "" {
				userId = mux.Vars(req)[idParam]
			}
			roles := user.GetAsArray(string(PipAuthRoles))
			admin := false
			for _, role := range roles.Value() {
				r, ok := role.(string)
				if ok && r == string(PipAuthAdmin) {
					admin = true
					break
				}
			}

			reqUserId, ok := req.Context().Value(PipAuthUserId).(string)
			if !ok || reqUserId != userId && !admin {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError("",
						"FORBIDDEN",
						"Only data owner can perform this operation",
					).WithStatus(403),
				)
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}
