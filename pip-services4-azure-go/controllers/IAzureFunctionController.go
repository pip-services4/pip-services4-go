package controllers

// An interface that allows to integrate Google Function services into Google Function containers
// and connect their actions to the function calls.
type IAzureFunctionController interface {

	// Get all actions supported by the service.
	// Returns an array with supported actions.
	GetActions() []*AzureFunctionAction
}
