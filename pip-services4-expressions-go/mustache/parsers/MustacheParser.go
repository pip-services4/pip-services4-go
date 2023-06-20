package parsers

import (
	"strings"

	merr "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache/errors"
	mtokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache/tokenizers"
	tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
)

// MustacheParser implements an mustache parser class.
type MustacheParser struct {
	tokenizer         tokenizers.ITokenizer
	template          string
	originalTokens    []*tokenizers.Token
	initialTokens     []*MustacheToken
	currentTokenIndex int
	variableNames     []string
	resultTokens      []*MustacheToken
}

func NewMustacheParser() *MustacheParser {
	c := &MustacheParser{
		tokenizer: mtokenizers.NewMustacheTokenizer(),
	}
	c.Clear()
	return c
}

// Template gets the mustache template.
func (c *MustacheParser) Template() string {
	return c.template
}

// SetTemplate sets the mustache template.
func (c *MustacheParser) SetTemplate(value string) error {
	return c.ParseString(value)
}

// OriginalTokens gets the list of original tokens
func (c *MustacheParser) OriginalTokens() []*tokenizers.Token {
	return c.originalTokens
}

// SetOriginalTokens sets the list of original tokens
func (c *MustacheParser) SetOriginalTokens(value []*tokenizers.Token) error {
	return c.ParseTokens(value)
}

// InitialTokens gets the list of original mustache tokens.
func (c *MustacheParser) InitialTokens() []*MustacheToken {
	return c.initialTokens
}

// ResultTokens sets the list of parsed mustache tokens.
func (c *MustacheParser) ResultTokens() []*MustacheToken {
	return c.resultTokens
}

// VariableNames gets the list of found variable names.
func (c *MustacheParser) VariableNames() []string {
	return c.variableNames
}

// ParseString sets a new mustache string and parses it into internal byte code.
//
//	Parameters:
//		mustache: A new mustache string.
func (c *MustacheParser) ParseString(mustache string) error {
	c.Clear()
	c.template = strings.Trim(mustache, " \t\r\n")
	c.originalTokens = c.tokenizeMustache(c.template)
	return c.performParsing()
}

func (c *MustacheParser) ParseTokens(tokens []*tokenizers.Token) error {
	c.Clear()
	c.originalTokens = tokens
	c.template = c.composeMustache(tokens)
	return c.performParsing()
}

// Clear parsing results.
func (c *MustacheParser) Clear() {
	c.template = ""
	c.originalTokens = []*tokenizers.Token{}
	c.initialTokens = []*MustacheToken{}
	c.resultTokens = []*MustacheToken{}
	c.currentTokenIndex = 0
	c.variableNames = []string{}
}

// hasMoreTokens checks are there more tokens for processing.
//
//	Returns: <code>true</code> if some tokens are present.
func (c *MustacheParser) hasMoreTokens() bool {
	return c.currentTokenIndex < len(c.initialTokens)
}

// checkForMoreTokens checks are there more tokens available and throws exception if no more tokens available.
func (c *MustacheParser) checkForMoreTokens() error {
	if !c.hasMoreTokens() {
		err := merr.NewMustacheError("", ErrCodeUnexpectedEnd, "Unexpected end of mustache", 0, 0)
		return err
	}
	return nil
}

// getCurrentToken gets the current token object.
//
//	Returns: The current token object.
func (c *MustacheParser) getCurrentToken() *MustacheToken {
	if c.currentTokenIndex < len(c.initialTokens) {
		return c.initialTokens[c.currentTokenIndex]
	}
	return nil
}

// getNextToken gets the next token object.
//
//	Returns: The next token object.
func (c *MustacheParser) getNextToken() *MustacheToken {
	if (c.currentTokenIndex + 1) < len(c.initialTokens) {
		return c.initialTokens[c.currentTokenIndex+1]
	}
	return nil
}

// moveToNextToken moves to the next token object.
func (c *MustacheParser) moveToNextToken() {
	c.currentTokenIndex++
}

// addTokenToResult adds a mustache to the result list
//
//	Parameters:
//		- type: The type of the token to be added.
//		- value: The value of the token to be added.
//		- line: The line where the token is.
//		- column: The column number where the token is.
func (c *MustacheParser) addTokenToResult(typ int, value string, line int, column int) *MustacheToken {
	token := NewMustacheToken(typ, value, line, column)
	c.resultTokens = append(c.resultTokens, token)
	return token
}

func (c *MustacheParser) tokenizeMustache(mustache string) []*tokenizers.Token {
	mustache = strings.Trim(mustache, " \t\r\n")
	if len(mustache) == 0 {
		return []*tokenizers.Token{}
	}

	c.tokenizer.SetSkipWhitespaces(true)
	c.tokenizer.SetSkipComments(true)
	c.tokenizer.SetSkipEof(true)
	c.tokenizer.SetDecodeStrings(true)
	return c.tokenizer.TokenizeBuffer(mustache)
}

func (c *MustacheParser) composeMustache(tokens []*tokenizers.Token) string {
	builder := strings.Builder{}
	for _, token := range tokens {
		builder.WriteString(token.Value())
	}
	return builder.String()
}

func (c *MustacheParser) performParsing() error {
	if len(c.originalTokens) > 0 {
		err := c.completeLexicalAnalysis()
		if err != nil {
			return err
		}

		err = c.performSyntaxAnalysis()
		if err != nil {
			return err
		}

		if c.hasMoreTokens() {
			token := c.getCurrentToken()
			err = merr.NewMustacheError(
				"",
				ErrCodeErrorNear,
				"Syntax error near "+token.Value(),
				token.Line(),
				token.Column(),
			)
			return err
		}

		c.lookupVariables()
	}
	return nil
}

// completeLexicalAnalysis tokenizes the given mustache and prepares an initial tokens list.
func (c *MustacheParser) completeLexicalAnalysis() error {
	state := StateValue
	closingBracket := ""
	operator1 := ""
	operator2 := ""
	variable := ""

	for _, token := range c.originalTokens {
		tokenType := TokenUnknown
		tokenValue := ""

		if state == StateComment {
			if token.Value() == "}}" || token.Value() == "}}}" {
				state = StateClosure
			} else {
				continue
			}
		}

		switch token.Type() {
		case tokenizers.Special:
			if state == StateValue {
				tokenType = TokenValue
				tokenValue = token.Value()
			}
			break
		case tokenizers.Symbol:
			if state == StateValue && (token.Value() == "{{" || token.Value() == "{{{") {
				closingBracket = "}}"
				if token.Value() == "{{{" {
					closingBracket = "}}}"
				}
				state = StateOperator1
				continue
			}
			if state == StateOperator1 && token.Value() == "!" {
				operator1 = token.Value()
				state = StateComment
				continue
			}
			if state == StateOperator1 && (token.Value() == "/" || token.Value() == "#" || token.Value() == "^") {
				operator1 = token.Value()
				state = StateOperator2
				continue
			}

			if state == StateVariable && (token.Value() == "}}" || token.Value() == "}}}") {
				if operator1 != "/" {
					variable = operator2
					operator2 = ""
				}
				state = StateClosure
				// Pass through
			}
			if state == StateClosure && (token.Value() == "}}" || token.Value() == "}}}") {
				if closingBracket != token.Value() {
					err := merr.NewMustacheError("", ErrCodeMismatchedBrackets, "Mismatched brackets. Expected '"+closingBracket+"'", token.Line(), token.Column())
					return err
				}

				if operator1 == "#" && (operator2 == "" || operator2 == "if") {
					tokenType = TokenSection
					tokenValue = variable
				}

				if operator1 == "#" && operator2 == "unless" {
					tokenType = TokenInvertedSection
					tokenValue = variable
				}

				if operator1 == "^" && operator2 == "" {
					tokenType = TokenInvertedSection
					tokenValue = variable
				}

				if operator1 == "/" {
					tokenType = TokenSectionEnd
					tokenValue = variable
				}

				if operator1 == "" {
					tokenType = TokenVariable
					if closingBracket == "}}}" {
						tokenType = TokenEscapedVariable
					}
					tokenValue = variable
				}

				if tokenType == TokenUnknown {
					err := merr.NewMustacheError("", ErrCodeInternal, "Internal error", token.Line(), token.Column())
					return err
				}

				operator1 = ""
				operator2 = ""
				variable = ""
				state = StateValue
			}
			break
		case tokenizers.Word:
			if state == StateOperator1 {
				state = StateVariable
			}
			if state == StateOperator2 && (token.Value() == "if" || token.Value() == "unless") {
				operator2 = token.Value()
				state = StateVariable
				continue
			}
			if state == StateOperator2 {
				state = StateVariable
			}
			if state == StateVariable {
				variable = token.Value()
				state = StateClosure
				continue
			}
			break
		case tokenizers.Whitespace:
			continue
		}
		if tokenType == TokenUnknown {
			err := merr.NewMustacheError("", ErrCodeUnexpectedSymbol, "Unexpected symbol '"+token.Value()+"'", token.Line(), token.Column())
			return err
		}
		c.initialTokens = append(c.initialTokens, NewMustacheToken(tokenType, tokenValue, token.Line(), token.Column()))
	}

	if state != StateValue {
		err := merr.NewMustacheError("", ErrCodeUnexpectedEnd, "Unexpected end of file", 0, 0)
		return err
	}

	return nil
}

// Performs a syntax analysis at level 0.
func (c *MustacheParser) performSyntaxAnalysis() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		c.moveToNextToken()

		if token.Type() == TokenSectionEnd {
			err = merr.NewMustacheError("", ErrCodeUnexpectedSectionEnd, "Unexpected section end for variable '"+token.Value()+"'", token.Line(), token.Column())
			return err
		}

		result := c.addTokenToResult(token.Type(), token.Value(), token.Line(), token.Column())

		if token.Type() == TokenSection || token.Type() == TokenInvertedSection {
			sectionTokens, err := c.performSyntaxAnalysisForSection(token.Value())
			if err != nil {
				return err
			}
			result.SetTokens(append(result.Tokens(), sectionTokens...))
		}
	}

	return nil
}

// Performs a syntax analysis for section
func (c *MustacheParser) performSyntaxAnalysisForSection(variable string) ([]*MustacheToken, error) {
	result := []*MustacheToken{}

	err := c.checkForMoreTokens()
	if err != nil {
		return nil, err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		c.moveToNextToken()

		if token.Type() == TokenSectionEnd && (token.Value() == variable || token.Value() == "") {
			return result, nil
		}

		if token.Type() == TokenSectionEnd {
			err := merr.NewMustacheError("", ErrCodeUnexpectedSectionEnd, "Unexpected section end for variable '"+variable+"'", token.Line(), token.Column())
			return nil, err
		}

		resultToken := NewMustacheToken(token.Type(), token.Value(), token.Line(), token.Column())

		if token.Type() == TokenSection || token.Type() == TokenInvertedSection {
			sectionTokens, err := c.performSyntaxAnalysisForSection(token.Value())
			if err != nil {
				return nil, err
			}
			resultToken.SetTokens(append(resultToken.Tokens(), sectionTokens...))
		}

		result = append(result, resultToken)
	}

	token := c.getCurrentToken()
	err = merr.NewMustacheError("", ErrCodeNotClosedSection, "Not closed section for variable '"+variable+"'", token.Line(), token.Column())
	return nil, err
}

// lookupVariables retrieves variables from the parsed output.
func (c *MustacheParser) lookupVariables() {
	if len(c.originalTokens) == 0 {
		return
	}

	c.variableNames = []string{}
	for _, token := range c.initialTokens {
		if token.Type() != TokenValue && token.Type() != TokenComment && token.Value() != "" {
			variableName := strings.ToLower(token.Value())
			found := false
			for _, v := range c.variableNames {
				if strings.ToLower(v) == variableName {
					found = true
				}
			}
			if !found {
				c.variableNames = append(c.variableNames, token.Value())
			}
		}
	}
}
