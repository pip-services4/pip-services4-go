package clients

import (
	"io/ioutil"
	"net/http"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
)

// HandleHttpResponse method helps handle http response body
//
//	Parameters:
//		- ctx context.Context
//		- traceId string (optional) transaction id to trace execution through call chain.
//	Returns: T any result, err error
func HandleHttpResponse[T any](r *http.Response, traceId string) (T, error) {
	var defaultValue T

	if r != nil {
		defer r.Body.Close()

		buffer, err := ioutil.ReadAll(r.Body)
		if err != nil {
			var defaultValue T
			return defaultValue, cerr.ApplicationErrorFactory.
				Create(&cerr.ErrorDescription{
					Type:     "Application",
					Category: "Application",
					Status:   r.StatusCode,
					Code:     "",
					Message:  err.Error(),
					TraceId:  traceId,
				}).
				WithCause(err)
		}

		return convert.NewDefaultCustomTypeJsonConvertor[T]().FromJson(string(buffer))

	}

	return defaultValue, nil
}
