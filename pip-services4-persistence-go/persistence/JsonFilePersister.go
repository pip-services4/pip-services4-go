package persistence

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
)

// JsonFilePersister is a persistence component that loads and saves data from/to flat file.
// It is used by FilePersistence, but can be useful on its own.
//
//	Important: this component is not thread save!
//	Configuration parameters:
//		- path to the file where data is stored
//	Typed params:
//		- T any type
//	Example:
//		persister := NewJsonFilePersister[MyData]("./data/data.json")
//		err := persister.Save(context.Background(), "123", []string{"A", "B", "C"})
//		if err == nil {
//			items, err := persister.Load("123")
//			if err == nil {
//				fmt.Println(items) // Result: ["A", "B", "C"]
//			}
//		}
//	Implements: ILoader, ISaver, IConfigurable
type JsonFilePersister[T any] struct {
	path      string
	convertor convert.IJSONEngine[[]T]
}

const ConfigParamPath = "path"

// NewJsonFilePersister creates a new instance of the persistence.
//
//	Typed params:
//		- T any type
//	Parameters: path string (optional) a path to the file where data is stored.
func NewJsonFilePersister[T any](path string) *JsonFilePersister[T] {
	return &JsonFilePersister[T]{
		path:      path,
		convertor: convert.NewDefaultCustomTypeJsonConvertor[[]T](),
	}
}

// Path gets the file path where data is stored.
//
//	Returns: the file path where data is stored.
func (c *JsonFilePersister[T]) Path() string {
	return c.path
}

// SetPath the file path where data is stored.
//
//	Parameters:
//		- value string the file path where data is stored.
func (c *JsonFilePersister[T]) SetPath(value string) {
	c.path = value
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- config: ConfigParams configuration parameters to be set.
func (c *JsonFilePersister[T]) Configure(ctx context.Context, config *config.ConfigParams) {
	c.path = config.GetAsStringWithDefault(ConfigParamPath, c.path)
}

// Load data items from external JSON file.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//
// Returns: []T, error loaded items or error.
func (c *JsonFilePersister[T]) Load(ctx context.Context) ([]T, error) {
	if c.path == "" {
		return nil, errors.NewConfigError("", "NO_PATH", "Data file path is not set")
	}

	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return nil, err
	}

	jsonStr, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, errors.NewFileError(
			cctx.GetTraceId(ctx),
			"READ_FAILED",
			"Failed to read data file: "+c.path).
			WithCause(err)
	}

	if len(jsonStr) == 0 {
		return nil, nil
	}

	if list, err := c.convertor.FromJson(string(jsonStr)); err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

// Save given data items to external JSON file.
//
//		Parameters:
//			- ctx context.Context execution context to trace execution through call chain.
//			- items []T list of data items to save
//	 Returns: error or nil for success.
func (c *JsonFilePersister[T]) Save(ctx context.Context, items []T) error {
	json, err := c.convertor.ToJson(items)
	if err != nil {
		err := errors.NewInternalError(cctx.GetTraceId(ctx), "CAN'T_CONVERT", "Failed convert to JSON")
		return err
	}
	if err := ioutil.WriteFile(c.path, ([]byte)(json), 0777); err != nil {
		return errors.NewFileError(
			cctx.GetTraceId(ctx),
			"WRITE_FAILED",
			"Failed to write data file: "+c.path).
			WithCause(err)
	}
	return nil
}
