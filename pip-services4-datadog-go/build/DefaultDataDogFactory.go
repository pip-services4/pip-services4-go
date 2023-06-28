package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	count "github.com/pip-services4/pip-services4-go/pip-services4-datadog-go/count"
	log "github.com/pip-services4/pip-services4-go/pip-services4-datadog-go/log"
)

// DefaultDataDogFactory are creates DataDog components by their descriptors.
// See DataDogLogger
type DefaultDataDogFactory struct {
	*cbuild.Factory
}

// NewDefaultDataDogFactory create a new instance of the factory.
// Retruns *DefaultDataDogFactory
// pointer on new factory
func NewDefaultDataDogFactory() *DefaultDataDogFactory {
	c := DefaultDataDogFactory{}
	c.Factory = cbuild.NewFactory()
	dataDogLoggerDescriptor := cref.NewDescriptor("pip-services", "logger", "datadog", "*", "1.0")
	dataDogCountersDescriptor := cref.NewDescriptor("pip-services", "counters", "datadog", "*", "1.0")

	c.RegisterType(dataDogLoggerDescriptor, log.NewDataDogLogger)
	c.RegisterType(dataDogCountersDescriptor, count.NewDataDogCounters)

	return &c
}
