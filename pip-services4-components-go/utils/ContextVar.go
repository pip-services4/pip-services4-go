package utils

type ContextVar string

var (
	TRACE_ID ContextVar = "trace_id"
	CLIENT   ContextVar = "client"
	USER     ContextVar = "user"
)
