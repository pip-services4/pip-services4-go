package context

import (
	"context"
	"errors"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

// ContextShutdownChan a channel to send
// default context feedback
type ContextShutdownChan chan int8

// ContextShutdownWithErrorChan a channel to send
// context feedback with error
type ContextShutdownWithErrorChan chan error

// ContextFeedbackWithCustomDataChan a channel to
// send context feedback with specific data.
type ContextFeedbackWithCustomDataChan[T any] chan T

// ContextValueType an enum to describe specific
// context feedback channel
//
//	Possible values:
//		- ContextShutdownChanType
//		- ContextShutdownWithErrorChanType
//		- ContextFeedbackChanWithCustomDataType
type ContextValueType string

const (
	ContextShutdownChanType               ContextValueType = "pip.ContextShutdownChan"
	ContextShutdownWithErrorChanType      ContextValueType = "pip.ContextShutdownWithErrorChan"
	ContextFeedbackWithCustomDataChanType ContextValueType = "pip.ContextFeedbackWithCustomDataChan"
)

const DefaultCancellationSignal int8 = 1

// AddShutdownChanToContext wrap context with ContextFeedbackChan
//
//	see context.WithValue
//	Parameters:
//		- context.Context parent context
//		- ContextShutdownChan - channel to put into context
//	Returns:
//		- context.Context is a context with value
//		- bool true if channel is not nil or false
func AddShutdownChanToContext(ctx context.Context, channel ContextShutdownChan) (context.Context, bool) {
	if channel == nil {
		return ctx, false
	}
	return context.WithValue(ctx, ContextShutdownChanType, channel), true
}

// AddErrShutdownChanToContext wrap context with ContextFeedbackChanWithError
//
//	see context.WithValue
//	Parameters:
//		- context.Context - parent context
//		- ContextShutdownWithErrorChan - channel to put into context
//	Returns:
//		- context.Context is a context with value
//		- bool true if channel is not nil or false
func AddErrShutdownChanToContext(ctx context.Context, channel ContextShutdownWithErrorChan) (context.Context, bool) {
	if channel == nil {
		return ctx, false
	}
	return context.WithValue(ctx, ContextShutdownWithErrorChanType, channel), true
}

// AddCustomDataChanToContext wrap context with ContextFeedbackChanWithCustomData
//
//	T is a custom data type
//	see context.WithValue
//	Parameters:
//		- context.Context - parent context
//		- ContextFeedbackWithCustomDataChan[T] - channel to put into context
//	Returns:
//		- context.Context is a context with value
//		- bool true if channel is not nil or false
func AddCustomDataChanToContext[T any](ctx context.Context, channel ContextFeedbackWithCustomDataChan[T]) (context.Context, bool) {
	if channel == nil {
		return ctx, false
	}
	return context.WithValue(ctx, ContextFeedbackWithCustomDataChanType, channel), true
}

// SendShutdownSignal sends interrupt signal up to the context owner
//
//	Parameters: context.Context is a current context
//	Returns: bool true if signal sends successful or false
func SendShutdownSignal(ctx context.Context) bool {
	if val := ctx.Value(ContextShutdownChanType); val != nil {
		if _chan, ok := val.(ContextShutdownChan); ok {
			select {
			case _chan <- DefaultCancellationSignal:
				return true
			default:
				return false
			}
		}
	}
	return false
}

// SendShutdownSignalWithErr sends error and interrupt signal up to the context owner
//
//	Parameters:
//		- context.Context is a current context
//		- error
//	Returns: bool true if signal sends successful or false
func SendShutdownSignalWithErr(ctx context.Context, err error) bool {
	if val := ctx.Value(ContextShutdownWithErrorChanType); val != nil {
		if _chan, ok := val.(ContextShutdownWithErrorChan); ok {
			select {
			case _chan <- err:
				return true
			default:
				return false
			}
		}
	}
	return false
}

// SendSignalWithCustomData sends custom data and interrupt signal up to the context owner
//
//	Parameters:
//		- context.Context is a current context
//		- T custom data
//	Returns: bool true if signal sends successful or false
func SendSignalWithCustomData[T any](ctx context.Context, data T) bool {
	if val := ctx.Value(ContextFeedbackWithCustomDataChanType); val != nil {
		if _chan, ok := val.(ContextFeedbackWithCustomDataChan[T]); ok {
			select {
			case _chan <- data:
				return true
			default:
				return false
			}
		}
	}
	return false
}

// DefaultErrorHandlerWithShutdown is a default error handler method which catch panic,
// parse error and send shutdown signal to main container.
// Gracefully shutdown all containers. Using only with defer operator!
//
//	see SendShutdownSignalWithErr
//	Examples:
//		func MyFunc(ctx context.Context) {
//			defer DefaultErrorHandlerWithShutdown(ctx)
//			...
//			panic("some error")
//		}
//	Parameters:
//		- ctx context.Context
func DefaultErrorHandlerWithShutdown(ctx context.Context) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			msg := convert.StringConverter.ToString(r)
			err = errors.New(msg)
		}
		SendShutdownSignalWithErr(ctx, err)
	}
}

func NewContextWithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, utils.TRACE_ID, traceId)
}

func GetTraceId(ctx context.Context) string {
	traceId := ctx.Value(utils.TRACE_ID)

	if traceId == nil || traceId == "" {
		traceId = ctx.Value("trace_id")
		if traceId == nil || traceId == "" {
			traceId = ctx.Value("traceId")
		}
	}

	if val, ok := traceId.(string); ok {
		return val
	} else {
		return ""
	}
}

func GetClient(ctx context.Context) string {
	client := ctx.Value(utils.CLIENT)

	if client == nil || client == "" {
		client = ctx.Value("client")
	}

	if val, ok := client.(string); ok {
		return val
	} else {
		return ""
	}
}

func GetUser(ctx context.Context) any {
	user := ctx.Value(utils.TRACE_ID)

	if user == nil || user == "" {
		user = ctx.Value("user")
	}

	return user
}
