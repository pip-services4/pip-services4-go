package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/connect"
)

// DefaultSqliteFactory helps creates Sqlite components by their descriptors.
//
//	see Factory
//	see SqliteConnection
type DefaultSqliteFactory struct {
	cbuild.Factory
}

// NewDefaultSqliteFactory are create a new instance of the factory.
//
//	Returns: *DefaultSqliteFactory
func NewDefaultSqliteFactory() *DefaultSqliteFactory {
	c := DefaultSqliteFactory{}

	sqliteConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "sqlite", "*", "1.0")

	c.RegisterType(sqliteConnectionDescriptor, conn.NewSqliteConnection)
	return &c
}
