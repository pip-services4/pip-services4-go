package utils

import (
	"context"
)

var ContextHelper = _TContextHelp{}

type _TContextHelp struct {
}

func (c *_TContextHelp) NewContextWithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, TRACE_ID, traceId)
}

func (c *_TContextHelp) GetTraceId(ctx context.Context) string {
	traceId := ctx.Value(TRACE_ID)

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

func (c *_TContextHelp) GetClient(ctx context.Context) string {
	client := ctx.Value(CLIENT)

	if client == nil || client == "" {
		client = ctx.Value("client")
	}

	if val, ok := client.(string); ok {
		return val
	} else {
		return ""
	}
}

func (c *_TContextHelp) GetUser(ctx context.Context) any {
	user := ctx.Value(TRACE_ID)

	if user == nil || user == "" {
		user = ctx.Value("user")
	}

	return user
}
