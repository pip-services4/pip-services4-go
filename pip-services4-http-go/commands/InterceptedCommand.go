package commands

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

// InterceptedCommand implements a command wrapped by an interceptor.
// It allows building command call chains.
// The interceptor can alter execution and delegate calls to a
// next command, which can be intercepted or concrete.
//
//	see ICommand
//	see ICommandInterceptor
//	Example:
//		type CommandLogger struct {
//			msg string
//		}
//
//		func (cl * CommandLogger) Name(command ICommand) string {
//			return command.Name();
//		}
//
//		func (cl * CommandLogger) Execute(ctx context.Context, command ICommand, args Parameters) (res any, err error){
//			fmt.Println("Executed command " + command.Name());
//			return command.Execute(ctx, args);
//		}
//
//		func (cl * CommandLogger) Validate(command ICommand, args Parameters) []*ValidationResult {
//			return command.Validate(args);
//		}
//
//		logger := CommandLogger{mgs:"CommandLogger"};
//		loggedCommand = NewInterceptedCommand(logger, command);
//
//		// Each called command will output: Executed command <command name>
type InterceptedCommand struct {
	interceptor ICommandInterceptor
	next        ICommand
}

// NewInterceptedCommand creates a new InterceptedCommand, which serves as a link in an execution chain.
// Contains information about the interceptor that is being used and the next command in the chain.
//
//	Parameters:
//		- interceptor: ICommandInterceptor the interceptor that is intercepting the command.
//		- next: ICommand (link to) the next command in the command's execution chain.
//	Returns: *InterceptedCommand
func NewInterceptedCommand(interceptor ICommandInterceptor, next ICommand) *InterceptedCommand {
	return &InterceptedCommand{
		interceptor: interceptor,
		next:        next,
	}
}

// Name Returns string the name of the command that is being intercepted.
func (c *InterceptedCommand) Name() string {
	return c.interceptor.Name(c.next)
}

// Execute the next command in the execution chain using the given parameters (arguments).
// see Parameters
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- args: Parameters the parameters (arguments) to pass to the command for execution.
//	Returns:
//		- err: error
//		- result: any
func (c *InterceptedCommand) Execute(ctx context.Context, args *exec.Parameters) (result any, err error) {
	return c.interceptor.Execute(ctx, c.next, args)
}

// Validate the parameters (arguments) that are to be passed to the command that is next in the execution chain.
//
//	see Parameters
//	see ValidationResult
//	Parameters: args the parameters (arguments) to validate for the next command.
//	Returns: []*ValidationResult an array of *ValidationResults.
func (c *InterceptedCommand) Validate(args *exec.Parameters) []*validate.ValidationResult {
	return c.interceptor.Validate(c.next, args)
}
