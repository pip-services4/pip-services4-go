package exec

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// Parameters contains map with execution parameters.
// In general, this map may contain non-serializable values.
// And in contrast with other maps, its getters and setters support dot
// notation and able to access properties in the entire object graph.
// This class is often use to pass execution and notification arguments,
// and parameterize classes before execution.
type Parameters struct {
	*data.AnyValueMap
}

// NewEmptyParameters creates a new instance of the map and assigns its value.
// Returns: *Parameters
func NewEmptyParameters() *Parameters {
	c := &Parameters{}
	c.AnyValueMap = data.InheritAnyValueMap(c)
	return c
}

// NewParameters creates a new instance of the map and assigns its value.
//
//	Parameters: values map[string]any
//	Returns: *Parameters
func NewParameters(values map[string]any) *Parameters {
	c := &Parameters{}
	c.AnyValueMap = data.InheritAnyValueMap(c)
	c.Append(values)
	return c
}

// NewParametersFromValue creates a new Parameters object filled with key-value pairs from specified object.
//
//	Parameters: value any an object with key-value pairs used to initialize a new Parameters.
//	Returns: *Parameters a new Parameters object.
func NewParametersFromValue(value any) *Parameters {
	result := NewEmptyParameters()
	result.SetAsSingleObject(value)
	return result
}

// NewParametersFromTuples creates a new Parameters object filled with provided key-value pairs called tuples.
// Tuples parameters contain a sequence of key1, value1, key2, value2, ... pairs.
//
//	see NewParametersFromTuplesArray
//	Parameters: tuples ...any the tuples to fill a new Parameters object.
//	Returns: *Parameters a new Parameters object.
func NewParametersFromTuples(tuples ...any) *Parameters {
	return NewParametersFromTuplesArray(tuples)
}

// NewParametersFromTuplesArray creates a new AnyValueMap from a list of key-value pairs called tuples.
// The method is similar to fromTuples but tuples are passed as array instead of parameters.
//
//	Parameters: tuples []any a list of values where odd elements
//		are keys and the following even elements are values
//	Returns: *Parameters a newly created Parameters.
func NewParametersFromTuplesArray(tuples []any) *Parameters {
	result := NewEmptyParameters()
	if len(tuples) == 0 {
		return result
	}

	for index := 0; index < len(tuples); index = index + 2 {
		if index+1 >= len(tuples) {
			break
		}

		name := convert.StringConverter.ToString(tuples[index])
		value := tuples[index+1]

		result.SetAsObject(name, value)
	}

	return result
}

// NewParametersFromMaps creates a new Parameters by merging two or more maps.
// Maps defined later in the list override values from previously defined maps.
//
//	Parameters: maps ...map[string]any an array of maps to be merged
//	Returns: *Parameters a newly created Parameters.
func NewParametersFromMaps(maps ...map[string]any) *Parameters {
	result := NewEmptyParameters()
	for index := 0; index < len(maps); index++ {
		result.Append(maps[index])
	}
	return result
}

// NewParametersFromConfig creates new Parameters from ConfigMap object.
//
//	see ConfigParams
//	Parameters: config: *config.ConfigParams a ConfigParams that contain parameters.
//	Returns: Parameters a new Parameters object.
func NewParametersFromConfig(config *config.ConfigParams) *Parameters {
	result := NewEmptyParameters()
	if values, ok := config.InnerValue().(map[string]string); ok {
		for key, value := range values {
			result.Put(key, value)
		}
	}
	return result
}

// Get gets a map element specified by its key.
// The key can be defined using dot notation and allows to recursively access elements of elements.
//
//	Parameters:  key string a key of the element to get.
//	Returns any the value of the map element.
func (c *Parameters) Get(key string) (any, bool) {
	if key == "" {
		return nil, false
	}

	if strings.Index(key, ".") > 0 {
		return reflect.RecursiveObjectReader.GetProperty(c.InnerValue(), key), true
	}

	return c.AnyValueMap.Get(key)
}

// Put puts a new value into map element specified by its key.
// The key can be defined using dot notation and allows to recursively access elements of elements.
//
//	Parameters:
//		- key string a key of the element to put.
//		- value any a new value for map element.
func (c *Parameters) Put(key string, value any) {
	if strings.Trim(key, "") == "" {
		return
	}
	if strings.Index(key, ".") > 0 {
		reflect.RecursiveObjectWriter.SetProperty(c.InnerValue(), key, value)
	} else {
		c.AnyValueMap.Put(key, value)
	}
}

// Remove removes a map element specified by its key
//
//	Parameters: key string a key of the element to remove.
func (c *Parameters) Remove(key string) {
	// Todo: Make this method recursive
	c.AnyValueMap.Remove(key)
}

// Contains checks if this map contains an element with specified key.
// The key can be defined using dot notation and allows to recursively access elements of elements.
//
//	Parameters: key string a key to be checked
//	Returns: bool true if this map contains the key or false otherwise.
func (c *Parameters) Contains(key string) bool {
	return reflect.RecursiveObjectReader.HasProperty(c.InnerValue(), key)
}

// GetAsNullableParameters converts map element into an Parameters or returns nil if conversion is not possible.
//
//	Parameters: key string a key of element to get.
//	Returns: *Parameters value of the element or nil if conversion is not supported.
func (c *Parameters) GetAsNullableParameters(key string) (*Parameters, bool) {
	if value, ok := c.GetAsNullableMap(key); ok {
		return NewParametersFromValue(value), true
	}
	return nil, false
}

// GetAsParameters converts map element into a Parameters or returns empty Parameters if conversion is not possible.
//
//	Parameters: key string a key of element to get.
//	Returns: *Parameters value of the element or empty Parameters if conversion is not supported.
func (c *Parameters) GetAsParameters(key string) *Parameters {
	return NewParametersFromValue(c.GetAsMap(key))
}

// GetAsParametersWithDefault converts map element into a Parameters or returns
// default value if conversion is not possible.
//
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue *Parameters the default value
//	Returns: *Parameters value of the element or default value if conversion is not supported.
func (c *Parameters) GetAsParametersWithDefault(key string, defaultValue *Parameters) *Parameters {
	if result, ok := c.GetAsNullableParameters(key); ok {
		return result
	}
	return defaultValue
}

// Clone creates a binary clone of this object.
//
//	Returns: any a clone of this object.
func (c *Parameters) Clone() *Parameters {
	return NewParametersFromValue(c.GetAsSingleObject())
}

// Override overrides parameters with new values from specified Parameters and returns a new Parameters object.
//
//	see SetDefaults
//	Parameters:
//		- parameters: Parameters with parameters to override the current values.
//		- recursive bool true to perform deep copy, and false for shallow copy. Default: false
//	Returns: *Parameters a new Parameters object.
func (c *Parameters) Override(parameters *Parameters, recursive bool) *Parameters {
	if parameters == nil {
		return c
	}

	result := NewEmptyParameters()
	if recursive {
		reflect.RecursiveObjectWriter.CopyProperties(result.InnerValue(), c.InnerValue())
		reflect.RecursiveObjectWriter.CopyProperties(result.InnerValue(), parameters.InnerValue())
	} else {
		reflect.ObjectWriter.SetProperties(result.InnerValue(), c.InnerValue().(map[string]any))
		reflect.ObjectWriter.SetProperties(result.InnerValue(), parameters.InnerValue().(map[string]any))
	}
	return result
}

// SetDefaults set default values from specified Parameters and returns a new Parameters object.
//
//	see Override
//	Parameters:
//		- defaultParameters *Parameters Parameters with default parameter values.
//		- recursive bool true to perform deep copy, and false for shallow copy. Default: false
//	Returns: *Parameters a new Parameters object.
func (c *Parameters) SetDefaults(defaultParameters *Parameters, recursive bool) *Parameters {
	if defaultParameters == nil {
		return c
	}

	result := NewEmptyParameters()
	if recursive {
		reflect.RecursiveObjectWriter.CopyProperties(result.InnerValue(), defaultParameters.InnerValue())
		reflect.RecursiveObjectWriter.CopyProperties(result.InnerValue(), c.InnerValue())
	} else {
		reflect.ObjectWriter.SetProperties(result.InnerValue(), defaultParameters.InnerValue().(map[string]any))
		reflect.ObjectWriter.SetProperties(result.InnerValue(), c.InnerValue().(map[string]any))
	}
	return result
}

// AssignTo assigns (copies over) properties from the specified value to this map.
//
//	Parameters: value any value whose properties shall be copied over.
func (c *Parameters) AssignTo(value any) {
	if value == nil {
		return
	}
	reflect.RecursiveObjectWriter.CopyProperties(value, c.InnerValue())
}

// Pick picks select parameters from this Parameters and returns them as a new Parameters object.
//
//	Parameters: paths ...string keys to be picked and copied over to new Parameters.
//	Returns: *Parameters a new Parameters object.
func (c *Parameters) Pick(paths ...string) *Parameters {
	result := NewEmptyParameters()
	for _, path := range paths {
		if c.Contains(path) {
			if val, ok := c.Get(path); ok {
				result.Put(path, val)
			}
		}
	}
	return result
}

// Omit selected parameters from this Parameters and returns the rest as a new Parameters object.
//
//	Parameters: paths ...string keys to be omitted from copying over to new Parameters.
//	Returns: *Parameters a new Parameters object.
func (c *Parameters) Omit(paths ...string) *Parameters {
	result := NewParametersFromValue(c.InnerValue())
	for _, path := range paths {
		result.Remove(path)
	}
	return result
}
