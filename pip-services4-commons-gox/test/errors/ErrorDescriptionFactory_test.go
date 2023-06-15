package test_errors

import (
	"errors"
	"testing"

	cerrors "github.com/pip-services4/pip-services4-commons-go/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateFromApplicationError(t *testing.T) {
	err := cerrors.NewInternalError("123", "CODE", "Error message")

	d := cerrors.ErrorDescriptionFactory.Create(err)

	assert.Equal(t, cerrors.Internal, d.Category)
	assert.Equal(t, "123", d.CorrelationId)
	assert.Equal(t, "CODE", d.Code)
	assert.Equal(t, "Error message", d.Message)
	assert.Equal(t, 500, d.Status)
}

func TestCreateFromError(t *testing.T) {
	err := errors.New("Message")

	d := cerrors.ErrorDescriptionFactory.Create(err)

	assert.Equal(t, cerrors.Unknown, d.Category)
	assert.Equal(t, "", d.CorrelationId)
	assert.Equal(t, "UNKNOWN", d.Code)
	assert.Equal(t, "Message", d.Message)
	assert.Equal(t, 500, d.Status)
}

func TestCreateFromString(t *testing.T) {
	d := cerrors.ErrorDescriptionFactory.Create("Message")

	assert.Equal(t, cerrors.Unknown, d.Category)
	assert.Equal(t, "", d.CorrelationId)
	assert.Equal(t, "UNKNOWN", d.Code)
	assert.Equal(t, "Message", d.Message)
	assert.Equal(t, 500, d.Status)
}
