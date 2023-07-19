package test_services

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	pservice "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/controllers"
	pcount "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/count"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMetricsController(t *testing.T) {
	ctx := context.Background()
	var service *pservice.PrometheusMetricsController
	var counters *pcount.PrometheusCounters

	var restConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	service = pservice.NewPrometheusMetricsController()
	service.Configure(ctx, restConfig)

	counters = pcount.NewPrometheusCounters()

	contextInfo := cctx.NewContextInfo()
	contextInfo.Name = "Test"
	contextInfo.Description = "This is a test container"

	references := cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
		cref.NewDescriptor("pip-services", "counters", "prometheus", "default", "1.0"), counters,
		cref.NewDescriptor("pip-services", "metrics-service", "prometheus", "default", "1.0"), service,
	)
	counters.SetReferences(ctx, references)
	service.SetReferences(ctx, references)

	opnErr := counters.Open(ctx)
	if opnErr == nil {
		service.Open(ctx)
	}

	defer service.Close(ctx)
	defer counters.Close(ctx)

	var url = "http://localhost:3000"

	counters.IncrementOne(ctx, "test.counter1")
	counters.Stats(ctx, "test.counter2", 2)
	counters.Last(ctx, "test.counter3", 3)
	counters.TimestampNow(ctx, "test.counter4")

	getRes, getErr := http.Get(url + "/metrics")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
	assert.True(t, getRes.StatusCode < 400)
	body, _ := ioutil.ReadAll(getRes.Body)
	assert.True(t, len(body) > 0)
}
