package data

// Interface for data objects that can be versioned.
//
// Versioning is often used as optimistic concurrency mechanism.
//
// The version doesn't have to be a number, but it is recommended to use sequential
// values to determine if one object has newer or older version than another one.
//
// It is a common pattern to use the time of change as the object version.
//	Example
//		type MyStruct struct {
//			...
//			version string
//		}
//
//		func (c *MyStruct) GetVersion() string {
//			return c.version
//		}
//		func (c *MyStruct) SetVersion(version string) {
//			c.version = version
//		}
//		func (c *MyStruct) UpdateData(ctx context.Context, item: MyData) {
//			if (item.version < oldItem.version) {
//				panic("VERSION_CONFLICT")
//			}
//		}
//
type IVersioned interface {
	// The object's version.
	GetVersion() string
}
