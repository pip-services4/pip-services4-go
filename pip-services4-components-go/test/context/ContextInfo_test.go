package test_info

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/stretchr/testify/assert"
)

func TestContextInfo(t *testing.T) {
	contextInfo := cctx.NewContextInfo()

	assert.Equal(t, cctx.ContextInfoNameUnknown, contextInfo.Name)
	assert.Equal(t, "", contextInfo.Description)
	assert.True(t, len(contextInfo.ContextId) > 0)

	contextInfo.Name = "new name"
	contextInfo.Description = "new description"
	contextInfo.ContextId = "new context id"

	assert.Equal(t, "new name", contextInfo.Name)
	assert.Equal(t, "new description", contextInfo.Description)
	assert.Equal(t, "new context id", contextInfo.ContextId)
}

func TestContextInfoFromConfig(t *testing.T) {
	cfg := config.NewConfigParamsFromTuples(
		cctx.ContextInfoParameterInfoName, "new name",
		cctx.ContextInfoParameterInfoDescription, "new description",
		cctx.ContextInfoSectionNameProperties+".access_key", "key",
		cctx.ContextInfoSectionNameProperties+".store_key", "store key",
	)

	contextInfo := cctx.NewContextInfoFromConfig(context.Background(), cfg)
	assert.Equal(t, "new name", contextInfo.Name)
	assert.Equal(t, "new description", contextInfo.Description)
}
