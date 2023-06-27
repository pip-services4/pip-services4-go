package clients

import (
	"encoding/json"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/protos"
	grpcproto "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/protos"
)

// HandleHttpResponse method helps handle http response body
//
//	Parameters:
//		- r *grpcproto.InvokeReply response object from grpc server
//	Returns: T any result, err error
func HandleHttpResponse[T any](r *grpcproto.InvokeReply) (T, error) {
	var defaultValue T

	// Handle error response
	if r.Error != nil {
		var errDesc cerr.ErrorDescription
		errDescJson, _ := json.Marshal(r.Error)
		json.Unmarshal(errDescJson, &errDesc)
		err := cerr.ApplicationErrorFactory.Create(&errDesc)
		return defaultValue, err
	}

	// Handle empty response
	if r.ResultEmpty || r.ResultJson == "" {
		return defaultValue, nil
	}

	return convert.NewDefaultCustomTypeJsonConvertor[T]().FromJson(r.ResultJson)
}

func ToError(obj *protos.ErrorDescription) error {
	if obj == nil || (obj.Category == "" && obj.Message == "") {
		return nil
	}

	description := &errors.ErrorDescription{
		// Type:          obj.Type,
		Category:   obj.Category,
		Code:       obj.Code,
		TraceId:    obj.TraceId,
		Status:     convert.IntegerConverter.ToInteger(obj.Status),
		Message:    obj.Message,
		Cause:      obj.Cause,
		StackTrace: obj.StackTrace,
		Details:    ToMap(obj.Details),
	}

	return errors.ApplicationErrorFactory.Create(description)
}

func ToMap(val map[string]string) map[string]interface{} {
	r := make(map[string]interface{}, 0)

	for k, v := range val {
		r[k] = v
	}

	return r
}
