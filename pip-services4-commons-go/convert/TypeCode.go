package convert

// TypeCode codes for the data types that can be converted using TypeConverter.
type TypeCode int

const (
	Unknown  TypeCode = iota
	String   TypeCode = iota
	Boolean  TypeCode = iota
	Integer  TypeCode = iota
	Long     TypeCode = iota
	Float    TypeCode = iota
	Double   TypeCode = iota
	DateTime TypeCode = iota
	Duration TypeCode = iota
	Object   TypeCode = iota
	Enum     TypeCode = iota
	Array    TypeCode = iota
	Map      TypeCode = iota
)
