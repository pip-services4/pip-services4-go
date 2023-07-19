package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	pcontrl "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/controllers"
	pcount "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/count"
)

// DefaultPrometheusFactory creates Prometheus components by their descriptors.
// See: Factory
// See: PrometheusCounters
// See: PrometheusMetricsController
type DefaultPrometheusFactory struct {
	*cbuild.Factory
}

// NewDefaultPrometheusFactory are create a new instance of the factory.
func NewDefaultPrometheusFactory() *DefaultPrometheusFactory {
	c := DefaultPrometheusFactory{}
	c.Factory = cbuild.NewFactory()

	prometheusCountersDescriptor := cref.NewDescriptor("pip-services", "counters", "prometheus", "*", "1.0")
	PrometheusMetricsControllerDescriptor := cref.NewDescriptor("pip-services", "metrics-service", "prometheus", "*", "1.0")

	c.RegisterType(prometheusCountersDescriptor, pcount.NewPrometheusCounters)
	c.RegisterType(PrometheusMetricsControllerDescriptor, pcontrl.NewPrometheusMetricsController)
	return &c
}
