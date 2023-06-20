package parsers

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/errors"
	ctokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
)

// Implements an expression parser class.
type ExpressionParser struct {
	tokenizer         tokenizers.ITokenizer
	expression        string
	originalTokens    []*tokenizers.Token
	initialTokens     []*ExpressionToken
	currentTokenIndex int
	variableNames     []string
	resultTokens      []*ExpressionToken
}

// Defines a list of operators.
var operators []string = []string{
	"(", ")", "[", "]", "+", "-", "*", "/", "%", "^",
	"=", "<>", "!=", ">", "<", ">=", "<=", "<<", ">>",
	"AND", "OR", "XOR", "NOT", "IS", "IN", "NULL", "LIKE", ",",
}

// Defines a list of operator token types.
var operatorTypes []int = []int{
	LeftBrace, RightBrace, LeftSquareBrace, RightSquareBrace,
	Plus, Minus, Star, Slash, Procent, Power, Equal, NotEqual,
	NotEqual, More, Less, EqualMore, EqualLess, ShiftLeft,
	ShiftRight, And, Or, Xor, Not, Is, In, Null, Like, Comma,
}

func NewExpressionParser() *ExpressionParser {
	c := &ExpressionParser{
		tokenizer:      ctokenizers.NewExpressionTokenizer(),
		originalTokens: []*tokenizers.Token{},
		initialTokens:  []*ExpressionToken{},
		variableNames:  []string{},
		resultTokens:   []*ExpressionToken{},
	}
	return c
}

// Gets the expression string.
func (c *ExpressionParser) Expression() string {
	return c.expression
}

// Sets the expression string.
func (c *ExpressionParser) SetExpression(value string) error {
	return c.ParseString(value)
}

// Gets the original tokens
func (c *ExpressionParser) OriginalTokens() []*tokenizers.Token {
	return c.originalTokens
}

// Sets the original tokens
func (c *ExpressionParser) SetOriginalTokens(value []*tokenizers.Token) error {
	return c.ParseTokens(value)
}

// Gets the list of original expression tokens.
func (c *ExpressionParser) InitialTokens() []*ExpressionToken {
	return c.initialTokens
}

// Gets the list of parsed expression tokens.
func (c *ExpressionParser) ResultTokens() []*ExpressionToken {
	return c.resultTokens
}

// Gets the list of found variable names.
func (c *ExpressionParser) VariableNames() []string {
	return c.variableNames
}

// Sets a new expression string and parses it into internal byte code.
//
// Parameters:
//   - expression: A new expression string.
func (c *ExpressionParser) ParseString(expression string) error {
	c.Clear()
	c.expression = strings.Trim(expression, " \t\r\n")
	c.originalTokens = c.tokenizeExpression(c.expression)
	return c.performParsing()
}

func (c *ExpressionParser) ParseTokens(tokens []*tokenizers.Token) error {
	c.Clear()
	c.originalTokens = tokens
	c.expression = c.composeExpression(tokens)
	return c.performParsing()
}

// Clears parsing results.
func (c *ExpressionParser) Clear() {
	c.expression = ""
	c.originalTokens = []*tokenizers.Token{}
	c.initialTokens = []*ExpressionToken{}
	c.resultTokens = []*ExpressionToken{}
	c.currentTokenIndex = 0
	c.variableNames = []string{}
}

// Checks are there more tokens for processing.
//
// Returns: <code>true</code> if some tokens are present.
func (c *ExpressionParser) hasMoreTokens() bool {
	return c.currentTokenIndex < len(c.initialTokens)
}

// Checks are there more tokens available and throws exception if no more tokens available.
func (c *ExpressionParser) checkForMoreTokens() error {
	if !c.hasMoreTokens() {
		return errors.NewSyntaxError("", errors.ErrUnexpectedEnd, "Unexpected end of expression.", 0, 0)
	}
	return nil
}

// Gets the current token object.
//
// Returns: The current token object.
func (c *ExpressionParser) getCurrentToken() *ExpressionToken {
	if c.currentTokenIndex < len(c.initialTokens) {
		return c.initialTokens[c.currentTokenIndex]
	}
	return nil
}

// Gets the next token object.
//
// Returns: The next token object.
func (c *ExpressionParser) getNextToken() *ExpressionToken {
	if (c.currentTokenIndex + 1) < len(c.initialTokens) {
		return c.initialTokens[c.currentTokenIndex+1]
	}
	return nil
}

// Moves to the next token object.
func (c *ExpressionParser) moveToNextToken() {
	c.currentTokenIndex = c.currentTokenIndex + 1
}

// Adds an expression to the result list
//
// Parameters:
//   - type: The type of the token to be added.
//   - value: The value of the token to be added.
//   - line: The line number where the token is.
//   - column: The column number where the token is.
func (c *ExpressionParser) addTokenToResult(typ int, value *variants.Variant, line int, column int) {
	c.resultTokens = append(c.resultTokens, NewExpressionToken(typ, value, line, column))
}

// Matches available tokens types with types from the list.
// If tokens matchs then shift the list.
//
// Parameters:
//   - types: A list of token types to compare.
//
// Returns: <code>true</code> if token types match.
func (c *ExpressionParser) matchTokensWithTypes(types ...int) bool {
	matches := false

	for i, typ := range types {
		if c.currentTokenIndex+i < len(c.initialTokens) {
			matches = c.initialTokens[c.currentTokenIndex+i].Type() == typ
		} else {
			matches = false
			break
		}
	}

	if matches {
		c.currentTokenIndex = c.currentTokenIndex + len(types)
	}

	return matches
}

func (c *ExpressionParser) tokenizeExpression(expression string) []*tokenizers.Token {
	expression = strings.Trim(expression, " \t\r\n")
	if len(expression) > 0 {
		c.tokenizer.SetSkipWhitespaces(true)
		c.tokenizer.SetSkipComments(true)
		c.tokenizer.SetSkipEof(true)
		c.tokenizer.SetDecodeStrings(true)
		return c.tokenizer.TokenizeBuffer(expression)
	}
	return []*tokenizers.Token{}
}

func (c *ExpressionParser) composeExpression(tokens []*tokenizers.Token) string {
	builder := strings.Builder{}
	for _, token := range tokens {
		builder.WriteString(token.Value())
	}
	return builder.String()
}

func (c *ExpressionParser) performParsing() error {
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
			err = errors.NewSyntaxError("", errors.ErrErrorNear, "Syntax error near "+token.Value().AsString(), token.Line(), token.Column())
			return err
		}
	}

	return nil
}

// Tokenizes the given expression and prepares an initial tokens list.
func (c *ExpressionParser) completeLexicalAnalysis() error {
	for _, token := range c.originalTokens {
		tokenType := Unknown
		tokenValue := variants.Empty

		switch token.Type() {
		case tokenizers.Comment:
		case tokenizers.Whitespace:
			continue
		case tokenizers.Keyword:
			{
				temp := strings.ToUpper(token.Value())
				if temp == "TRUE" {
					tokenType = Constant
					tokenValue = variants.VariantFromBoolean(true)
				} else if temp == "FALSE" {
					tokenType = Constant
					tokenValue = variants.VariantFromBoolean(false)
				} else {
					for index := 0; index < len(operators); index++ {
						if temp == operators[index] {
							tokenType = operatorTypes[index]
							break
						}
					}
				}
				break
			}
		case tokenizers.Word:
			{
				tokenType = Variable
				tokenValue = variants.VariantFromString(token.Value())
				break
			}
		case tokenizers.Integer:
			{
				tokenType = Constant
				tokenValue = variants.VariantFromInteger(convert.IntegerConverter.ToInteger(token.Value()))
				break
			}
		case tokenizers.Float:
			{
				tokenType = Constant
				tokenValue = variants.VariantFromFloat(convert.FloatConverter.ToFloat(token.Value()))
				break
			}
		case tokenizers.Quoted:
			{
				tokenType = Constant
				tokenValue = variants.VariantFromString(token.Value())
				break
			}
		case tokenizers.Symbol:
			{
				temp := strings.ToUpper(token.Value())
				for i := 0; i < len(operators); i++ {
					if temp == operators[i] {
						tokenType = operatorTypes[i]
						break
					}
				}
				break
			}
		}

		if tokenType == Unknown {
			err := errors.NewSyntaxError("", errors.ErrUnknownSymbol, "Unknown symbol "+token.Value(), token.Line(), token.Column())
			return err
		}

		c.initialTokens = append(c.initialTokens, NewExpressionToken(tokenType, tokenValue, token.Line(), token.Column()))
	}

	return nil
}

// Performs a syntax analysis at level 0.
func (c *ExpressionParser) performSyntaxAnalysis() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	err = c.performSyntaxAnalysisAtLevel1()
	if err != nil {
		return err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		if token.Type() == And || token.Type() == Or || token.Type() == Xor {
			c.moveToNextToken()

			err = c.performSyntaxAnalysisAtLevel1()
			if err != nil {
				return err
			}

			c.addTokenToResult(token.Type(), variants.Empty, token.Line(), token.Column())
			continue
		}
		break
	}

	return nil
}

// Performs a syntax analysis at level 1.
func (c *ExpressionParser) performSyntaxAnalysisAtLevel1() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	token := c.getCurrentToken()
	if token.Type() == Not {
		c.moveToNextToken()

		err = c.performSyntaxAnalysisAtLevel2()
		if err != nil {
			return err
		}

		c.addTokenToResult(token.Type(), variants.Empty, token.Line(), token.Column())
	} else {
		err = c.performSyntaxAnalysisAtLevel2()
		if err != nil {
			return err
		}
	}

	return nil
}

// Performs a syntax analysis at level 2.
func (c *ExpressionParser) performSyntaxAnalysisAtLevel2() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	err = c.performSyntaxAnalysisAtLevel3()
	if err != nil {
		return err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		if token.Type() == Equal || token.Type() == NotEqual || token.Type() == More ||
			token.Type() == Less || token.Type() == EqualMore || token.Type() == EqualLess {
			c.moveToNextToken()

			err = c.performSyntaxAnalysisAtLevel3()
			if err != nil {
				return err
			}

			c.addTokenToResult(token.Type(), variants.Empty, token.Line(), token.Column())
			continue
		}
		break
	}

	return nil
}

// Performs a syntax analysis at level 3.
func (c *ExpressionParser) performSyntaxAnalysisAtLevel3() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	err = c.performSyntaxAnalysisAtLevel4()
	if err != nil {
		return err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		if token.Type() == Plus || token.Type() == Minus || token.Type() == Like {
			c.moveToNextToken()

			err = c.performSyntaxAnalysisAtLevel4()
			if err != nil {
				return err
			}

			c.addTokenToResult(token.Type(), variants.Empty, token.Line(), token.Column())
		} else if c.matchTokensWithTypes(Not, Like) {
			err = c.performSyntaxAnalysisAtLevel4()
			if err != nil {
				return err
			}

			c.addTokenToResult(NotLike, variants.Empty, token.Line(), token.Column())
		} else if c.matchTokensWithTypes(Is, Null) {
			c.addTokenToResult(IsNull, variants.Empty, token.Line(), token.Column())
		} else if c.matchTokensWithTypes(Is, Not, Null) {
			c.addTokenToResult(IsNotNull, variants.Empty, token.Line(), token.Column())
		} else if c.matchTokensWithTypes(Not, In) {
			err = c.performSyntaxAnalysisAtLevel4()
			if err != nil {
				return err
			}

			c.addTokenToResult(NotIn, variants.Empty, token.Line(), token.Column())
		} else {
			break
		}
	}

	return nil
}

// Performs a syntax analysis at level 4.
func (c *ExpressionParser) performSyntaxAnalysisAtLevel4() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	err = c.performSyntaxAnalysisAtLevel5()
	if err != nil {
		return err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		if token.Type() == Star || token.Type() == Slash || token.Type() == Procent {
			c.moveToNextToken()

			err = c.performSyntaxAnalysisAtLevel5()
			if err != nil {
				return err
			}

			c.addTokenToResult(token.Type(), variants.Empty, token.Line(), token.Column())
			continue
		}
		break
	}

	return nil
}

// Performs a syntax analysis at level 5.
func (c *ExpressionParser) performSyntaxAnalysisAtLevel5() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	err = c.performSyntaxAnalysisAtLevel6()
	if err != nil {
		return err
	}

	for c.hasMoreTokens() {
		token := c.getCurrentToken()
		if token.Type() == Power || token.Type() == In ||
			token.Type() == ShiftLeft || token.Type() == ShiftRight {
			c.moveToNextToken()

			err = c.performSyntaxAnalysisAtLevel6()
			if err != nil {
				return err
			}

			c.addTokenToResult(token.Type(), variants.Empty, token.Line(), token.Column())
			continue
		}
		break
	}

	return nil
}

// Performs a syntax analysis at level 6.
func (c *ExpressionParser) performSyntaxAnalysisAtLevel6() error {
	err := c.checkForMoreTokens()
	if err != nil {
		return err
	}

	// Process unary '+' or '-'.
	unaryToken := c.getCurrentToken()
	if unaryToken.Type() == Plus {
		unaryToken = nil
		c.moveToNextToken()
	} else if unaryToken.Type() == Minus {
		unaryToken = NewExpressionToken(Unary, unaryToken.Value(), unaryToken.Line(), unaryToken.Column())
		c.moveToNextToken()
	} else {
		unaryToken = nil
	}

	err = c.checkForMoreTokens()
	if err != nil {
		return err
	}

	// Identify function calls.
	primitiveToken := c.getCurrentToken()
	nextToken := c.getNextToken()
	if primitiveToken.Type() == Variable &&
		nextToken != nil && nextToken.Type() == LeftBrace {
		primitiveToken = NewExpressionToken(Function, primitiveToken.Value(), primitiveToken.Line(), primitiveToken.Column())
	}

	if primitiveToken.Type() == Constant {
		c.moveToNextToken()
		c.addTokenToResult(primitiveToken.Type(), primitiveToken.Value(), primitiveToken.Line(), primitiveToken.Column())
	} else if primitiveToken.Type() == Variable {
		c.moveToNextToken()

		temp := primitiveToken.Value().AsString()
		found := false
		for _, v := range c.variableNames {
			if temp == v {
				found = true
				break
			}
		}
		if !found {
			c.variableNames = append(c.variableNames, temp)
		}

		c.addTokenToResult(primitiveToken.Type(), primitiveToken.Value(), primitiveToken.Line(), primitiveToken.Column())
	} else if primitiveToken.Type() == LeftBrace {
		c.moveToNextToken()

		err = c.performSyntaxAnalysis()
		if err != nil {
			return err
		}

		err = c.checkForMoreTokens()
		if err != nil {
			return err
		}

		primitiveToken = c.getCurrentToken()
		if primitiveToken.Type() != RightBrace {
			err = errors.NewSyntaxError(
				"",
				errors.ErrMissedCloseParenthesis,
				"Expected ')' was not found",
				primitiveToken.Line(),
				primitiveToken.Column(),
			)
			return err
		}

		c.moveToNextToken()
	} else if primitiveToken.Type() == Function {
		c.moveToNextToken()
		token := c.getCurrentToken()

		if token.Type() != LeftBrace {
			err = errors.NewSyntaxError(
				"",
				errors.ErrInternal,
				"Internal error",
				token.Line(),
				token.Column(),
			)
			return err
		}

		paramCount := 0
		for true {
			c.moveToNextToken()
			token = c.getCurrentToken()
			if token == nil || token.Type() == RightBrace {
				break
			}

			paramCount = paramCount + 1

			err = c.performSyntaxAnalysis()
			if err != nil {
				return err
			}

			token = c.getCurrentToken()

			if token == nil || token.Type() != Comma {
				break
			}
		}

		err = c.checkForMoreTokens()
		if err != nil {
			return err
		}

		if token.Type() != RightBrace {
			err = errors.NewSyntaxError("", errors.ErrMissedCloseParenthesis, "Expected ')' was not found", token.Line(), token.Column())
			return err
		}

		c.moveToNextToken()

		c.addTokenToResult(Constant, variants.VariantFromInteger(paramCount), primitiveToken.Line(), primitiveToken.Column())
		c.addTokenToResult(primitiveToken.Type(), primitiveToken.Value(), primitiveToken.Line(), primitiveToken.Column())
	} else {
		err = errors.NewSyntaxError("", errors.ErrErrorAt, "Syntax error at "+primitiveToken.Value().AsString(), primitiveToken.Line(), primitiveToken.Column())
		return err
	}

	if unaryToken != nil {
		c.addTokenToResult(unaryToken.Type(), variants.Empty, unaryToken.Line(), unaryToken.Column())
	}

	// Process [] operator.
	if c.hasMoreTokens() {
		primitiveToken = c.getCurrentToken()
		if primitiveToken.Type() == LeftSquareBrace {
			c.moveToNextToken()

			err = c.performSyntaxAnalysis()
			if err != nil {
				return err
			}

			err = c.checkForMoreTokens()
			if err != nil {
				return err
			}

			primitiveToken := c.getCurrentToken()
			if primitiveToken.Type() != RightSquareBrace {
				err = errors.NewSyntaxError("", errors.ErrMissedCloseSquareBracket, "Expected ']' was not found", primitiveToken.Line(), primitiveToken.Column())
			}

			c.moveToNextToken()
			c.addTokenToResult(Element, variants.Empty, 0, 0)
		}
	}

	return nil
}
