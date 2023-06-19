package test_refer

import (
	"context"
	"testing"

	conf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/stretchr/testify/assert"
)

func TestResolveDependencies(t *testing.T) {
	ref1 := "AAA"
	ref2 := "BBB"
	refs := refer.NewReferencesFromTuples(
		context.Background(),
		"Reference1", ref1,
		refer.NewDescriptor("pip-services-commons", "reference", "object", "ref2", "1.0"), ref2,
	)

	resolver := refer.NewDependencyResolverFromTuples(
		context.Background(),
		"ref1", "Reference1",
		"ref2", refer.NewDescriptor("pip-services-commons", "reference", "*", "*", "*"),
	)
	resolver.SetReferences(context.Background(), refs)

	ref, _ := resolver.GetOneRequired("ref1")
	assert.Equal(t, ref1, ref)
	ref, _ = resolver.GetOneRequired("ref2")
	assert.Equal(t, ref2, ref)
	assert.Nil(t, resolver.GetOneOptional("ref3"))
}

func TestConfigureDependencies(t *testing.T) {
	ref1 := "AAA"
	ref2 := "BBB"
	refs := refer.NewReferencesFromTuples(
		context.Background(),
		"Reference1", ref1,
		refer.NewDescriptor("pip-services-commons", "reference", "object", "ref2", "1.0"), ref2,
	)

	config := conf.NewConfigParamsFromTuples(
		"dependencies.ref1", "Reference1",
		"dependencies.ref2", "pip-services-commons:reference:*:*:*",
		"dependencies.ref3", "",
	)

	resolver := refer.NewDependencyResolverWithParams(context.Background(), config, nil)
	resolver.SetReferences(context.Background(), refs)

	ref, _ := resolver.GetOneRequired("ref1")
	assert.Equal(t, ref1, ref)
	ref, _ = resolver.GetOneRequired("ref2")
	assert.Equal(t, ref2, ref)
	assert.Nil(t, resolver.GetOneOptional("ref3"))
}
