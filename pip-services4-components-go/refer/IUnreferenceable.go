package refer

import "context"

// IUnreferenceable Interface for components that require explicit clearing of references to dependent components.
//	see IReferences
//	see IReferenceable
//	Example
//		type MyController  {
//			_persistence IMyPersistence;
//		}
//		func (mc* MyController) SetReferences(ctx context.Context, references *IReferences) {
//			mc._persistence = references.GetOneRequired(
//				NewDescriptor("mygroup", "persistence", "*", "*", "1.0"),
//			);
//		}
//
//		func (mc* MyController) UnsetReferences(ctx context.Context) {
//			mc._persistence = nil;
//		}
type IUnreferenceable interface {
	// UnsetReferences (clears) previously set references to dependent components.
	UnsetReferences(ctx context.Context)
}
