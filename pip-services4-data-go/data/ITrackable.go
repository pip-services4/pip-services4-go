package data

import "time"

// Interface for data objects that can track their changes, including logical deletion.
//
//	Example
//		type MyStruct struct {
//			...
//			changeTime time.Time
//			createTime time.Time
//			deleted bool
//		}
//
//		func (c *MyStruct) GetChangeTime() string {
//			return c.changeTime
//		}
//		func (c *MyStruct) SetDeleted(deleted bool) {
//			c.deleted = deleted
//		}
//
type ITrackable interface {
	// The UTC time at which the object was created.
	GetCreateTime() time.Time
	// The UTC time at which the object was last changed (created, updated, or deleted).
	GetChangeTime() time.Time
	// The logical deletion flag. True when object is deleted and null or false otherwise
	GetDeleted() bool
}
