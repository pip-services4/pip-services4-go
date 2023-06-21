package commands

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

// ICommand An interface for Commands, which are part of the Command design pattern.
// Each command wraps a method or function and allows to call them in uniform and safe manner.
type ICommand interface {
	exec.IExecutable
	// Name gets the command name.
	//	Returns: string the command name.
	Name() string
	// Validate validates command arguments before execution using defined schema.
	//	see Parameters
	//	see ValidationResult
	//	Parameters: args: Parameters the parameters (arguments) to validate.
	//	Returns: ValidationResult[] an array of ValidationResults.
	Validate(args *exec.Parameters) []*validate.ValidationResult
}
