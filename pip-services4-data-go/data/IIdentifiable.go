package data

// IIdentifiable interface for data objects to operate with ids.
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
type IIdentifiable[K any] interface {
	GetId() K
}
