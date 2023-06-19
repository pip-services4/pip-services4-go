package exec

import "context"

// IExecutable interface for components that can be called to execute work.
//	Example:
//		type EchoComponent {}
//		...
//		func  (ec* EchoComponent) Execute(ctx context.Context, args *Parameters) (result any, err error) {
//			return nil, result = args.getAsObject("message")
//		}
//		echo := EchoComponent{};
//		message = "Test";
//		res, err = echo.Execute(ctx, NewParametersFromTuples("message", message));
//		fmt.Println(res);
type IExecutable interface {
	// Execute component with arguments and receives execution result.
	//	Parameters:
	//		- ctx context.Context a context to trace execution through call chain.
	//		- args *Parameters execution arguments.
	//	Returns: any, error result or execution and error
	Execute(ctx context.Context, args *Parameters) (result any, err error)
}
