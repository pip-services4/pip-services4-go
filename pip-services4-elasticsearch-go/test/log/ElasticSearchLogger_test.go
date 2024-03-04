package test_log

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	elog "github.com/pip-services4/pip-services4-go/pip-services4-elasticsearch-go/log"
	"github.com/stretchr/testify/assert"
)

func TestElasticSearchLogger(t *testing.T) {
	var logger *elog.ElasticSearchLogger
	var fixture *LoggerFixture

	var host = os.Getenv("ELASTICSEARCH_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	var port = os.Getenv("ELASTICSEARCH_SERVICE_PORT")
	if port == "" {
		port = "9200"
	}

	logger = elog.NewElasticSearchLogger()
	fixture = NewLoggerFixture(logger.CachedLogger)

	config := cconf.NewConfigParamsFromTuples(
		"source", "test",
		"index", "log",
		"daily", true,
		"connection.host", host,
		"connection.port", port,
		"connection.protocol", "http",
		// "options.include_type_name", true, // Elasticsearch 6.x
	)
	logger.Configure(context.Background(), config)

	opnErr := logger.Open(context.Background())

	assert.Nil(t, opnErr)

	defer logger.Close(context.Background())

	t.Run("Log Level", fixture.TestLogLevel)
	t.Run("Simple Logging", fixture.TestSimpleLogging)
	t.Run("Error Logging", fixture.TestErrorLogging)

}
