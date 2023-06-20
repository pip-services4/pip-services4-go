package parsers

// Define types of mustache tokens
const (
	TokenUnknown = iota
	TokenValue
	TokenVariable
	TokenEscapedVariable
	TokenSection
	TokenInvertedSection
	TokenSectionEnd
	TokenPartial
	TokenComment
)
