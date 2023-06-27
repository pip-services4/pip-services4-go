package test_controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
	"github.com/stretchr/testify/assert"
)

func TestHttpEndpoint(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", HttpEndpointControllertPort)

	getResponse, getErr := http.Get(url + "/api/v1/dummies")
	assert.Nil(t, getErr)
	resBody, bodyErr := ioutil.ReadAll(getResponse.Body)
	assert.Nil(t, bodyErr)
	var dummies *cdata.DataPage[tdata.Dummy]
	jsonErr := json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.False(t, dummies.HasData())
	assert.Len(t, dummies.Data, 0)
}
