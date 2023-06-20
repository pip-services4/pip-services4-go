package count

import "context"

// Interface for a callback to end measurement of execution elapsed time.

// ITimingCallback ends measurement of execution elapsed time and updates specified counter.
//	see Timing.EndTiming
//	Parameters:
//		- ctx context.Context
//		- name string a counter name
//		- elapsed float32 execution elapsed time in milliseconds to update the counter.
type ITimingCallback interface {
	EndTiming(ctx context.Context, name string, elapsed float64)
}
