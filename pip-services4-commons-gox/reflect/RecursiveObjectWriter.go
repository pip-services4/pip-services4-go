package reflect

import (
	"strings"

	"github.com/pip-services4/pip-services4-commons-go/convert"
)

// RecursiveObjectWriter Helper class to perform property introspection and dynamic writing.
//
// It is similar to ObjectWriter but writes properties recursively through the entire object graph.
// Nested property names are defined using dot notation as "object.subobject.property"
var RecursiveObjectWriter = &_TRecursiveObjectWriter{}

type _TRecursiveObjectWriter struct{}

func (c *_TRecursiveObjectWriter) createProperty(obj any, names []string, nameIndex int) any {
	// Todo: Complete implementation
	// If next field is index then create an array
	subField := ""
	if len(names) > nameIndex+1 {
		subField = names[nameIndex+1]
	}
	if _, ok := convert.IntegerConverter.ToNullableInteger(subField); ok {
		return make([]any, 0)
	}

	// Else create a dictionary
	return make(map[string]any)
}

func (c *_TRecursiveObjectWriter) performSetProperty(obj any, names []string, nameIndex int, value any) {
	if nameIndex < len(names)-1 {
		subObj := ObjectReader.GetProperty(obj, names[nameIndex])
		if subObj != nil {
			c.performSetProperty(subObj, names, nameIndex+1, value)
		} else {
			subObj = c.createProperty(obj, names, nameIndex)
			if subObj != nil {
				c.performSetProperty(subObj, names, nameIndex+1, value)
				ObjectWriter.SetProperty(obj, names[nameIndex], subObj)
			}
		}
	} else {
		ObjectWriter.SetProperty(obj, names[nameIndex], value)
	}
}

// SetProperty recursively sets value of object and its subobjects property specified by its name.
// The object can be a user defined object, map or array.
// The property name correspondently must be object property, map key or array index.
// If the property does not exist or introspection fails this method doesn't do anything and doesn't any throw errors.
//
//	Parameters:
//		- obj any an object to write property to.
//		- name string a name of the property to set.
//		- value any a new value for the property to set.
func (c *_TRecursiveObjectWriter) SetProperty(obj any, name string, value any) {
	if obj == nil || name == "" {
		return
	}

	names := strings.Split(name, ".")
	if len(names) == 0 {
		return
	}

	c.performSetProperty(obj, names, 0, value)
}

// SetProperties recursively sets values of some (all) object and its subobjects properties.
// The object can be a user defined object, map or array.
// Property values correspondently are object properties, map key-pairs or array elements with their indexes.
// If some properties do not exist or introspection fails they are just silently skipped and no errors thrown.
//
//	see SetProperty
//	Parameters:
//		- obj any an object to write properties to.
//		- values map[string]any a map, containing property names and their values.
func (c *_TRecursiveObjectWriter) SetProperties(obj any, values map[string]any) {
	if len(values) == 0 {
		return
	}

	for key, value := range values {
		c.SetProperty(obj, key, value)
	}
}

// CopyProperties copies content of one object to another object by recursively reading all
// properties from source object and then recursively writing them to destination object.
//
//	Parameters:
//		- dest any a destination object to write properties to.
//		- src any a source object to read properties from
func (c *_TRecursiveObjectWriter) CopyProperties(dest any, src any) {
	if dest == nil || src == nil {
		return
	}

	values := RecursiveObjectReader.GetProperties(src)
	c.SetProperties(dest, values)
}
