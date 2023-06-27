package test_clients

import (
	"context"
	"testing"

	test_sample "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/stretchr/testify/assert"
)

func TestRetriesRestClient(t *testing.T) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "12345",

		"options.retries", "4",
		"options.timeout", "100",
		"options.connect_timeout", "100",
	)

	var client *DummyRestClient

	client = NewDummyRestClient()

	client.Configure(context.Background(), restConfig)
	client.SetReferences(context.Background(), cref.NewEmptyReferences())
	client.Open(context.Background())

	res, err := client.GetDummyById(context.Background(), "1")
	assert.NotNil(t, err)
	assert.Equal(t, test_sample.Dummy{}, res)

}
