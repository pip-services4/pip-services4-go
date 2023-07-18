package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/connect"
)

// DefaultMySqlFactory creates MySql components by their descriptors.
//
//	see Factory
//	see MySqlConnection
type DefaultMySqlFactory struct {
	*cbuild.Factory
}

// Create a new instance of the factory.
func NewDefaultMySqlFactory() *DefaultMySqlFactory {

	c := &DefaultMySqlFactory{}
	c.Factory = cbuild.NewFactory()

	mysqlConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "mysql", "*", "1.0")
	c.RegisterType(mysqlConnectionDescriptor, conn.NewMySqlConnection)

	return c
}
