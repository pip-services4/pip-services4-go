package parsers

// Define states in mustache lexical analysis
const (
	StateValue = iota
	StateOperator1
	StateOperator2
	StateVariable
	StateComment
	StateClosure
)
