package config

import "strings"

// NameResolver is a helper class that allows to extract component name from configuration parameters.
// The name can be defined in "id", "name" parameters or inside a component descriptor.
//	Example:
//		config := NewConfigParamsFromTuples(
//			"descriptor", "myservice:connector:aws:connector1:1.0",
//			"param1", "ABC",
//			"param2", 123
//		);
//		name := NameResolver.Resolve(config); // Result: connector1
var NameResolver = &_TNameResolver{}

type _TNameResolver struct{}

// Resolve resolves a component name from configuration parameters. The name can be stored in "id",
// "name" fields or inside a component descriptor. If name cannot be determined it returns a empty string.
//	Parameters: config: ConfigParams configuration parameters that may contain a component name.
//	Returns string resolved name or "" if the name cannot be determined.
func (c *_TNameResolver) Resolve(config *ConfigParams) string {
	return c.ResolveWithDefault(config, "")
}

// ResolveWithDefault resolves a component name from configuration parameters. The name can be stored in "id",
// "name" fields or inside a component descriptor. If name cannot be determined it returns a defaultName.
//	Parameters:
//		- config: ConfigParams configuration parameters that may contain a component name.
//		- defaultName: string a default component name.
//	Returns: string resolved name or default name if the name cannot be determined.
func (c *_TNameResolver) ResolveWithDefault(config *ConfigParams, defaultName string) string {
	var name = config.GetAsString("name")

	if name == "" {
		name = config.GetAsString("id")
	}

	if name == "" {
		descriptorStr := config.GetAsString("descriptor")
		if descriptorStr != "" {
			tokens := strings.Split(descriptorStr, ":")
			if len(tokens) == 5 {
				name = tokens[3]
			}
		}
	}

	if name == "" {
		name = defaultName
	}

	return name
}
