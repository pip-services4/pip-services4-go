package data

// Interface for data objects that have human-readable names.
//	Example
//		type MyStruct struct {
//			...
//			name string
//		}
//
//		func (c *MyStruct) GetName() string {
//			return c.name
//		}
//		func (c *MyStruct) SetName(name string) {
//			c.name = name
//		}
type INamed interface {
	GetName() string
}
