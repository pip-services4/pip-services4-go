package config

// OptionsResolver is a helper class to parameters from "options" configuration section.
//
//	Example:
//		config := NewConfigParamsFromTuples(
//			...
//			"options.param1", "ABC",
//			"options.param2", 123
//		);
//
//	options := OptionsResolver.resolve(config); // Result: param1=ABC;param2=123
var OptionsResolver = &_TOptionsResolver{}

type _TOptionsResolver struct{}

const SectionOptions = "options"

// Resolve configuration section from component configuration parameters.
//	Parameters: config: ConfigParams configuration parameters
//	Returns: *ConfigParams configuration parameters from "options" section
func (c *_TOptionsResolver) Resolve(config *ConfigParams) *ConfigParams {
	var options = config.GetSection(SectionOptions)
	return options
}

// ResolveWithDefault an "options" configuration section from component configuration parameters.
//	Parameters:
//		- config: ConfigParams configuration parameters
//		- configAsDefault: boolean
//			When set true the method returns the entire parameter
//			set when "options" section is not found. Default: false
//	Returns: *ConfigParams configuration parameters from "options" section
func (c *_TOptionsResolver) ResolveWithDefault(config *ConfigParams) *ConfigParams {
	var options = c.Resolve(config)

	if options.Len() == 0 {
		options = config
	}

	return options
}
