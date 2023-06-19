package test_test

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/test"
	"github.com/stretchr/testify/assert"
)

func TestShutdown(t *testing.T) {
	sd := test.NewShutdown()

	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()

	sd.Shutdown(context.Background())
}
