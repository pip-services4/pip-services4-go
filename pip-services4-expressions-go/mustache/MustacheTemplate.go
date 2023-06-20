package mustache

import (
	"strings"

	merr "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache/errors"
	mparsers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache/parsers"
	tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
)

// MustacheTemplate mustache templating engine
type MustacheTemplate struct {
	defaultVariables map[string]string
	parser           *mparsers.MustacheParser
	autoVariables    bool
}

// NewMustacheTemplate creates a new instance of mustache template
func NewMustacheTemplate() *MustacheTemplate {
	c := MustacheTemplate{
		defaultVariables: map[string]string{},
		parser:           mparsers.NewMustacheParser(),
		autoVariables:    true,
	}
	return &c
}

// NewMustacheTemplateFromString creates and initializes an instance of mustache template
//
//	Parameters:
//		template: the mustache template
func NewMustacheTemplateFromString(template string) (*MustacheTemplate, error) {
	c := NewMustacheTemplate()
	err := c.SetTemplate(template)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Template gets the mustache template.
func (c *MustacheTemplate) Template() string {
	return c.parser.Template()
}

// SetTemplate sets the mustache template.
func (c *MustacheTemplate) SetTemplate(value string) error {
	err := c.parser.SetTemplate(value)
	if err != nil {
		return err
	}
	if c.autoVariables {
		c.CreateVariables(&c.defaultVariables)
	}
	return nil
}

func (c *MustacheTemplate) OriginalTokens() []*tokenizers.Token {
	return c.parser.OriginalTokens()
}

func (c *MustacheTemplate) SetOriginalTokens(value []*tokenizers.Token) error {
	err := c.parser.SetOriginalTokens(value)
	if err != nil {
		return err
	}

	if c.autoVariables {
		c.CreateVariables(&c.defaultVariables)
	}
	return nil
}

// AutoVariables gets the flag to turn on auto creation of variables for specified mustache.
func (c *MustacheTemplate) AutoVariables() bool {
	return c.autoVariables
}

// SetAutoVariables sets the flag to turn on auto creation of variables for specified mustache.
func (c *MustacheTemplate) SetAutoVariables(value bool) {
	c.autoVariables = value
}

// DefaultVariables gets the list with default variables.
func (c *MustacheTemplate) DefaultVariables() map[string]string {
	return c.defaultVariables
}

// SetDefaultVariables sets the list with default variables.
func (c *MustacheTemplate) SetDefaultVariables(value map[string]string) {
	c.defaultVariables = value
}

// InitialTokens gets the list of original mustache tokens.
func (c *MustacheTemplate) InitialTokens() []*mparsers.MustacheToken {
	return c.parser.InitialTokens()
}

// Gets the list of processed mustache tokens.
func (c *MustacheTemplate) ResultTokens() []*mparsers.MustacheToken {
	return c.parser.ResultTokens()
}

// GetVariable gets a variable value from the collection of variables
//
//	Parameters:
//		- variables: a collection of variables.
//		- name: a variable name to get.
//	Returns: a variable value or <code>undefined</code>
func (c *MustacheTemplate) GetVariable(variables map[string]string, name string) *string {
	if variables == nil || name == "" {
		return nil
	}

	name = strings.ToLower(name)
	var result *string = nil

	for propName, propValue := range variables {
		if strings.ToLower(propName) == name {
			result = &propValue
			break
		}
	}

	return result
}

// CreateVariables populates the specified variables list with variables from parsed mustache.
//
//	Parameters:
//		variables: The list of variables to be populated.
func (c *MustacheTemplate) CreateVariables(variables *map[string]string) {
	if variables == nil {
		return
	}

	for _, variableName := range c.parser.VariableNames() {
		found := c.GetVariable(*variables, variableName) != nil
		if !found {
			(*variables)[variableName] = ""
		}
	}
}

// Clear cleans up this calculator from all data.
func (c *MustacheTemplate) Clear() {
	c.parser.Clear()
	c.defaultVariables = map[string]string{}
}

// Evaluate this mustache template using default variables.
//
//	Returns: the evaluated template
func (c *MustacheTemplate) Evaluate() (string, error) {
	return c.EvaluateWithVariables(nil)
}

// EvaluateWithVariables evaluates this mustache using specified variables.
//
//	Parameters:
//		variables: The collection of variables
//	Returns: the evaluated template
func (c *MustacheTemplate) EvaluateWithVariables(variables map[string]string) (string, error) {
	if variables == nil {
		variables = c.defaultVariables
	}

	return c.evaluateTokens(c.parser.ResultTokens(), variables)
}

func (c *MustacheTemplate) isDefinedVariable(variables map[string]string, name string) bool {
	value := c.GetVariable(variables, name)
	if value == nil {
		return false
	}
	return (*value) != ""
}

func (c *MustacheTemplate) escapeString(value string) string {
	if value == "" {
		return ""
	}

	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "/", "\\/")
	value = strings.ReplaceAll(value, "\b", "\\b")
	value = strings.ReplaceAll(value, "\f", "\\f")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\t", "\\t")

	return value
}

func (c *MustacheTemplate) evaluateTokens(tokens []*mparsers.MustacheToken, variables map[string]string) (string, error) {
	if tokens == nil {
		return "", nil
	}

	builder := strings.Builder{}

	for _, token := range tokens {
		switch token.Type() {
		case mparsers.TokenComment:
			// Skip
			break
		case mparsers.TokenValue:
			builder.WriteString(token.Value())
			break
		case mparsers.TokenVariable:
			value := c.GetVariable(variables, token.Value())
			if value != nil {
				builder.WriteString(*value)
			}
			break
		case mparsers.TokenEscapedVariable:
			value := c.GetVariable(variables, token.Value())
			if value != nil {
				escValue := c.escapeString(*value)
				builder.WriteString(escValue)
			}
			break
		case mparsers.TokenSection:
			defined := c.isDefinedVariable(variables, token.Value())
			if defined && token.Tokens() != nil {
				value, err := c.evaluateTokens(token.Tokens(), variables)
				if err != nil {
					return "", err
				}
				builder.WriteString(value)
			}
			break
		case mparsers.TokenInvertedSection:
			defined := c.isDefinedVariable(variables, token.Value())
			if !defined && token.Tokens() != nil {
				value, err := c.evaluateTokens(token.Tokens(), variables)
				if err != nil {
					return "", err
				}
				builder.WriteString(value)
			}
			break
		case mparsers.TokenPartial:
			err := merr.NewMustacheError("", "PARTIALS_NOT_SUPPORTED", "Partials are not supported", token.Line(), token.Column())
			return "", err
		default:
			err := merr.NewMustacheError("", "INTERNAL", "Internal error", token.Line(), token.Column())
			return "", err
		}
	}

	return builder.String(), nil
}
