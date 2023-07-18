package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/connect"
)

// DefaultPostgresFactory creates Postgres components by their descriptors.
//
//	see Factory
//	see PostgresConnection
type DefaultPostgresFactory struct {
	*cbuild.Factory
}

// Create a new instance of the factory.
func NewDefaultPostgresFactory() *DefaultPostgresFactory {

	c := &DefaultPostgresFactory{}
	c.Factory = cbuild.NewFactory()

	postgresConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "postgres", "*", "1.0")
	c.RegisterType(postgresConnectionDescriptor, conn.NewPostgresConnection)

	return c
}
