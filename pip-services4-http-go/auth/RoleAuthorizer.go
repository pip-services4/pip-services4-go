package auth

import (
	"net/http"
	"strings"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type RoleAuthorizer struct {
}

func (c *RoleAuthorizer) UserInRoles(roles []string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		user, ok := req.Context().Value(PipAuthUser).(cdata.AnyValueMap)
		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("", "NOT_SIGNED",
					"User must be signed in to perform this operation").WithStatus(401))
		} else {
			authorized := false
			userRoles := user.GetAsArray(string(PipAuthRoles))

			if userRoles == nil {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError("", "NOT_SIGNED",
						"User must be signed in to perform this operation").WithStatus(401))
				return
			}

			for _, role := range roles {
				for _, userRole := range userRoles.Value() {
					r, ok := userRole.(string)
					if ok && role == r {
						authorized = true
					}
				}
			}

			if !authorized {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError(
						"", "NOT_IN_ROLE",
						"User must be "+strings.Join(roles, " or ")+" to perform this operation").WithDetails(string(PipAuthRoles), roles).WithStatus(403))
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}

func (c *RoleAuthorizer) UserInRole(role string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return c.UserInRoles([]string{role})
}

func (c *RoleAuthorizer) Admin() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return c.UserInRole(string(PipAuthAdmin))
}
