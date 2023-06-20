package test_tracer

import (
	"context"
	"errors"
	"testing"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
)

func newLogTracer() *ctrace.LogTracer {
	tracer := trace.NewLogTracer()
	ctx := context.Background()
	tracer.SetReferences(
		ctx,
		cref.NewReferencesFromTuples(
			ctx,
			cref.NewDescriptor("pip-services", "logger", "null", "default", "1.0"),
			clog.NewNullLogger()))
	return tracer
}

func TestSimpleTracing(t *testing.T) {
	tracer := newLogTracer()
	tracer.Trace(context.Background(), "mycomponent", "mymethod", 123456)
	tracer.Failure(context.Background(), "mycomponent", "mymethod", errors.New("Test error"), 123456)
}

func TestTraceTiming(t *testing.T) {
	tracer := newLogTracer()
	var timing = tracer.BeginTrace(context.Background(), "mycomponent", "mymethod")
	timing.EndTrace()

	timing = tracer.BeginTrace(context.Background(), "mycomponent", "mymethod")
	timing.EndFailure(errors.New("Test error"))
}
