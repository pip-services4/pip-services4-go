package commands

// ICommandable an interface for commandable objects, which are part of the command design pattern.
// The commandable object exposes its functonality as commands and events groupped into a CommandSet.
// This interface is typically implemented by controllers and is used to auto generate external interfaces.
//	Example:
//		type MyDataController {
//			_commandSet  CommandSet;
//		}
//		func (dc *MyDataController) getCommandSet() CommandSet {
//			if (dc._commandSet == nil) {
//				dc._commandSet = NewDataCommandSet();
//			}
//			return dc._commandSet;
//		}
type ICommandable interface {
	// GetCommandSet gets a command set with all supported commands and events.
	//	see CommandSet
	//	Returns: *CommandSet a command set with commands and events.
	GetCommandSet() *CommandSet
}
