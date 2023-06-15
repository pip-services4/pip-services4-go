package reflect

import (
	"fmt"
	refl "reflect"
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// ObjectReader Helper class to perform property introspection and dynamic reading.
//
// In contrast to PropertyReflector which only introspects regular objects,
// this ObjectReader is also able to handle maps and arrays.
// For maps properties are key-pairs identified by string keys,
// For arrays properties are elements identified by integer index.
//
// This class has symmetric implementation across all languages supported by
// Pip.Services toolkit and used to support dynamic data processing.
//
// Because all languages have different casing and case sensitivity rules,
// this ObjectReader treats all property names as case insensitive.
//
//	see PropertyReflector
//
//	Example:
//		myObj := MyObject{}
//
//		properties := ObjectReader.GetPropertyNames()
//		ObjectReader.HasProperty(myObj, "myProperty")
//		value := PropertyReflector.GetProperty(myObj, "myProperty")
//
//		myMap := { key1: 123, key2: "ABC" }
//		ObjectReader.HasProperty(myMap, "key1")
//		value := ObjectReader.GetProperty(myMap, "key1")
//
//		myArray := [1, 2, 3]
//		ObjectReader.HasProperty(myArrat, "0")
//		value := ObjectReader.GetProperty(myArray, "0")
var ObjectReader = &_TObjectReader{}

type _TObjectReader struct{}

// GetValue gets a real object value. If object is a wrapper, it unwraps the value behind it. Otherwise it returns the same object value.
//
//	Parameters: obj any an object to unwrap
//	Returns: any an actual (unwrapped) object value.
func (c *_TObjectReader) GetValue(obj any) any {
	if wrap, ok := obj.(IValueWrapper); ok {
		obj = wrap.InnerValue()
	}
	return obj
}

// HasProperty checks if object has a property with specified name.
// The object can be a user defined object, map or array.
// The property name correspondently must be object property, map key or array index.
//
//	Parameters:
//		- obj any an object to introspect.
//		- name string a name of the property to check.
//	Returns: bool true if the object has the property and false if it doesn't.
func (c *_TObjectReader) HasProperty(obj any, name string) bool {
	if obj == nil || name == "" {
		return false
	}

	obj = c.GetValue(obj)
	val := refl.ValueOf(obj)

	if val.Kind() == refl.Map {
		name = strings.ToLower(name)
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			key = strings.ToLower(key)
			if name == key {
				return true
			}
		}
		return false
	}

	if val.Kind() == refl.Slice || val.Kind() == refl.Array {
		index := convert.IntegerConverter.ToIntegerWithDefault(name, -1)
		return index >= 0 && index < val.Len()
	}

	return PropertyReflector.HasProperty(obj, name)
}

// GetProperty gets value of object property specified by its name.
// The object can be a user defined object, map or array.
// The property name correspondently must be object property, map key or array index.
//
//	Parameters:
//		- obj interface an object to read property from.
//		- name string a name of the property to get.
//	Returns: any the property value or null if property doesn't exist or introspection failed.
func (c *_TObjectReader) GetProperty(obj any, name string) any {
	if obj == nil || name == "" {
		return nil
	}

	obj = c.GetValue(obj)
	val := refl.ValueOf(obj)

	if val.Kind() == refl.Map {
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

	if val.Kind() == refl.Slice || val.Kind() == refl.Array {
		index := convert.IntegerConverter.ToIntegerWithDefault(name, -1)
		if index >= 0 && index < val.Len() {
			return val.Index(index).Interface()
		}
		return nil
	}

	return PropertyReflector.GetProperty(obj, name)
}

// GetPropertyNames gets names of all properties implemented in specified object.
// The object can be a user defined object, map or array.
// Returned property name correspondently are object properties, map keys or array indexes.
//
//	Parameters: obj any an objec to introspect.
//	Returns: []string a list with property names.
func (c *_TObjectReader) GetPropertyNames(obj any) []string {
	if obj == nil {
		return nil
	}

	obj = c.GetValue(obj)
	val := refl.ValueOf(obj)
	properties := make([]string, 0)

	if val.Kind() == refl.Map {
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			properties = append(properties, key)
		}
		return properties
	}

	if val.Kind() == refl.Slice || val.Kind() == refl.Array {
		strFmt := c.GetStrIndexFormat(val.Len())
		for index := 0; index < val.Len(); index++ {
			properties = append(properties, fmt.Sprintf(strFmt, index))
		}
		return properties
	}

	return PropertyReflector.GetPropertyNames(obj)
}

// GetProperties get values of all properties in specified object and returns them as a map.
// The object can be a user defined object, map or array.
// Returned properties correspondently are object properties, map key-pairs or array elements with their indexes.
//
//	Parameters: obj any an object to get properties from.
//	Returns: map[string]any a map, containing the names of the object's properties and their values.
func (c *_TObjectReader) GetProperties(obj any) map[string]any {
	if obj == nil {
		return nil
	}

	obj = c.GetValue(obj)
	val := refl.ValueOf(obj)
	values := make(map[string]any)

	if val.Kind() == refl.Map {
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			values[key] = val.MapIndex(v).Interface()
		}
		return values
	}

	if val.Kind() == refl.Slice || val.Kind() == refl.Array {
		strFmt := c.GetStrIndexFormat(val.Len())
		for index := 0; index < val.Len(); index++ {
			values[fmt.Sprintf(strFmt, index)] = val.Index(index).Interface()
		}
		return values
	}

	return PropertyReflector.GetProperties(obj)
}

func (c *_TObjectReader) GetStrIndexFormat(len int) string {
	count := 0
	for len > 0 {
		count++
		len = len / 10
	}

	return fmt.Sprintf("%%0%dd", count)
}
