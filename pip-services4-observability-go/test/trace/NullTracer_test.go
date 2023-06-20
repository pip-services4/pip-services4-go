package test_tracer

import (
	"context"
	"errors"
	"testing"

	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
)

func newNullTracer() *ctrace.NullTracer {
	return ctrace.NewNullTracer()
}

func TestSimpleNullTracing(t *testing.T) {
	tracer := newNullTracer()
	tracer.Trace(context.Background(), "mycomponent", "mymethod", 123456)
	tracer.Failure(context.Background(), "mycomponent", "mymethod", errors.New("Test error"), 123456)
}

func TestTraceNullTiming(t *testing.T) {
	tracer := newNullTracer()
	timing := tracer.BeginTrace(context.Background(), "mycomponent", "mymethod")
	timing.EndTrace()

	timing = tracer.BeginTrace(context.Background(), "mycomponent", "mymethod")
	timing.EndFailure(errors.New("Test error"))
}
