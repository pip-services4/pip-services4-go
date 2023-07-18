package build

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	pcount "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/count"
	pservices "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/services"
)

// DefaultPrometheusFactory creates Prometheus components by their descriptors.
// See: Factory
// See: PrometheusCounters
// See: PrometheusMetricsService
type DefaultPrometheusFactory struct {
	*cbuild.Factory
}

// NewDefaultPrometheusFactory are create a new instance of the factory.
func NewDefaultPrometheusFactory() *DefaultPrometheusFactory {
	c := DefaultPrometheusFactory{}
	c.Factory = build.NewFactory()

	prometheusCountersDescriptor := cref.NewDescriptor("pip-services", "counters", "prometheus", "*", "1.0")
	prometheusMetricsServiceDescriptor := cref.NewDescriptor("pip-services", "metrics-service", "prometheus", "*", "1.0")

	c.RegisterType(prometheusCountersDescriptor, pcount.NewPrometheusCounters)
	c.RegisterType(prometheusMetricsServiceDescriptor, pservices.NewPrometheusMetricsService)
	return &c
}
