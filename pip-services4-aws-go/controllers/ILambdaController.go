package controllers

// An interface that allows to integrate lambda services into lambda function containers
// and connect their actions to the function calls.
type ILambdaController interface {
	// Get all actions supported by the service.
	// Returns an array with supported actions.
	GetActions() []*LambdaAction
}
