package test_commands

import (
	"context"
	"testing"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	"github.com/stretchr/testify/assert"
)

type TestListener struct{}

func (c *TestListener) OnEvent(ctx context.Context, e commands.IEvent, value *exec.Parameters) {
	if cctx.GetTraceId(ctx) == "wrongId" {
		panic("Test error")
	}
}

func TestGetEventName(t *testing.T) {
	event := commands.NewEvent("name")

	assert.NotNil(t, event)
	assert.Equal(t, "name", event.Name())
}

func TestEventNotify(t *testing.T) {
	event := commands.NewEvent("name")

	listener := &TestListener{}
	event.AddListener(listener)
	assert.Equal(t, 1, len(event.Listeners()))

	event.Notify(context.Background(), nil)

	event.RemoveListener(listener)
	assert.Equal(t, 0, len(event.Listeners()))
}
