package commands

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

// Command Concrete implementation of ICommand interface. Command allows to call a
// method or function using Command pattern.
//
//	Example:
//		command := NewCommand("add", null, func (ctx context.Context args *exec.Parameters)(any, err) {
//			param1 := args.getAsFloat("param1");
//			param2 := args.getAsFloat("param2");
//			return (param1 + param2), nil;
//		});
//
//		result, err := command.Execute(ctx, Parameters.NewParametersFromTuples("param1", 2, "param2", 2))
//		if (err) {
//			fmt.Println(err)
//		} else {
//			fmt.Println("2 + 2 = " + result)
//		}
//		// Console output: 2 + 2 = 4
type Command struct {
	schema validate.ISchema
	action func(ctx context.Context, args *exec.Parameters) (any, error)
	name   string
}

// NewCommand creates a new command object and assigns it's parameters.
//
//	Parameters
//		- name: string - the command name.
//		- schema: validate.ISchema the schema to validate command arguments.
//		- function: func(ctx context.Context, args *exec.Parameters) (any, error)
//			the function to be executed by this command.
//	Returns: *Command
func NewCommand(name string, schema validate.ISchema,
	action func(ctx context.Context, args *exec.Parameters) (any, error)) *Command {

	if name == "" {
		panic("Name cannot be empty")
	}
	if action == nil {
		panic("Action cannot be nil")
	}

	return &Command{
		name:   name,
		schema: schema,
		action: action,
	}
}

// Name gets the command name.
//
//	Returns: string - the name of this command.
func (c *Command) Name() string {
	return c.name
}

// Execute the command. Before execution, it validates args using the defined schema.
// The command execution intercepts exceptions raised by
// the called function and returns them as an error in callback.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- args: exec.Parameters - the parameters (arguments) to pass to this command for execution.
//	Returns: (any, error)
func (c *Command) Execute(ctx context.Context, args *exec.Parameters) (any, error) {
	if c.schema != nil {
		if err := c.schema.ValidateAndReturnError(utils.ContextHelper.GetTraceId(ctx), args, false); err != nil {
			return nil, err
		}
	}

	var err error

	// Execute in inner function to capture errors
	result, err2 := func() (any, error) {
		// Intercepting unhandled errors
		defer func() {
			if r := recover(); r != nil {
				tempMessage := convert.StringConverter.ToString(r)
				tempError := errors.NewInvocationError(
					utils.ContextHelper.GetTraceId(ctx),
					"EXEC_FAILED",
					"Execution "+c.Name()+" failed: "+tempMessage,
				).WithDetails("command", c.Name())

				if cause, ok := r.(error); ok {
					tempError.WithCause(cause)
				}

				err = tempError
			}
		}()

		return c.action(ctx, args)
	}()

	if err2 != nil {
		err = err2
	}

	return result, err
}

// Validate the command args before execution using the defined schema.
//
//	Parameters: args: exec.Parameters - the parameters
//		(arguments) to validate using this command's schema.
//	Returns: []*validate.ValidationResult an array of
//		ValidationResults or an empty array (if no schema is set).
func (c *Command) Validate(args *exec.Parameters) []*validate.ValidationResult {
	if c.schema != nil {
		results := c.schema.Validate(args)
		if results == nil {
			results = []*validate.ValidationResult{}
		}
		return results
	}

	return []*validate.ValidationResult{}
}

// GetSchema methods return validation schema for this command
func (c *Command) GetSchema() validate.ISchema {
	return c.schema
}
