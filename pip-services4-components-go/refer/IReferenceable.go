package refer

import "context"

// IReferenceable interface for components that depends on other components.
// If component requires explicit notification to unset references it shall
// additionally implement IUnreferenceable interface.
//	see IReferences
//	see IUnreferenceable
//	see Referencer
//	Example
//		type MyController {
//			_persistence IPersistence
//		}
//		func (mc* MyController) SetReferences(ctx context.Context, references IReferences) {
//			mc._persistence = references.GetOneRequired(
//				NewDescriptor("mygroup", "persistence", "*", "*", "1.0"))
//			);
//		}
//		...
type IReferenceable interface {
	// SetReferences sets references to dependent components.
	//	see IReferences
	//	Parameters:
	//		- ctx context.Context
	//		- references IReferences references to locate the component dependencies.
	SetReferences(ctx context.Context, references IReferences)
}
