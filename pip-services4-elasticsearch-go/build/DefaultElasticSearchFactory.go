package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	elog "github.com/pip-services4/pip-services4-go/pip-services4-elasticsearch-go/log"
)

/*
DefaultElasticSearchFactory are creates ElasticSearch components by their descriptors.
See ElasticSearchLogger
*/
type DefaultElasticSearchFactory struct {
	*cbuild.Factory
}

// NewDefaultElasticSearchFactory create a new instance of the factory.
// Retruns *DefaultElasticSearchFactory
// pointer on new factory
func NewDefaultElasticSearchFactory() *DefaultElasticSearchFactory {
	c := DefaultElasticSearchFactory{}
	c.Factory = cbuild.NewFactory()

	elasticSearchLoggerDescriptor := cref.NewDescriptor("pip-services", "logger", "elasticsearch", "*", "1.0")

	c.RegisterType(elasticSearchLoggerDescriptor, elog.NewElasticSearchLogger)

	return &c
}
