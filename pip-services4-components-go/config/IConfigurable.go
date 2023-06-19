package config

import "context"

// IConfigurable An interface to set configuration parameters to an object.
//
// It can be added to any existing class by implementing a single configure() method.
//
// If you need to emphasis the fact that Configure() method can be called multiple
// times to change object configuration in runtime, use IReconfigurable interface instead.
//
// Example:
//		type MyStruct struct  {
//		     myParam string
//		}
//
//		func NewMyStruct() *MyStruct {
//		    return &MyStruct{
//		        myParam: "default value",
//		    },
//		}
//
//		// Implement configure
//		func (c* MyStruct) Configure(ctx context.Context, config *cconf.ConfigParams)  {
//		    c.myParam = config.GetAsStringWithDefault("options.param", myParam);
//		    ...
//		}
type IConfigurable interface {

	// Configure object by passing configuration parameters.
	//	Parameters:
	//		- ctx context.Context
	//		- config: ConfigParams configuration parameters to be set.
	Configure(ctx context.Context, config *ConfigParams)
}
