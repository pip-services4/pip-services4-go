package reflect

import (
	"strings"

	"github.com/pip-services4/pip-services4-commons-go/convert"
)

// RecursiveObjectReader Helper class to perform property introspection and dynamic reading.
//
// It is similar to ObjectReader but reads properties recursively through the entire object graph.
// Nested property names are defined using dot notation as "object.subobject.property"
var RecursiveObjectReader = &TRecursiveObjectReader{}

type TRecursiveObjectReader struct{}

func (c *TRecursiveObjectReader) performHasProperty(obj any, names []string, nameIndex int) bool {
	if nameIndex < len(names)-1 {
		value := ObjectReader.GetProperty(obj, names[nameIndex])
		if value != nil {
			return c.performHasProperty(value, names, nameIndex+1)
		}
		return false
	}

	return ObjectReader.HasProperty(obj, names[nameIndex])
}

// HasProperty checks recursively if object or its subobjects has a property with specified name.
// The object can be a user defined object, map or array.
// The property name correspondently must be object property, map key or array index.
//
//	Parameters:
//		- obj any an object to introspect.
//		- name string a name of the property to check.
//	Returns: boolean true if the object has the property and false if it doesn't.
func (c *TRecursiveObjectReader) HasProperty(obj any, name string) bool {
	if obj == nil || name == "" {
		return false
	}

	names := strings.Split(name, ".")
	if len(names) == 0 {
		return false
	}

	return c.performHasProperty(obj, names, 0)
}

func (c *TRecursiveObjectReader) performGetProperty(obj any, names []string, nameIndex int) any {
	if nameIndex < len(names)-1 {
		value := ObjectReader.GetProperty(obj, names[nameIndex])
		if value != nil {
			return c.performGetProperty(value, names, nameIndex+1)
		}

		return nil
	}

	return ObjectReader.GetProperty(obj, names[nameIndex])
}

// GetProperty recursively gets value of object or its subobjects property specified by its name.
// The object can be a user defined object, map or array.
// The property name correspondently must be object property, map key or array index.
//
//	Parameters:
//		- obj any an object to read property from.
//		- name string a name of the property to get.
//	Returns: any the property value or null if property doesn't exist or introspection failed.
func (c *TRecursiveObjectReader) GetProperty(obj any, name string) any {
	if obj == nil || name == "" {
		return nil
	}

	names := strings.Split(name, ".")
	if len(names) == 0 {
		return nil
	}

	return c.performGetProperty(obj, names, 0)
}

func (c *TRecursiveObjectReader) isSimpleValue(value any) bool {
	code := convert.TypeConverter.ToTypeCode(value)
	return code != convert.Array && code != convert.Map && code != convert.Object
}

func (c *TRecursiveObjectReader) contains(values []any, obj any) bool {
	for _, value := range values {
		if value == obj {
			return true
		}
	}
	return false
}

func (c *TRecursiveObjectReader) performGetPropertyNames(obj any, path string,
	result []string, cycleDetect []any) []string {
	values := ObjectReader.GetProperties(obj)

	if len(values) != 0 && len(cycleDetect) < 100 {
		savedDetect := cycleDetect
		cycleDetect = append(cycleDetect, &obj)

		for key, value := range values {
			// Prevent cycles
			if c.contains(cycleDetect, &value) {
				continue
			}

			if path != "" {
				key = path + "." + key
			}

			// Add simple values directly
			if c.isSimpleValue(value) {
				result = append(result, key)
			} else {
				// Recursively go to elements
				result = c.performGetPropertyNames(value, key, result, cycleDetect)
			}
		}

		cycleDetect = savedDetect
	} else {
		if path != "" {
			result = append(result, path)
		}
	}

	return result
}

// GetPropertyNames Recursively gets names of all properties implemented in specified object and its subobjects.
// The object can be a user defined object, map or array.
// Returned property name correspondently are object properties, map keys or array indexes.
//
//	Parameters: obj any an object to introspect.
//	Returns: []string a list with property names.
func (c *TRecursiveObjectReader) GetPropertyNames(obj any) []string {
	propertyNames := make([]string, 0)

	if obj == nil {
		return propertyNames
	}

	cycleDetect := make([]any, 0)
	propertyNames = c.performGetPropertyNames(obj, "", propertyNames, cycleDetect)
	return propertyNames
}

func (c *TRecursiveObjectReader) performGetProperties(obj any, path string, result map[string]any,
	cycleDetect []any) map[string]any {

	values := ObjectReader.GetProperties(obj)

	if len(values) != 0 && len(cycleDetect) < 100 {
		savedDetect := cycleDetect
		cycleDetect = append(cycleDetect, &obj)

		for key, value := range values {
			// Prevent cycles
			if c.contains(cycleDetect, &value) {
				continue
			}

			if path != "" {
				key = path + "." + key
			}

			// Add simple values directly
			if c.isSimpleValue(value) {
				result[key] = value
			} else {
				// Recursively go to elements
				result = c.performGetProperties(value, key, result, cycleDetect)
			}
		}

		cycleDetect = savedDetect
	} else {
		if path != "" {
			result[path] = obj
		}
	}
	return result
}

// GetProperties get values of all properties in specified object and its subobjects and returns them as a map.
// The object can be a user defined object, map or array.
// Returned properties correspondently are object properties, map key-pairs or array elements with their indexes.
//
//	Parameters: obj any an object to get properties from.
//	Returns: map[string]any a map, containing the names of the object's properties and their values.
func (c *TRecursiveObjectReader) GetProperties(obj any) map[string]any {
	properties := map[string]any{}

	if obj == nil {
		return properties
	}

	cycleDetect := []any{}
	properties = c.performGetProperties(obj, "", properties, cycleDetect)
	return properties
}
