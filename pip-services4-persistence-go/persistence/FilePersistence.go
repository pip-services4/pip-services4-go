package persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// FilePersistence is an abstract persistence component that stores data in flat files
// and caches them in memory.
//
// FilePersistence is the most basic persistence component that is only
// able to store data items of any type. Specific CRUD operations
// over the data items must be implemented in child structs by
// accessing fp._items property and calling Save method.
//
//	see MemoryPersistence
//	see JsonFilePersister
//
//	Configuration parameters:
//		- path to the file where data is stored
//	References:
//		- *:logger:*:*:1.0  (optional) ILogger components to pass log messages
//	Typed params:
//		- T cdata.ICloneable[T] any type that implemented
//			ICloneable interface of getting element
//
//	Example:
//		type MyJsonFilePersistence struct {
//			*FilePersistence[*MyData]
//		}
//
//		func NewMyJsonFilePersistence(path string) *MyJsonFilePersistence {
//			return &MyJsonFilePersistence{
//				FilePersistence: NewFilePersistence(NewJsonFilePersister[*MyData](path)),
//			}
//		}
//
//		func (c *MyJsonFilePersistence) GetByName(ctx context.Context,
//			name string) (*MyData, error) {
//			for _, v := range c.Items {
//				if v.Name == name {
//					return v, nil
//				}
//			}
//
//			var defaultValue *MyData
//			return defaultValue, nil
//		}
//
//		func (c *MyData) Clone() *MyData {
//			return &MyData{Id: c.Id, Name: c.Name}
//		}
//
//		type MyData struct {
//			Id   string
//			Name string
//		}
//
//	Extends: MemoryPersistence
//	Implements: IConfigurable
type FilePersistence[T cdata.ICloneable[T]] struct {
	*MemoryPersistence[T]
	Persister *JsonFilePersister[T]
}

// NewFilePersistence creates a new instance of the persistence.
//
//	Parameters:
//		- persister (optional) a persister component that loads and saves data from/to flat file.
//	Typed params:
//		- T cdata.ICloneable[T] any type that implemented
//			ICloneable interface of getting element
//
// Returns: *FilePersistence[T] pointer on new FilePersistence instance
func NewFilePersistence[T cdata.ICloneable[T]](persister *JsonFilePersister[T]) *FilePersistence[T] {
	c := &FilePersistence[T]{}
	c.MemoryPersistence = NewMemoryPersistence[T]()
	if persister == nil {
		persister = NewJsonFilePersister[T]("")
	}
	c.Loader = persister
	c.Saver = persister
	c.Persister = persister
	return c
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *FilePersistence[T]) Configure(ctx context.Context, conf *config.ConfigParams) {
	c.Persister.Configure(ctx, conf)
}
