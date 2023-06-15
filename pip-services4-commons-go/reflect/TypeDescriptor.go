package reflect

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
)

// TypeDescriptor is a descriptor that points to specific object
// type by it's name and optional library (or module) where this type is defined.
// This class has symmetric implementation across all
// languages supported by Pip.Services toolkit and used to support dynamic data processing.
type TypeDescriptor struct {
	name string
	pkg  string
}

// NewTypeDescriptor creates a new instance of the type descriptor and sets its values.
// Parameters:
//   - name string a name of the object type.
//   - library string a library or module where this object type is implemented.
//     Returns: *TypeDescriptor
func NewTypeDescriptor(name string, pkg string) *TypeDescriptor {
	return &TypeDescriptor{
		name: name,
		pkg:  pkg,
	}
}

// Name get the name of the object type.
//
//	Returns: string the name of the object type.
func (c *TypeDescriptor) Name() string {
	return c.name
}

// Package gets the name of the package or module where the object type is defined.
// Returns: string the name of the package or module.
func (c *TypeDescriptor) Package() string {
	return c.pkg
}

// Equals compares this descriptor to a value. If the value is also a
// TypeDescriptor it compares their name and library fields.
// Otherwise, this method returns false.
//
//	Parameters:
//		- obj any a value to compare.
//	Returns: bool true if value is identical TypeDescriptor and false otherwise.
func (c *TypeDescriptor) Equals(descriptor *TypeDescriptor) bool {
	if descriptor == nil {
		return false
	}
	if strings.Compare(c.name, descriptor.name) != 0 {
		return false
	}
	if strings.Compare(c.pkg, descriptor.pkg) == 0 {
		return true
	}

	return false
}

// String gets a string representation of the object. The result has format name[,package]
//
//	Returns: string a string representation of the object.
func (c *TypeDescriptor) String() string {
	builder := strings.Builder{}

	builder.WriteString(c.name)

	if c.pkg != "" {
		builder.WriteString(",")
		builder.WriteString(c.pkg)
	}

	return builder.String()
}

// ParseTypeDescriptorFromString parses a string to get descriptor fields and returns them as a Descriptor.
// The string must have format name[,package]
// throws a ConfigError if the descriptor string is of a wrong format.
//
//	Parameters:
//		- value string a string to parse.
//	Returns: *TypeDescriptor a newly created Descriptor.
func ParseTypeDescriptorFromString(value string) (*TypeDescriptor, error) {
	if value == "" {
		return nil, nil
	}

	tokens := strings.Split(value, ",")

	if len(tokens) == 1 {
		return NewTypeDescriptor(strings.Trim(tokens[0], " "), ""), nil
	} else if len(tokens) == 2 {
		return NewTypeDescriptor(strings.Trim(tokens[0], " "), strings.Trim(tokens[1], " ")), nil
	}

	return nil, errors.NewConfigError(
		"",
		"BAD_DESCRIPTOR",
		"Type descriptor "+value+" is in wrong format",
	)
}
