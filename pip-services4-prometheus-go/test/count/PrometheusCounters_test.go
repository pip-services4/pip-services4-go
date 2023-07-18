package test_count

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	pcount "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/count"
	pfixture "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/test/fixture"
)

func TestPrometheusCounters(t *testing.T) {
	var counters *pcount.PrometheusCounters
	var fixture *pfixture.CountersFixture

	host := os.Getenv("PUSHGATEWAY_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PUSHGATEWAY_SERVICE_PORT")
	if port == "" {
		port = "9091"
	}
	counters = pcount.NewPrometheusCounters()
	fixture = pfixture.NewCountersFixture(counters.CachedCounters)

	config := cconf.NewConfigParamsFromTuples(
		"source", "test",
		"connection.host", host,
		"connection.port", port,
		"connection.protocol", "http",
	)
	counters.Configure(context.Background(), config)

	counters.Open(context.Background())

	defer counters.Close(context.Background())

	t.Run("Simple Counters", fixture.TestSimpleCounters)
	t.Run("Measure Elapsed Time", fixture.TestMeasureElapsedTime)
}
