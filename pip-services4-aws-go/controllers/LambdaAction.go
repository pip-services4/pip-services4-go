package controllers

import (
	"context"

	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type LambdaAction struct {

	// Command to call the action
	Cmd string

	// Schema to validate action parameters
	Schema *cvalid.Schema

	// Action to be executed
	Action func(ctx context.Context, params map[string]any) (any, error)
}
