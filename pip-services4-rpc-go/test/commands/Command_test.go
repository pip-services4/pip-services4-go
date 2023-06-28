package test_commands

import (
	"context"
	"testing"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	"github.com/stretchr/testify/assert"
)

func commandExec(ctx context.Context, args *exec.Parameters) (any, error) {
	if cctx.GetTraceId(ctx) == "wrongId" {
		panic("Test error")
	}

	return nil, nil
}

func TestGetCommandName(t *testing.T) {
	command := commands.NewCommand("name", nil, commandExec)

	// Check match by individual fields
	assert.NotNil(t, command)
	assert.Equal(t, "name", command.Name())
}

func TestExecuteCommand(t *testing.T) {
	command := commands.NewCommand("name", nil, commandExec)

	_, err := command.Execute(context.Background(), nil)
	assert.Nil(t, err)

	_, err = command.Execute(cctx.NewContextWithTraceId(context.Background(), "wrongId"), nil)
	assert.NotNil(t, err)
}
