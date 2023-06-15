/*
Contains "soft" data converters. Soft data converters differ from the data conversion
algorithms found in typical programming language, due to the fact that they support rare
conversions between various data types (such as integer to timespan, timespan to string, and so on).

These converters are necessary, due to the fact that data in enterprise systems is represented in
various forms and conversion is often necessary – at times in very difficult combinations.
*/

package convert
