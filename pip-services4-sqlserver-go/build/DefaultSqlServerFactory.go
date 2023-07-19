package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/connect"
)

// DefaultSqlServerFactory creates SqlServer components by their descriptors.
//
//	see Factory
//	see SqlServerConnection
type DefaultSqlServerFactory struct {
	*cbuild.Factory
}

// Create a new instance of the factory.
func NewDefaultSqlServerFactory() *DefaultSqlServerFactory {

	c := &DefaultSqlServerFactory{}
	c.Factory = cbuild.NewFactory()

	sqlserverConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "sqlserver", "*", "1.0")
	c.RegisterType(sqlserverConnectionDescriptor, conn.NewSqlServerConnection)

	return c
}
