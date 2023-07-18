package persistence

import (
	"reflect"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

func ItemsToAnySlice[T any](items []T) []any {
	ln := len(items)
	result := make([]any, ln)
	for i := range items {
		result[i] = items[i]
	}
	return result
}

func GenerateObjectMapIdIfNotExists(objectMap map[string]any) {
	if id, ok := objectMap["id"]; ok {
		if reflect.ValueOf(id).IsZero() && reflect.TypeOf(id).Kind() == reflect.String {
			objectMap["id"] = keys.IdGenerator.NextLong()
		}
	}
}

func GenerateObjectIdIfNotExists[T any](obj any) T {
	if _item, ok := obj.(cdata.IStringIdentifiable); ok {
		if _item.GetId() == "" {
			_item.SetId(keys.IdGenerator.NextLong())
		}
		return _item.(T)
	}
	cpersist.GenerateObjectId(&obj)
	return obj.(T)
}

func GetObjectId[K any](obj any) (id K) {
	if _obj, ok := obj.(cdata.IIdentifiable[K]); ok {
		return _obj.GetId()
	}
	objId := cpersist.GetObjectId(obj)
	if _id, ok := objId.(K); ok {
		return _id
	}
	return
}
