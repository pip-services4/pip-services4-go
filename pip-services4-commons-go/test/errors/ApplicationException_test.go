package test_errors

import (
	"errors"
	"testing"

	cerrors "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/stretchr/testify/assert"
)

func TestDefaultError(t *testing.T) {
	err := cerrors.NewError("")

	assert.Equal(t, "Unknown error", err.Message)
	assert.Equal(t, "UNKNOWN", err.Code)
	assert.Equal(t, 500, err.Status)
}

func TestWithCause(t *testing.T) {
	cause := errors.New("Cause error")
	err := cerrors.NewError("").WithCause(cause)

	assert.Equal(t, cause.Error(), err.Cause)
}

func TestWithTraceId(t *testing.T) {
	err := cerrors.NewError("").WithTraceId("123")

	assert.Equal(t, "123", err.TraceId)
}

func TestWithStatus(t *testing.T) {
	err := cerrors.NewError("").WithStatus(300)

	assert.Equal(t, 300, err.Status)
}
