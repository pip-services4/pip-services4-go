package clients

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/lambda"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// ConvertComandResult method helps get correct result from JSON by prototype
// Parameters:
//   - comRes any  input JSON string
//   - prototype reflect.Type output object prototype
//
// Returns: convRes any, err error
func HandleLambdaResponse[T any](data *lambda.InvokeOutput) (convRes T, err error) {
	if data.Payload != nil && len(data.Payload) > 0 {

		unesccapedResult, err := strconv.Unquote((string)(data.Payload))
		if err != nil {
			unesccapedResult = (string)(data.Payload)
		}

		return cconv.NewDefaultCustomTypeJsonConvertor[T]().FromJson(unesccapedResult)
	}

	return convRes, nil
}
