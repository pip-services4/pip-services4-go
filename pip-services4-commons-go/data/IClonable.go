package data

// ICloneable Interface for data objects that are able to create their full binary copy.
//	Example
//		type MyStruct struct {
//			...
//		}
//
//		func (c *MyStruct) Clone() *MyStruct {
//			cloneObj := new(MyStruct)
//			// Copy every attribute from this to cloneObj here.
//			...
//			return cloneObj
//		}
type ICloneable[T any] interface {
	// Clone creates a binary clone of this object.
	// returns: a clone of this object.
	Clone() T
}
