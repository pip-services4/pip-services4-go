package clients_test

import (
	"context"
	"os"
	"testing"
	"time"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	clients1 "github.com/pip-services4/pip-services4-go/pip-services4-datadog-go/clients"
	"github.com/stretchr/testify/assert"
)

func TestDataDogLogClient(t *testing.T) {
	var client *clients1.DataDogLogClient
	ctx := context.Background()

	apiKey := os.Getenv("DATADOG_API_KEY")
	if apiKey == "" {
		apiKey = "3eb3355caf628d4689a72084425177ac"
	}

	client = clients1.NewDataDogLogClient(nil)

	config := cconf.NewConfigParamsFromTuples(
		"source", "test",
		"credential.access_key", apiKey,
	)
	client.Configure(ctx, config)

	err := client.Open(ctx)
	assert.Nil(t, err)

	defer client.Close(ctx)

	t.Run("Send Logs", func(t *testing.T) {
		messages := []clients1.DataDogLogMessage{
			{
				Time:    time.Now().UTC(),
				Service: "TestService",
				Host:    "TestHost",
				Status:  clients1.Debug,
				Message: "Test trace message",
			},
			{
				Time:    time.Now().UTC(),
				Service: "TestService",
				Host:    "TestHost",
				Status:  clients1.Info,
				Message: "Test info message",
			},
			{
				Time:       time.Now().UTC(),
				Service:    "TestService",
				Host:       "TestHost",
				Status:     clients1.Error,
				Message:    "Test error message",
				ErrorKind:  "Exception",
				ErrorStack: "Stack trace...",
			},
			{
				Time:       time.Now().UTC(),
				Service:    "TestService",
				Host:       "TestHost",
				Status:     clients1.Emergency,
				Message:    "Test fatal message",
				ErrorKind:  "Exception",
				ErrorStack: "Stack trace...",
			},
		}

		err := client.SendLogs(ctx, messages)
		assert.Nil(t, err)

	})

}
