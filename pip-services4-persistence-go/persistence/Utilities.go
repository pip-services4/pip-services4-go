package persistence

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jinzhu/copier"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	refl "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

func toFieldType(obj any) reflect.Type {
	// Unwrap value
	wrap, ok := obj.(refl.IValueWrapper)
	if ok {
		obj = wrap.InnerValue()
	}

	// Move from pointer to real struct
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

func matchField(field reflect.StructField, name string) bool {
	// Field must be public and match to name as case insensitive
	r, _ := utf8.DecodeRuneInString(field.Name)
	return unicode.IsUpper(r) &&
		strings.ToLower(field.Name) == strings.ToLower(name)
}

func getValue(obj any) any {
	wrap, ok := obj.(refl.IValueWrapper)
	if ok {
		obj = wrap.InnerValue()
	}

	return obj
}

// GetProperty value of object property specified by its name.
//
//	Parameters:
//		- obj any an object to read property from.
//		- name string a name of the property to get.
//	Returns: any the property value or null if property doesn't exist or introspection failed.
func GetProperty(obj any, name string) any {
	if obj == nil || name == "" {
		return nil
	}

	obj = getValue(obj)
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Map {
		name = strings.ToLower(name)
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			key = strings.ToLower(key)
			if name == key {
				return val.MapIndex(v).Interface()
			}
		}
		return nil
	}

	defer func() {
		// Do nothing and return nil
		recover()
	}()

	fieldType := toFieldType(obj)
	if fieldType.Kind() == reflect.Struct {
		return getPropertyRecursive(fieldType, obj, name)
	}

	return nil
}

func getPropertyRecursive(fieldType reflect.Type, obj any, name string) any {
	for index := 0; index < fieldType.NumField(); index++ {
		field := fieldType.Field(index)
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		switch field.Type.Kind() {
		default:
			if matchField(field, name) {
				return val.Field(index).Interface()
			}
		case reflect.Struct:
			if item := getPropertyRecursive(field.Type, val.Field(index).Interface(), name); item != nil {
				return item
			}
		}
	}

	return nil
}

// SetProperty value of object property specified by its name.
// If the property does not exist or introspection fails this method doesn't do anything and doesn't any throw errors.
//
//	Parameters:
//		- obj any an object to write property to.
//		- name string a name of the property to set.
//		- value any a new value for the property to set.
func SetProperty(obj any, name string, value any) {
	if obj == nil || name == "" {
		return
	}

	obj = getValue(obj)
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Map {
		name = strings.ToLower(name)
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			key = strings.ToLower(key)
			if name == key {
				val.SetMapIndex(v, reflect.ValueOf(value))
				return
			}
		}
		val.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))
		return
	}

	defer func() {
		// Do nothing and return nil
		if err := recover(); err != nil {
			fmt.Printf("Error while set property %v", err)
		}
	}()

	fieldType := toFieldType(obj)
	if fieldType.Kind() == reflect.Struct {
		setPropertyRecursive(fieldType, obj, name, value)
	}

}

func setPropertyRecursive(fieldType reflect.Type, obj any, name string, value any) {
	for index := 0; index < fieldType.NumField(); index++ {
		field := fieldType.Field(index)
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		switch field.Type.Kind() {
		default:
			if matchField(field, name) {
				val.Field(index).Set(reflect.ValueOf(value))
				return
			}
		case reflect.Struct:
			setPropertyRecursive(field.Type, val.Field(index).Addr().Interface(), name, value)
		}
	}
}

// GetObjectId value
//
//	Parameters:
//		- item any an object to read property from.
//	Returns: any the property value or nil if property doesn't exist or introspection failed.
func GetObjectId(item any) any {
	return GetProperty(item, "Id")
}

// SetObjectId is set object Id value
//
//	Parameters:
//		- item *any a pointer on object to set id property
//		- id any id value for set
//	Results: saved in input object
func SetObjectId(item *any, id any) {
	value := *item
	var isPointer bool
	if reflect.ValueOf(value).Kind() == reflect.Map {
		//refl.ObjectWriter.SetProperty(value, "Id", id)
		SetProperty(value, "Id", id)
	} else {
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			value = reflect.ValueOf(value).Elem().Interface()
			isPointer = true
		}
		typePointer := reflect.New(reflect.TypeOf(value))
		typePointer.Elem().Set(reflect.ValueOf(value))
		typeInterface := typePointer.Interface()
		//refl.ObjectWriter.SetProperty(typeInterface, "Id", id)
		SetProperty(typeInterface, "Id", id)

		if isPointer {
			*item = reflect.ValueOf(typeInterface).Interface()
		} else {
			*item = reflect.ValueOf(typeInterface).Elem().Interface()
		}
	}
}

// GenerateObjectId is generates a new id value when it's empty
//
//	Parameters:
//		- item *any a pointer on object to set id property
//	Results: saved in input object
func GenerateObjectId(item *any) {
	value := *item
	idField := GetProperty(value, "Id")
	if idField != nil {
		if reflect.ValueOf(idField).IsZero() {
			SetObjectId(item, cdata.IdGenerator.NextLong())
		}
	} else {
		panic("'Id' or 'ID' field doesn't exist")
	}

}

// CloneObject is clones object function
//
//	Parameters:
//		- item any an object to clone
//	Return any copy of input item
func CloneObject(item any, proto reflect.Type) any {
	var dest any
	var src = item

	if reflect.ValueOf(src).Kind() == reflect.Map {
		itemType := reflect.TypeOf(src)
		mapType := reflect.MapOf(itemType.Key(), itemType.Elem())
		newMap := reflect.MakeMap(mapType)
		dest = newMap.Interface()
		err := copier.CopyWithOption(&dest, src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}

	} else {
		var destPtr reflect.Value
		if proto.Kind() == reflect.Ptr {
			destPtr = reflect.New(proto.Elem())
		} else {
			destPtr = reflect.New(proto)
		}
		if reflect.TypeOf(src).Kind() == reflect.Ptr {
			src = reflect.ValueOf(src).Elem().Interface()
		}
		err := copier.CopyWithOption(destPtr.Interface(), src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}

		dest = destPtr.Elem().Interface()
	}
	return dest
}

// CloneObjectForResult is clones object for result function
//
//	Parameters:
//		- item any an object to clone
//		- proto reflect.Type of returned value, need for detect object or pointer returned type
//	Returns: any copy of input item
func CloneObjectForResult(src any, proto reflect.Type) any {
	var dest any

	if reflect.ValueOf(src).Kind() == reflect.Map {
		itemType := reflect.TypeOf(src)
		mapType := reflect.MapOf(itemType.Key(), itemType.Elem())
		newMap := reflect.MakeMap(mapType)
		dest = newMap.Interface()
		err := copier.CopyWithOption(&dest, src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}
	} else {
		var destPtr reflect.Value
		if proto.Kind() == reflect.Ptr {
			destPtr = reflect.New(proto.Elem())
		} else {
			destPtr = reflect.New(proto)
		}
		err := copier.CopyWithOption(destPtr.Interface(), src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}
		// make pointer on clone object, if proto is ptr
		dest = destPtr.Elem().Interface()
		if proto.Kind() == reflect.Ptr {
			dest = destPtr.Interface()
		}
	}

	return dest
}

// CompareValues are ompares two values
//
//	Parameters:
//		- value1 any an object one for compare
//		- value2 any an object two for compare
//	Returns: bool true if value1 equal value2 and false otherwise
func CompareValues(value1 any, value2 any) bool {
	// Todo: Implement proper comparison
	return value1 == value2
}

// Convert methods

// FromIds method convert ids string array to array of any object
//
//	Parameters:
//		- ids - []string array of ids
//	Returns: []any array of ids
func FromIds(ids []string) []any {
	result := make([]any, len(ids))
	for i, v := range ids {
		result[i] = v
	}
	return result
}

// ToPublicMap method convert any object to map[string]any
//
//	Parameters:
//		- value any input object to convert
//	Returns: map[string]any converted object to map
func ToPublicMap(value any) map[string]any {
	if value != nil {
		result, _ := value.(map[string]any)
		return result
	}
	return nil
}

// ToPublicArray method convert array of any object to array of map[string]any
//
//	Parameters:
//		- value []any input object to convert
//	Returns: []map[string]any converted map array
func ToPublicArray(values []any) []map[string]any {
	if values == nil {
		return nil
	}

	result := make([]map[string]any, len(values))
	for i, v := range values {
		result[i] = ToPublicMap(v)
	}
	return result
}
