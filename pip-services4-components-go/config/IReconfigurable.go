package config

// IReconfigurable an interface to set configuration parameters to an object.
//
// It is similar to IConfigurable interface, but emphasises the fact
// that Configure() method can be called more than once to change object configuration in runtime.
type IReconfigurable interface {
	IConfigurable
}
