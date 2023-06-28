package commands

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

// CommandSet contains a set of commands and events supported by a commandable object.
// The CommandSet supports command interceptors to extend and the command call chain.
// CommandSets can be used as alternative commandable interface to a business object.
// It can be used to auto generate multiple external services for the business object without writing much code.
//
//	see Command
//	see Event
//	see ICommandable
//	Example:
//		type MyDataCommandSet struct {
//			*CommandSet
//			_controller IMyDataController
//		}
//
//		// Any data controller interface
//		func (dcs *MyDataCommandSet) CreateMyDataCommandSet(controller IMyDataController) {
//			dcs._controller = controller
//			dcs.AddCommand(dcs.makeGetMyDataCommand())
//		}
//		func (dcs *MyDataCommandSet) makeGetMyDataCommand() ICommand {
//			return NewCommand(
//				"get_mydata",
//				nil,
//				func(ctx context.Context, args *exec.Parameters) (any, err) {
//					var param = args.GetAsString("param")
//					return dcs._controller.GetMyData(ctx, param)
//				},
//			)
//		}
type CommandSet struct {
	commands       []ICommand
	events         []IEvent
	interceptors   []ICommandInterceptor
	commandsByName map[string]ICommand
	eventsByName   map[string]IEvent
}

// NewCommandSet creates an empty CommandSet object.
//
//	Returns: *CommandSet
func NewCommandSet() *CommandSet {
	return &CommandSet{
		commands:       []ICommand{},
		events:         []IEvent{},
		interceptors:   []ICommandInterceptor{},
		commandsByName: map[string]ICommand{},
		eventsByName:   map[string]IEvent{},
	}
}

// Commands gets all commands registered in this command set.
//
//	see ICommand
//	Returns: []ICommand a list of commands.
func (c *CommandSet) Commands() []ICommand {
	return c.commands
}

// Events gets all events registered in this command set.
//
//	see IEvent
//	Returns: []IEvent a list of events.
func (c *CommandSet) Events() []IEvent {
	return c.events
}

// FindCommand searches for a command by its name.
//
//	see ICommand
//	Parameters: commandName: string the name of the command to search for.
//	Returns: ICommand the command, whose name matches the provided name.
func (c *CommandSet) FindCommand(commandName string) ICommand {
	return c.commandsByName[commandName]
}

// FindEvent searches for an event by its name in this command set.
//
//	see IEvent
//	Parameters: eventName: string the name of the event to search for.
//	Returns: IEvent the event, whose name matches the provided name.
func (c *CommandSet) FindEvent(eventName string) IEvent {
	return c.eventsByName[eventName]
}

func (c *CommandSet) buildCommandChain(command ICommand) {
	next := command

	for i := len(c.interceptors) - 1; i >= 0; i-- {
		next = NewInterceptedCommand(c.interceptors[i], next)
	}

	c.commandsByName[next.Name()] = next
}

func (c *CommandSet) rebuildAllCommandChains() {
	c.commandsByName = map[string]ICommand{}

	for _, command := range c.commands {
		c.buildCommandChain(command)
	}
}

// AddCommand adds a command to this command set.
//
//	see ICommand
//	Parameters: command: ICommand the command to add.
func (c *CommandSet) AddCommand(command ICommand) {
	c.commands = append(c.commands, command)
	c.buildCommandChain(command)
}

// AddCommands adds multiple commands to this command set.
//
//	see ICommand
//	Parameters: []ICommand the array of commands to add.
func (c *CommandSet) AddCommands(commands []ICommand) {
	for _, command := range commands {
		c.AddCommand(command)
	}
}

// AddEvent adds an event to this command set.
//
//	see IEvent
//	Parameters: IEvent the event to add.
func (c *CommandSet) AddEvent(event IEvent) {
	c.events = append(c.events, event)
	c.eventsByName[event.Name()] = event
}

// AddEvents adds multiple events to this command set.
//
//	see IEvent
//	Parameters: []IEvent the array of events to add.
func (c *CommandSet) AddEvents(events []IEvent) {
	for _, event := range events {
		c.AddEvent(event)
	}
}

// AddCommandSet adds all the commands and events from specified command set into this one.
//
//	Parameters: commandSet: *CommandSet the CommandSet to add.
func (c *CommandSet) AddCommandSet(commandSet *CommandSet) {
	c.AddCommands(commandSet.Commands())
	c.AddEvents(commandSet.Events())
}

// AddListener фdds a listener to receive notifications on fired events.
//
//	see IEventListener
//	Parameters: listener: IEventListener the listener to add.
func (c *CommandSet) AddListener(listener IEventListener) {
	for _, event := range c.events {
		event.AddListener(listener)
	}
}

// RemoveListener removes previosly added listener.
//
//	see IEventListener
//	Parameters: IEventListener the listener to remove.
func (c *CommandSet) RemoveListener(listener IEventListener) {
	for _, event := range c.events {
		event.RemoveListener(listener)
	}
}

// AddInterceptor adds a command interceptor to this command set.
//
//	see ICommandInterceptor
//	Parameters: ICommandInterceptor the interceptor to add.
func (c *CommandSet) AddInterceptor(interceptor ICommandInterceptor) {
	c.interceptors = append(c.interceptors, interceptor)
	c.rebuildAllCommandChains()
}

// Execute a command specified by its name.
//
//	see ICommand
//	see Parameters
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- commandName: string the name of that command that is to be executed.
//		- args: Parameters the parameters (arguments) to pass to the command for execution.
//	Returns:
//		- result: any
//		- err: error
func (c *CommandSet) Execute(ctx context.Context, commandName string, args *exec.Parameters) (result any, err error) {
	cref := c.FindCommand(commandName)

	traceid := cctx.GetTraceId(ctx)

	if cref == nil {
		err := errors.NewBadRequestError(
			traceid,
			"CMD_NOT_FOUND",
			"Request command does not exist",
		).WithDetails("command", commandName)
		return nil, err
	}

	if traceid == "" {
		traceid = keys.IdGenerator.NextShort()
	}

	// Validate parameters
	results := cref.Validate(args)
	if len(results) > 0 {
		err := validate.NewValidationErrorFromResults(traceid, results, false)
		return nil, err
	}

	return cref.Execute(ctx, args)
}

// Validate args for command specified by its name using defined schema. If validation schema is
// not defined than the methods returns no errors. It returns validation error if the command is not found.
//
//	see Command
//	see Parameters
//	see ValidationResult
//	Parameters:
//		- commandName: string the name of the command for which the 'args' must be validated.
//		- args: Parameters the parameters (arguments) to validate.
//	Returns: []ValidationResult an array of ValidationResults.
//		If no command is found by the given name,
//		then the returned array of ValidationResults will contain a single entry,
//		whose type will be ValidationResultType.Error.
func (c *CommandSet) Validate(commandName string, args *exec.Parameters) []*validate.ValidationResult {
	cref := c.FindCommand(commandName)

	if cref == nil {
		return []*validate.ValidationResult{
			validate.NewValidationResult(
				"",
				validate.Error,
				"CMD_NOT_FOUND",
				"Requested command does not exist",
				nil,
				nil,
			),
		}
	}

	return cref.Validate(args)
}

// Notify fires event specified by its name and notifies all registered listeners
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- eventName: string the name of the event that is to be fired.
//		- args: Parameters the event arguments (parameters).
func (c *CommandSet) Notify(ctx context.Context, eventName string, args *exec.Parameters) {
	if event := c.FindEvent(eventName); event != nil {
		event.Notify(ctx, args)
	}
}
