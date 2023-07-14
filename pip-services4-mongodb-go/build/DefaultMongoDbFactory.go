package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-mongodb-go/connect"
)

// DefaultMongoDbFactory helps creates MongoDb components by their descriptors.
//
//	see Factory
//	see MongoDbConnection
type DefaultMongoDbFactory struct {
	cbuild.Factory
}

// NewDefaultMongoDbFactory are create a new instance of the factory.
//
//	Returns: *DefaultMongoDbFactory
func NewDefaultMongoDbFactory() *DefaultMongoDbFactory {
	c := DefaultMongoDbFactory{}

	mongoDbConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "mongodb", "*", "1.0")

	c.RegisterType(mongoDbConnectionDescriptor, conn.NewMongoDbConnection)
	return &c
}
