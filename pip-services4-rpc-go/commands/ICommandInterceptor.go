package commands

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

// ICommandInterceptor An interface for stackable command intercepters,
// which can extend and modify the command call chain.
// This mechanism can be used for authentication, logging, and other functions.
//
//	see ICommand
//	see InterceptedCommand
type ICommandInterceptor interface {
	// Name gets the name of the wrapped command.
	// The interceptor can use this method to override the command name.
	// Otherwise, it shall just delegate the call to the wrapped command.
	//	Parameters:  command: ICommand the next command in the call chain.
	//	Returns: string the name of the wrapped command.
	Name(command ICommand) string

	// Execute the wrapped command with specified arguments.
	// The interceptor can use this method to intercept and alter the command execution.
	// Otherwise, it shall just delete the call to the wrapped command.
	//	see Parameters
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- command: ICommand the next command in the call chain that is to be executed.
	//		- args: Parameters the function that is to be called once execution is complete.
	//			If an exception is raised, then it will be called with the error.
	//	Returns:
	//		- result: any
	//		- err: error
	Execute(ctx context.Context, command ICommand, args *exec.Parameters) (any, error)

	// Validate arguments of the wrapped command before its execution.
	// The interceptor can use this method to intercept and alter validation of the command arguments.
	// Otherwise, it shall just delegate the call to the wrapped command.
	//	see Parameters
	//	see ValidationResult
	//	Parameters:
	//		- command: ICommand the next command in the call chain to be validated against.
	//		- args: Parameters the parameters (arguments) to validate.
	//	Returns: []*ValidationResult an array of *ValidationResults.
	Validate(command ICommand, args *exec.Parameters) []*validate.ValidationResult
}
