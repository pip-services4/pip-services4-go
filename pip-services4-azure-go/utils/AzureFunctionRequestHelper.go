package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
)

// Helper struct that allow prepare of requests data
var AzureFunctionRequestHelper = _TAzureFunctionRequestHelper{}

type _TAzureFunctionRequestHelper struct {
}

// Returns traceId from request struct
// Parameters:
//   - req	request struct
//
// Returns trace id string or empty
func (c *_TAzureFunctionRequestHelper) GetTraceId(req *http.Request) string {
	traceId := req.URL.Query().Get("trace_id")
	if traceId == "" {
		traceId = req.Header.Get("trace_id")
	}
	return traceId
}

// Returns command from request struct
// Parameters:
//   - req	request struct
//
// Returns command string or empty
func (c *_TAzureFunctionRequestHelper) GetCommand(req *http.Request) (string, error) {
	cmd := req.URL.Query().Get("cmd")

	if cmd == "" {
		var body map[string]any

		err := c.DecodeBody(req, &body)

		if err != nil {
			return "", err
		}

		if val, ok := body["cmd"].(string); ok {
			cmd = val
		}
	}

	return cmd, nil
}

// Returns body of request
// Parameters:
//   - req	request struct
//   - target	the target instance to which the result will be written
//
// Returns error
func (c *_TAzureFunctionRequestHelper) DecodeBody(req *http.Request, target any) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, target)

	if err != nil {
		return err
	}

	_ = req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

// Get body of request as Parameters struct
// Parameters:
//   - req	request struct
//
// Returns Parameters
func (c *_TAzureFunctionRequestHelper) GetParameters(req *http.Request) *cexec.Parameters {
	var params map[string]any

	_ = c.DecodeBody(req, &params) // Ignore the error

	return cexec.NewParametersFromValue(params)
}
