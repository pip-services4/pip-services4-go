package persistence

import (
	"context"
	"reflect"
	"sync"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	refl "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"

	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// IdentifiableMemoryPersistence Abstract persistence component that stores data in memory
// and implements a number of CRUD operations over data items with unique ids.
//
// In basic scenarios' child structs shall only override GetPageByFilter,
// GetListByFilter or DeleteByFilter operations with specific filter function.
// All other operations can be used out of the box.
//
// In complex scenarios' child structs can implement additional operations by
// accessing cached items via c.Items property and calling Save method
// on updates.
//
//	Important:
//		- this component is a thread save!
//		- the data items must implement IDataObject interface
//
//	see MemoryPersistence
//
//	Configuration parameters:
//		- options
//		- max_page_size maximum number of items returned in a single page (default: 100)
//	References:
//		- *:logger:*:*:1.0 (optional) ILogger components to pass log messages
//	Typed params:
//		- T cdata.IDataObject[T, K] any type that implemented
//			IDataObject interface of getting element
//		- K any type if id (key)
//	Examples:
//		type MyMemoryPersistence struct {
//			*IdentifiableMemoryPersistence[*MyData, string]
//		}
//
//		func NewMyMemoryPersistence() *MyMemoryPersistence {
//			return &MyMemoryPersistence{IdentifiableMemoryPersistence: NewIdentifiableMemoryPersistence[*MyData, string]()}
//		}
//		func (c *MyMemoryPersistence) composeFilter(filter cdata.FilterParams) func(item *MyData) bool {
//			name, _ := filter.GetAsNullableString("Name")
//			return func(item *MyData) bool {
//				if name != "" && item.Name != name {
//					return false
//				}
//				return true
//			}
//		}
//
//		func (c *MyMemoryPersistence) GetPageByFilter(ctx context.Context,
//			filter FilterParams, paging PagingParams) (page cdata.DataPage[*MyData], err error) {
//			return c.GetPageByFilter(ctx, c.composeFilter(filter), paging, nil, nil)
//		}
//
//		func f() {
//			persistence := NewMyMemoryPersistence()
//
//			item, err := persistence.Create(context.Background(), "123", &MyData{Id: "1", Name: "ABC"})
//			// ...
//			page, err := persistence.GetPageByFilter(context.Background(), *NewFilterParamsFromTuples("Name", "ABC"), nil)
//			if err != nil {
//				panic("Error can't get data")
//			}
//			data := page.Data
//			fmt.Println(data) // Result: { Id: "1", Name: "ABC" }
//			item, err = persistence.DeleteById(context.Background(), "123", "1")
//			// ...
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
//	Implements: IConfigurable, IWriter, IGetter, ISetter
type IdentifiableMemoryPersistence[T any, K any] struct {
	*MemoryPersistence[T]
	Mtx sync.RWMutex
}

const IdentifiableMemoryPersistenceConfigParamOptionsMaxPageSize = "options.max_page_size"

// NewIdentifiableMemoryPersistence creates a new empty instance of the persistence.
//
//	Typed params:
//		- T cdata.IDataObject[T, K] any type that implemented
//			IDataObject interface of getting element
//		- K any type if id (key)
//
// Returns: *IdentifiableMemoryPersistence created empty IdentifiableMemoryPersistence
func NewIdentifiableMemoryPersistence[T any, K any]() (c *IdentifiableMemoryPersistence[T, K]) {
	c = &IdentifiableMemoryPersistence[T, K]{
		MemoryPersistence: NewMemoryPersistence[T](),
	}
	c.Logger = log.NewCompositeLogger()
	c.MaxPageSize = 100
	return c
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *IdentifiableMemoryPersistence[T, K]) Configure(ctx context.Context, config *config.ConfigParams) {
	c.MaxPageSize = config.GetAsIntegerWithDefault(IdentifiableMemoryPersistenceConfigParamOptionsMaxPageSize, c.MaxPageSize)
}

// GetListByIds gets a list of data items retrieved by given unique ids.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- ids []K ids of data items to be retrieved
//	Returns: []T, error data list or error.
func (c *IdentifiableMemoryPersistence[T, K]) GetListByIds(ctx context.Context,
	ids []K) ([]T, error) {

	filter := func(item T) bool {
		itemId := c.getItemId(item)
		for _, _id := range ids {
			if c.isEqualIds(itemId, _id) {
				return true
			}
		}
		return false
	}
	return c.GetListByFilter(ctx, filter, nil, nil)
}

// GetOneById gets a data item by its unique id.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- id K an id of data item to be retrieved.
//
// Returns: T, error data item or error.
func (c *IdentifiableMemoryPersistence[T, K]) GetOneById(ctx context.Context, id K) (T, error) {

	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	for _, item := range c.Items {
		itemId := c.getItemId(item)
		if c.isEqualIds(itemId, id) {
			c.Logger.Trace(ctx, "Retrieved item %s", id)
			return c.cloneItem(item), nil
		}
	}

	c.Logger.Trace(ctx, "Cannot find item by %s", id)

	var defaultObject T
	return defaultObject, nil
}

// GetIndexById get index by "Id" field
//
//	Parameters:
//		- id K id parameter of data struct
//	Returns: index number
func (c *IdentifiableMemoryPersistence[T, K]) GetIndexById(id K) int {
	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	for i, item := range c.Items {
		if c.isEqualIds(c.getItemId(item), id) {
			return i
		}
	}
	return -1
}

// Create a data item.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- item T an item to be created.
//	Returns: T, error created item or error.
func (c *IdentifiableMemoryPersistence[T, K]) Create(ctx context.Context, item T) (T, error) {
	c.Mtx.Lock()

	newItem := c.cloneItem(item)
	if _item, ok := c.setItemId(newItem, c.getItemId(newItem)).(T); ok {
		newItem = _item
	}

	c.Items = append(c.Items, newItem)

	c.Mtx.Unlock()
	c.Logger.Trace(ctx, "Created item %s", c.getItemId(newItem))

	if err := c.Save(ctx); err != nil {
		return c.cloneItem(newItem), err
	}

	return c.cloneItem(newItem), nil
}

// Set a data item. If the data item exists it updates it,
// otherwise it creates a new data item.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- item T a item to be set.
//
// Returns: T, error updated item or error.
func (c *IdentifiableMemoryPersistence[T, K]) Set(ctx context.Context, item T) (T, error) {
	newItem := c.cloneItem(item)
	if _item, ok := c.setItemId(newItem, c.getItemId(newItem)).(T); ok {
		newItem = _item
	}

	index := c.GetIndexById(c.getItemId(item))

	c.Mtx.Lock()
	if index < 0 {
		c.Items = append(c.Items, newItem)
	} else {
		c.Items[index] = newItem
	}

	c.Mtx.Unlock()
	c.Logger.Trace(ctx, "Set item %s", c.getItemId(newItem))

	if err := c.Save(ctx); err != nil {
		return c.cloneItem(newItem), err
	}

	return c.cloneItem(newItem), nil
}

// Update a data item.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- item T an item to be updated.
//
// Returns: T, error updated item or error.
func (c *IdentifiableMemoryPersistence[T, K]) Update(ctx context.Context, item T) (T, error) {
	var defaultObject T

	index := c.GetIndexById(c.getItemId(item))
	if index < 0 {
		c.Logger.Trace(ctx, "Item %s was not found", c.getItemId(item))
		return defaultObject, nil
	}
	newItem := c.cloneItem(item)

	c.Mtx.Lock()
	c.Items[index] = newItem
	c.Mtx.Unlock()

	c.Logger.Trace(ctx, "Updated item %s", c.getItemId(item))

	if err := c.Save(ctx); err != nil {
		return c.cloneItem(newItem), err
	}

	return c.cloneItem(newItem), nil
}

// UpdatePartially only few selected fields in a data item.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- id K an id of data item to be updated.
//		- data  cdata.AnyValueMap a map with fields to be updated.
//
// Returns: T, error updated item or error.
func (c *IdentifiableMemoryPersistence[T, K]) UpdatePartially(ctx context.Context,
	id K, data cdata.AnyValueMap) (T, error) {

	var defaultObject T

	index := c.GetIndexById(id)
	if index < 0 {
		c.Logger.Trace(ctx, "Item %s was not found", id)
		return defaultObject, nil
	}

	c.Mtx.Lock()

	newItem := c.cloneItem(c.Items[index])

	if reflect.ValueOf(newItem).Kind() == reflect.Map {
		refl.ObjectWriter.SetProperties(newItem, data.Value())
	} else {
		var intPointer any = newItem
		if reflect.TypeOf(newItem).Kind() != reflect.Pointer {
			objPointer := reflect.New(reflect.TypeOf(newItem))
			objPointer.Elem().Set(reflect.ValueOf(newItem))
			intPointer = objPointer.Interface()
		}
		refl.ObjectWriter.SetProperties(intPointer, data.Value())
		if _newItem, ok := reflect.ValueOf(intPointer).Elem().Interface().(T); ok {
			newItem = _newItem
		}
	}

	c.Items[index] = newItem

	c.Mtx.Unlock()
	c.Logger.Trace(ctx, "Partially updated item %s", id)

	if err := c.Save(ctx); err != nil {
		return c.cloneItem(newItem), err
	}

	return c.cloneItem(newItem), nil
}

// DeleteById a data item by it's unique id.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- id K an id of the item to be deleted
//	Returns: T, error deleted item or error.
func (c *IdentifiableMemoryPersistence[T, K]) DeleteById(ctx context.Context, id K) (T, error) {

	var defaultObject T

	index := c.GetIndexById(id)
	if index < 0 {
		c.Logger.Trace(ctx, "Item %s was not found", id)
		return defaultObject, nil
	}

	c.Mtx.Lock()

	oldItem := c.Items[index]
	if index == len(c.Items) {
		c.Items = c.Items[:index-1]
	} else {
		c.Items = append(c.Items[:index], c.Items[index+1:]...)
	}

	c.Mtx.Unlock()

	c.Logger.Trace(ctx, "Deleted item by %s", id)

	if err := c.Save(ctx); err != nil {
		return oldItem, err
	}
	return oldItem, nil
}

// DeleteByIds multiple data items by their unique ids.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- ids []K ids of data items to be deleted.
//	Returns: error or null for success.
func (c *IdentifiableMemoryPersistence[T, K]) DeleteByIds(ctx context.Context, ids []K) error {
	filterFunc := func(item T) bool {
		itemId := c.getItemId(item)
		for _, id := range ids {
			if c.isEqualIds(itemId, id) {
				return true
			}
		}
		return false
	}

	return c.DeleteByFilter(ctx, filterFunc)
}

func (c *IdentifiableMemoryPersistence[T, K]) isEqualIds(idA, idB any) bool {
	return CompareValues(idA, idB)
}

func (c *IdentifiableMemoryPersistence[T, K]) getItemId(item any) K {
	if _item, ok := item.(data.IIdentifiable[K]); ok {
		return _item.GetId()
	}

	if _id, ok := GetObjectId(item).(K); ok {
		return _id
	}

	var defaultValue K
	return defaultValue
}

func (c *IdentifiableMemoryPersistence[T, K]) setItemId(item any, id any) any {
	newId := id
	if c.isEmptyId(id) {
		newId = keys.IdGenerator.NextLong()
	}
	SetObjectId(&item, newId)
	return item
}

func (c *IdentifiableMemoryPersistence[T, K]) isEmptyId(id any) bool {
	if _id, ok := id.(data.IIdentifier[K]); ok {
		return _id.Empty()
	}

	return reflect.ValueOf(id).IsZero()
}
