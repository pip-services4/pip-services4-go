package exec

import "context"

// IParameterized interface for components that require execution parameters.
type IParameterized interface {
	// SetParameters sets execution parameters.
	//	Parameters:
	//		- ctx context.Context
	//		- parameters *Parameters execution parameters.
	SetParameters(ctx context.Context, parameters *Parameters)
}
