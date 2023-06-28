package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/process"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTags(t *testing.T) {
	tag := process.TagsProcessor.NormalizeTag("  A_b#c ")
	assert.Equal(t, "A b c", tag)

	tags := process.TagsProcessor.NormalizeTags([]string{"  A_b#c ", "d__E f"})
	assert.Len(t, tags, 2)
	assert.Equal(t, "A b c", tags[0])
	assert.Equal(t, "d E f", tags[1])

	tags = process.TagsProcessor.NormalizeTagList("  A_b#c ,d__E f;;")
	assert.Len(t, tags, 3)
	assert.Equal(t, "A b c", tags[0])
	assert.Equal(t, "d E f", tags[1])
}

func TestCompressTags(t *testing.T) {
	tag := process.TagsProcessor.CompressTag("  A_b#c ")
	assert.Equal(t, "abc", tag)

	tags := process.TagsProcessor.CompressTags([]string{"  A_b#c ", "d__E f"})
	assert.Len(t, tags, 2)
	assert.Equal(t, "abc", tags[0])
	assert.Equal(t, "def", tags[1])

	tags = process.TagsProcessor.CompressTagList("  A_b#c ,d__E f;;")
	assert.Len(t, tags, 3)
	assert.Equal(t, "abc", tags[0])
	assert.Equal(t, "def", tags[1])
}

func TestExtractHashTags(t *testing.T) {
	tags := process.TagsProcessor.ExtractHashTags("  #Tag_1  #TAG2#tag3 ")
	assert.Len(t, tags, 3)
	assert.Equal(t, "tag1", tags[0])
	assert.Equal(t, "tag2", tags[1])
	assert.Equal(t, "tag3", tags[2])
}
