package errors

// Defines standard error categories to application exceptions supported by PipServices toolkit.
//
//	BadRequest - Errors due to incorrectly specified invocation parameters.
//	For example: missing or incorrect parameters.
//
//	Conflict - Errors raised by conflicts between object versions that were posted by the user and those
//	that are stored on the server.
//
//	FailedInvocation - Errors caused by remote calls failed due to unidenfied reasons.
//
//	FileError - Errors in read/write local disk operations.
//
//	Internal - Internal errors caused by programming mistakes.
//
//	InvalidState - Errors caused by incorrect object state..
//	For example: business calls when the component is not ready.
//
//	Misconfiguration - Errors related to mistakes in user-defined configurations.
//
//	NoResponse - Errors caused by remote calls timeouted and not returning results. It allows to clearly separate communication related problems from other application errors.
//
//	NotFound - Errors caused by attempts to access missing objects.
//
//	Unauthorized - Access errors caused by missing user identity (authentication error) or incorrect security permissions (authorization error).
//
//	Unknown - Unknown or unexpected errors.
//
//	Unsupported - Errors caused by calls to unsupported or not yet implemented functionality.
const (
	Unknown          = "Unknown"
	Internal         = "Internal"
	Misconfiguration = "Misconfiguration"
	InvalidState     = "InvalidState"
	NoResponse       = "NoResponse"
	FailedInvocation = "FailedInvocation"
	FileError        = "FileError"
	BadRequest       = "BadRequest"
	Unauthorized     = "Unauthorized"
	NotFound         = "NotFound"
	Conflict         = "Conflict"
	Unsupported      = "Unsupported"
)
