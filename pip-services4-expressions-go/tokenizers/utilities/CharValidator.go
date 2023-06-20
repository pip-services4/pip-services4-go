package utilities

// TCharValidator validates characters that are processed by Tokenizers.
type _TCharValidator struct {
}

var CharValidator = &_TCharValidator{}

const Eof rune = -1

func (c *_TCharValidator) IsEof(value rune) bool {
	return value == -1
}

func (c *_TCharValidator) IsEol(value rune) bool {
	return value == 10 || value == 13
}

func (c *_TCharValidator) IsDigit(value rune) bool {
	return value >= '0' && value <= '9'
}
