package config

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

// ConfigParams Contains a key-value map with configuration parameters.
// All values stored as strings and can be serialized
// as JSON or string forms. When retrieved the values can be automatically converted on read using GetAsXXX methods.
// The keys are case-sensitive, so it is recommended to use consistent C-style as: "my_param"
//
// Configuration parameters can be broken into sections and subsections using dot notation as:
// "section1.subsection1.param1". Using GetSection method all parameters
// from specified section can be extracted from a ConfigMap.
//
// The ConfigParams supports serialization from/to plain strings as:
// "key1=123;key2=ABC;key3=2016-09-16T00:00:00.00Z"
//
// ConfigParams are used to pass configurations to IConfigurable objects.
// They also serve as a basis for more concrete configurations such as ConnectionParams or
// CredentialParams (in the Pip.Services components package).
//
//	see IConfigurable
//	see StringValueMap
//
//	Example:
//		config := NewConfigParamsFromTuples(
//			"section1.key1", "AAA",
//			"section1.key2", 123,
//			"section2.key1", true
//		);
//
//		config.GetAsString("section1.key1"); // Result: AAA
//		config.GetAsInteger("section1.key1"); // Result: 0
//
//		section1 = config.GetSection("section2");
//		section1.GetAsString("key1"); // Result: true
type ConfigParams struct {
	*data.StringValueMap
}

// NewEmptyConfigParams creates a new empty ConfigParams object.
// Returns: *ConfigParams a new empty ConfigParams object.
func NewEmptyConfigParams() *ConfigParams {
	return &ConfigParams{
		StringValueMap: data.NewEmptyStringValueMap(),
	}
}

// NewConfigParams creates a new ConfigParams from map.
//
//	Parameters: values ...map[string]string
//	Returns: *ConfigParams a newly created ConfigParams.
func NewConfigParams(values map[string]string) *ConfigParams {
	return &ConfigParams{
		StringValueMap: data.NewStringValueMap(values),
	}
}

// NewConfigParamsFromValue creates a new ConfigParams object filled with key-value pairs from specified object.
//
//	Parameters: value any an object with key-value pairs used to initialize a new ConfigParams.
//	Returns: *ConfigParams a new ConfigParams object.
func NewConfigParamsFromValue(value any) *ConfigParams {
	values := reflect.RecursiveObjectReader.GetProperties(value)
	return &ConfigParams{
		StringValueMap: data.NewStringValueMapFromValue(values),
	}
}

// NewConfigParamsFromTuples creates a new ConfigParams object filled with provided key-value pairs called tuples.
// Tuples parameters contain a sequence of key1, value1, key2, value2, ... pairs.
//
//	see convert.StringValueMap.fromTuplesArray
//	Parameters: tuples ...any the tuples to fill a new ConfigParams object.
//	Returns ConfigParams a new ConfigParams object.
func NewConfigParamsFromTuples(tuples ...any) *ConfigParams {
	return &ConfigParams{
		StringValueMap: data.NewStringValueMapFromTuplesArray(tuples),
	}
}

// NewConfigParamsFromTuplesArray creates a new StringValueMap from a list of key-value pairs called tuples.
// The method is similar to fromTuples but tuples are passed as array instead of parameters.
//
//	Parameters: tuples []any a list of values where odd elements
//		are keys and the following even elements are values
//	Returns *ConfigParams a newly created ConfigParams.
func NewConfigParamsFromTuplesArray(tuples []any) *ConfigParams {
	return &ConfigParams{
		StringValueMap: data.NewStringValueMapFromTuplesArray(tuples),
	}
}

// NewConfigParamsFromString creates a new ConfigParams object filled with key-value pairs serialized as a string.
//
//	see convert.StringValueMap.fromString
//	Parameters: line: string a string with serialized key-value pairs as "key1=value1;key2=value2;..."
//	Example: "Key1=123;Key2=ABC;Key3=2016-09-16T00:00:00.00Z"
//	Returns: *ConfigParams a new ConfigParams object.
func NewConfigParamsFromString(line string) *ConfigParams {
	return &ConfigParams{
		StringValueMap: data.NewStringValueMapFromString(line),
	}
}

// NewConfigParamsFromMaps creates a new ConfigParams by merging two or more maps.
// Maps defined later in the list override values from previously defined maps.
//
//	Parameters: maps ...map[string]string an array of maps to be merged
//	Returns: *ConfigParams a newly created ConfigParams.
func NewConfigParamsFromMaps(maps ...map[string]string) *ConfigParams {
	return &ConfigParams{
		StringValueMap: data.NewStringValueMapFromMaps(maps...),
	}
}

// GetSectionNames gets a list with all 1st level section names.
//
//	Returns: []string a list of section names stored in this ConfigMap.
func (c *ConfigParams) GetSectionNames() []string {
	sections := make([]string, 0)

	for key := range c.Value() {
		pos := strings.Index(key, ".")
		section := key
		if pos > 0 {
			section = key[0:pos]
		}

		// Perform case sensitive search
		found := false
		for index := 0; index < len(sections); index++ {
			if section == sections[index] {
				found = true
				break
			}
		}

		if !found {
			sections = append(sections, section)
		}
	}

	return sections
}

// GetSection gets parameters from specific section stored in this ConfigMap. The section name is removed from parameter keys.
//
//	Parameters: section: string name of the section to retrieve configuration parameters from.
//	Returns: *ConfigParams all configuration parameters that belong to the section named 'section'.
func (c *ConfigParams) GetSection(section string) *ConfigParams {
	result := NewEmptyConfigParams()
	prefix := section + "."

	for key := range c.Value() {
		// Prevents exception on the next line
		if len(key) < len(prefix) {
			continue
		}

		// Perform case sensitive match
		keyPrefix := key[0:len(prefix)]
		if keyPrefix == prefix {
			sectionKey := key[len(prefix):]
			val, _ := c.Get(key)
			result.Put(sectionKey, val)
		}
	}

	return result
}

// AddSection adds parameters into this ConfigParams under specified section.
// Keys for the new parameters are appended with section dot prefix.
//
//	Parameters:
//		- section: string name of the section where add new parameters
//		- sectionParams: *ConfigParams new parameters to be added.
func (c *ConfigParams) AddSection(section string, sectionParams *ConfigParams) {
	if section == "" {
		panic("Section name cannot be empty")
	}

	if sectionParams != nil {
		for key := range sectionParams.Value() {
			sectionKey := key

			if len(sectionKey) > 0 {
				sectionKey = section + "." + sectionKey
			} else {
				sectionKey = section
			}

			var value any
			if val, ok := (*sectionParams).Get(key); ok {
				value = val
			}

			c.Put(sectionKey, value)
		}
	}
}

// Override overrides parameters with new values from specified ConfigParams and returns a new ConfigParams object.
//
//	see NewConfigParamsFromMaps
//	Parameters: configParams: *ConfigParams ConfigMap with parameters to override the current values.
//	Returns *ConfigParams a new ConfigParams object.
func (c *ConfigParams) Override(configParams *ConfigParams) *ConfigParams {
	return NewConfigParamsFromMaps(c.Value(), configParams.Value())
}

// SetDefaults set default values from specified ConfigParams and returns a new ConfigParams object.
//
//	see NewConfigParamsFromMaps
//	Parameters: defaultConfigParams: *ConfigParams ConfigMap with default parameter values.
//	Returns *ConfigParams a new ConfigParams object.
func (c *ConfigParams) SetDefaults(defaults *ConfigParams) *ConfigParams {
	return NewConfigParamsFromMaps(defaults.Value(), c.Value())
}
