// Container configuration serves as a recipe for instantiating and configuring components inside the container.
// External configurations (stored as YAML or JSON) are passed to the container and define
// the structure of objects that need to be recreated in the container.
//	Objects can be defined in two ways:
//		using descriptors (using which registered factories can recreate the object)
//		using hard-coded types (objects are recreated directly, based on their type, bypassing factories).
// In addition, various configurations are stored for each object. The container recreates the objects and, if they implement the IConfigurable interface, passes them their configurations.

package config
