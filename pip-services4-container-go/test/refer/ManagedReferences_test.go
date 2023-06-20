package test_refer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-container-go/refer"
	obuild "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/build"
)

func TestAutoCreateComponent(t *testing.T) {
	refs := crefer.NewEmptyManagedReferences()

	factory := obuild.NewDefaultObservabilityFactory()
	refs.Put(context.Background(), nil, factory)

	logger, err := refs.GetOneRequired(
		refer.NewDescriptor("*", "logger", "*", "*", "*"),
	)

	assert.Nil(t, err)
	assert.NotNil(t, logger)
}

func TestStringLocator(t *testing.T) {
	refs := crefer.NewEmptyManagedReferences()

	factory := obuild.NewDefaultObservabilityFactory()
	refs.Put(context.Background(), nil, factory)

	logger := refs.GetOneOptional("ABC")

	assert.Nil(t, logger)
}

func TestNilLocator(t *testing.T) {
	refs := crefer.NewEmptyManagedReferences()

	factory := obuild.NewDefaultObservabilityFactory()
	refs.Put(context.Background(), nil, factory)

	logger := refs.GetOneOptional(nil)

	assert.Nil(t, logger)
}
