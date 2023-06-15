// Portable and localizable Exceptions classes.
// Each Exception, in addition to a description and stack trace has a unique string code,
// details array (which can be used for creating localized strings).
//
// Way to use:
//
//	An existing exception class can be used.
//	A child class that extends ApplicationException can we written.
//	A exception can be wrapped around (into?) an existing application exception.
//	Exceptions are serializable.
//		The exception classes themselves are not serializable,
//		but they can be converted to ErrorDescriptions, which are serializable in one language,
//		transferred to the receiving side, and deserialized in another language.
//		After deserialization, the initial exception class can be restored.
//
//	Additionally: when transferring an exception from one language to another,
//		the exception type that is closest to the initial exception type is
//		chosen from the exceptions available in the target language.

package errors
