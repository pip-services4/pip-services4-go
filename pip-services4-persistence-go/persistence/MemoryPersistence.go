package persistence

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/read"
	"github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/write"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// MemoryPersistence abstract persistence component that stores data in memory.
//
//	This is the most basic persistence component that is only
//	able to store data items of any type. Specific CRUD operations
//	over the data items must be implemented in child struct by
//	accessing Items property and calling Save method.
//
//	The component supports loading and saving items from another data source.
//	That allows to use it as a base struct for file and other types
//	of persistence components that cache all data in memory.
//
//	Important:
//		- this component is a thread save!
//		- if data object will implement ICloneable interface, it rises speed of execution
//	References:
//		*:logger:*:*:1.0    ILogger components to pass log messages
//	Typed params:
//		- T cdata.ICloneable[T] any type that implemented
//			ICloneable interface of getting element
//	Example:
//		type MyMemoryPersistence struct {
//			*MemoryPersistence[MyData]
//		}
//
//		func (c *MyMemoryPersistence) GetByName(ctx context.Context,
//			name string) (MyData, error) {
//			for _, v := range c.Items {
//				if v.Name == name {
//					return v
//				}
//			}
//			var defaultValue T
//			return defaultValue, nil
//		}
//
//	Implements: IReferenceable, IOpenable, ICleanable
type MemoryPersistence[T any] struct {
	Logger      *log.CompositeLogger
	Items       []T
	Loader      read.ILoader[T]
	Saver       write.ISaver[T]
	Mtx         sync.RWMutex
	opened      bool
	MaxPageSize int
	convertor   convert.IJSONEngine[T]
}

// NewMemoryPersistence creates a new instance of the MemoryPersistence
//
//	Typed params:
//		- T cdata.ICloneable[T] any type that implemented
//			ICloneable interface of getting element
//	Return *MemoryPersistence[T]
func NewMemoryPersistence[T any]() *MemoryPersistence[T] {
	c := &MemoryPersistence[T]{
		convertor: convert.NewDefaultCustomTypeJsonConvertor[T](),
	}
	c.Logger = log.NewCompositeLogger()
	c.Items = make([]T, 0, 10)
	return c
}

// SetReferences references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references refer.IReferences references to locate the component dependencies.
func (c *MemoryPersistence[T]) SetReferences(ctx context.Context, references refer.IReferences) {
	c.Logger.SetReferences(ctx, references)
}

// IsOpen checks if the component is opened.
//
//	Returns true if the component has been opened and false otherwise.
func (c *MemoryPersistence[T]) IsOpen() bool {
	c.Mtx.RLock()
	defer c.Mtx.RUnlock()
	return c.opened
}

// Open the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error or null no errors occurred.
func (c *MemoryPersistence[T]) Open(ctx context.Context) error {
	c.Mtx.Lock()
	defer c.Mtx.Unlock()

	if c.Loader == nil {
		return nil
	}

	items, err := c.Loader.Load(ctx)
	if err == nil && items != nil {
		c.Items = make([]T, len(items))
		for i, v := range items {
			c.Items[i] = c.cloneItem(v)
		}
		length := len(c.Items)
		c.Logger.Trace(ctx, "Loaded %d items", length)
	}
	c.opened = true
	return nil
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error or null no errors occurred.
func (c *MemoryPersistence[T]) Close(ctx context.Context) error {
	err := c.Save(ctx)
	c.Mtx.Lock()
	defer c.Mtx.Unlock()
	c.opened = false
	return err
}

// Save items to external data source using configured saver component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error or null for success.
func (c *MemoryPersistence[T]) Save(ctx context.Context) error {
	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	if c.Saver == nil {
		return nil
	}

	err := c.Saver.Save(ctx, c.Items)
	if err == nil {
		length := len(c.Items)
		c.Logger.Trace(ctx, "Saved %d items", length)
	}
	return err
}

// Clear component state.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error or null no errors occurred.
func (c *MemoryPersistence[T]) Clear(ctx context.Context) error {
	if err := c.Save(ctx); err != nil {
		return err
	}
	c.Mtx.Lock()
	defer c.Mtx.Unlock()

	c.Items = make([]T, 0, 5)
	c.Logger.Trace(ctx, "Cleared items")

	return nil
}

// GetPageByFilter gets a page of data items retrieved by a given filter and sorted
// according to sort parameters.
// method shall be called by a func (imp* IdentifiableMemoryPersistence)
// getPageByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- filter func(any) bool (optional) a filter function to filter items
//		- paging cdata.PagingParams (optional) paging parameters
//		- sortFunc func(a, b T) bool (optional) sorting compare function func Less (a, b T) bool
//			see sort.Interface Less function
//		- selectFunc func(in T}) (out interface{}) (optional) projection parameters
//	Return cdata.DataPage[T], error data page or error.
func (c *MemoryPersistence[T]) GetPageByFilter(ctx context.Context,
	filterFunc func(T) bool,
	paging cquery.PagingParams,
	sortFunc func(T, T) bool,
	selectFunc func(T) T) (cquery.DataPage[T], error) {

	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	items := make([]T, 0, len(c.Items))

	// Apply filtering
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				items = append(items, c.cloneItem(v))
			}
		}
	} else {
		for _, v := range c.Items {
			items = append(items, c.cloneItem(v))
		}
	}

	// Apply sorting
	if sortFunc != nil {
		localSort := sorter[T]{items: items, compFunc: sortFunc}
		sort.Sort(localSort)
	}

	// Extract a page
	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c.MaxPageSize))
	var total int64
	if paging.Total {
		total = (int64)(len(items))
	}
	if skip > 0 {
		_len := (int64)(len(items))
		if skip >= _len {
			skip = _len
		}
		items = items[skip:]
	}
	if (int64)(len(items)) >= take {
		items = items[:take]
	}

	// Get projection
	if selectFunc != nil {
		for i, v := range items {
			items[i] = selectFunc(v)
		}
	}

	c.Logger.Trace(ctx, "Retrieved %d items", len(items))

	return *cquery.NewDataPage[T](items, int(total)), nil
}

// GetListByFilter gets a list of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c * IdentifiableMemoryPersistence)
// GetListByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- filter func(T) bool (optional) a filter function to filter items
//		- sortFunc func(a, b T) bool (optional) sorting compare function
//			func Less (a, b T) bool  see sort.Interface Less function
//		- selectFunc func(in T) (out T) (optional) projection parameters
//	Returns: []T, error array of items and error
func (c *MemoryPersistence[T]) GetListByFilter(ctx context.Context,
	filterFunc func(T) bool,
	sortFunc func(T, T) bool,
	selectFunc func(T) T) ([]T, error) {

	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	// Apply filter
	items := make([]T, 0, len(c.Items))

	// Apply filtering
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				items = append(items, c.cloneItem(v))
			}
		}
	} else {
		for _, v := range c.Items {
			items = append(items, c.cloneItem(v))
		}
	}

	if len(items) == 0 {
		return nil, nil
	}

	// Apply sorting
	if sortFunc != nil {
		localSort := sorter[T]{items: items, compFunc: sortFunc}
		sort.Sort(localSort)
	}

	// Get projection
	if selectFunc != nil {
		for i, v := range items {
			items[i] = selectFunc(v)
		}
	}

	c.Logger.Trace(ctx, "Retrieved %d items", len(items))

	return items, nil
}

// GetOneRandom gets a random item from items that match to a given filter.
// This method shall be called by a func (c* IdentifiableMemoryPersistence) GetOneRandom method from child type that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- filter func(T) bool (optional) a filter function to filter items.
//	Returns: T, error random item or error.
func (c *MemoryPersistence[T]) GetOneRandom(ctx context.Context,
	filterFunc func(T) bool) (T, error) {

	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	// Apply filter
	items := make([]T, 0, len(c.Items))

	// Apply filtering
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				items = append(items, c.cloneItem(v))
			}
		}
	} else {
		for _, v := range c.Items {
			items = append(items, c.cloneItem(v))
		}
	}
	rand.Seed(time.Now().UnixNano())

	var item *T = nil
	if len(items) > 0 {
		item = &items[rand.Intn(len(items))]
	}

	if item != nil {
		c.Logger.Trace(ctx, "Retrieved a random item")
	} else {
		c.Logger.Trace(ctx, "Nothing to return as random item")

		var defaultValue T
		return defaultValue, nil
	}

	return *item, nil
}

// Create a data item.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- item T an item to be created.
//	Returns: T, error created item or error.
func (c *MemoryPersistence[T]) Create(ctx context.Context, item T) (T, error) {

	c.Mtx.Lock()

	c.Items = append(c.Items, c.cloneItem(item))

	c.Logger.Trace(ctx, "Created item")

	c.Mtx.Unlock()

	if err := c.Save(ctx); err != nil {
		return c.cloneItem(item), err
	}

	return c.cloneItem(item), nil
}

// DeleteByFilter data items that match to a given filter.
// this method shall be called by a func (c* IdentifiableMemoryPersistence)
// DeleteByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- filter  filter func(T) bool (optional) a filter function to filter items.
//	Returns: error or nil for success.
func (c *MemoryPersistence[T]) DeleteByFilter(ctx context.Context,
	filterFunc func(T) bool) error {

	c.Mtx.Lock()

	deleted := 0
	for i := 0; i < len(c.Items); {
		if filterFunc(c.Items[i]) {
			if i == len(c.Items)-1 {
				c.Items = c.Items[:i]
			} else {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}
			deleted++
		} else {
			i++
		}
	}
	c.Mtx.Unlock()

	if deleted == 0 {
		return nil
	}

	c.Logger.Trace(ctx, "Deleted %s items", deleted)

	return c.Save(ctx)
}

// GetCountByFilter gets a count of data items retrieved by a given filter.
// this method shall be called by a func (imp* IdentifiableMemoryPersistence)
// getCountByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- filter func(T) bool (optional) a filter function to filter items
//	Return int64, error data count or error.
func (c *MemoryPersistence[T]) GetCountByFilter(ctx context.Context,
	filterFunc func(T) bool) (int64, error) {

	c.Mtx.RLock()
	defer c.Mtx.RUnlock()

	var count int64

	// Apply filtering
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				count++
			}
		}
	}
	c.Logger.Trace(ctx, "Find %d items", count)
	return count, nil
}

func (c *MemoryPersistence[T]) cloneItem(item any) T {
	if cloneableItem, ok := item.(data.ICloneable[T]); ok {
		return cloneableItem.Clone()
	}

	strObject, _ := c.convertor.ToJson(item.(T))
	newItem, _ := c.convertor.FromJson(strObject)
	return newItem
}
