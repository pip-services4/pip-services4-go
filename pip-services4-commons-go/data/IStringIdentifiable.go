package data

// IStringIdentifiable interface for data objects to operate with ids.
//	Example
//		type MyStruct struct {
//			...
//			id string
//		}
//
//		func (c *MyStruct) GetId() string {
//			return c.id
//		}
//		func (c *MyStruct) SetId(id string) {
//			c.id = id
//		}
type IStringIdentifiable interface {
	GetId() string
	SetId(string)
}
