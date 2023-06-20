package data

import "time"

// Interface for data objects that contain their latest change time.
//	Example
//		type MyStruct struct {
//			...
//			changeTime time.Time
//		}
//
//		func (c *MyStruct) GetChangeTime() time.Time {
//			return c.changeTime
//		}
//		func (c *MyStruct) SetGetChangeTime(changeTime time.Time) {
//			c.changeTime = changeTime
//		}
type IChangeable interface {
	// The UTC time at which the object was last changed (created or updated).
	GetChangeTime() time.Time
}
