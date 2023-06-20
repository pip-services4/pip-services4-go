package parsers

// General syntax errors
const (
	// ErrCodeUnknown the unknown
	ErrCodeUnknown = "UNKNOWN"

	// ErrCodeInternal the internal error
	ErrCodeInternal = "INTERNAL"

	// ErrCodeUnexpectedEnd the unexpected end.
	ErrCodeUnexpectedEnd = "UNEXPECTED_END"

	// ErrCodeErrorNear the error near
	ErrCodeErrorNear = "ERROR_NEAR"

	// ErrCodeErrorAt the error at
	ErrCodeErrorAt = "ERROR_AT"

	// ErrCodeUnexpectedSymbol the unexpected symbol
	ErrCodeUnexpectedSymbol = "UNEXPECTED_SYMBOL"

	// ErrCodeMismatchedBrackets the mismatched brackets
	ErrCodeMismatchedBrackets = "MISTMATCHED_BRACKETS"

	// ErrCodeMissingVariable the missing variable
	ErrCodeMissingVariable = "MISSING_VARIABLE"

	// ErrCodeNotClosedSection not closed section
	ErrCodeNotClosedSection = "NOT_CLOSED_SECTION"

	// ErrCodeUnexpectedSectionEnd unexpected section end
	ErrCodeUnexpectedSectionEnd = "UNEXPECTED_SECTION_END"
)
