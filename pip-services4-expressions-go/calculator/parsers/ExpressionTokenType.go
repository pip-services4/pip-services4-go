package parsers

// Define types of expression tokens.
const (
	Unknown = iota
	LeftBrace
	RightBrace
	LeftSquareBrace
	RightSquareBrace
	Plus
	Minus
	Star
	Slash
	Procent
	Power
	Equal
	NotEqual
	More
	Less
	EqualMore
	EqualLess
	ShiftLeft
	ShiftRight
	And
	Or
	Xor
	Is
	In
	NotIn
	Element
	Null
	Not
	Like
	NotLike
	IsNull
	IsNotNull
	Comma
	Unary
	Function
	Variable
	Constant
)
