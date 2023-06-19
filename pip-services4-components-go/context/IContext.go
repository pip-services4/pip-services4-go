package context

type IContext interface {
	// Configure configures component by passing configuration parameters.
	//
	//	Parameters:
	//		- key string a key of the element to get.
	// Returns: any the value of the map element.
	Get(key string) any

	// Gets a trace (trace) id.
	//
	// Returns: a trace id or empty string if it is not defined.
	GetTraceId() string

	// Gets a client name.
	//
	// Returns: a client name or <code>null</code> if it is not defined.
	GetClient() string

	// Gets a reference to user object.
	//
	// Returns: a user reference or empty string if it is not defined.
	GetUser() string
}
