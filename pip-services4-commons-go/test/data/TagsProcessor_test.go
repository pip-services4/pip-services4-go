package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTags(t *testing.T) {
	tag := data.TagsProcessor.NormalizeTag("  A_b#c ")
	assert.Equal(t, "A b c", tag)

	tags := data.TagsProcessor.NormalizeTags([]string{"  A_b#c ", "d__E f"})
	assert.Len(t, tags, 2)
	assert.Equal(t, "A b c", tags[0])
	assert.Equal(t, "d E f", tags[1])

	tags = data.TagsProcessor.NormalizeTagList("  A_b#c ,d__E f;;")
	assert.Len(t, tags, 3)
	assert.Equal(t, "A b c", tags[0])
	assert.Equal(t, "d E f", tags[1])
}

func TestCompressTags(t *testing.T) {
	tag := data.TagsProcessor.CompressTag("  A_b#c ")
	assert.Equal(t, "abc", tag)

	tags := data.TagsProcessor.CompressTags([]string{"  A_b#c ", "d__E f"})
	assert.Len(t, tags, 2)
	assert.Equal(t, "abc", tags[0])
	assert.Equal(t, "def", tags[1])

	tags = data.TagsProcessor.CompressTagList("  A_b#c ,d__E f;;")
	assert.Len(t, tags, 3)
	assert.Equal(t, "abc", tags[0])
	assert.Equal(t, "def", tags[1])
}

func TestExtractHashTags(t *testing.T) {
	tags := data.TagsProcessor.ExtractHashTags("  #Tag_1  #TAG2#tag3 ")
	assert.Len(t, tags, 3)
	assert.Equal(t, "tag1", tags[0])
	assert.Equal(t, "tag2", tags[1])
	assert.Equal(t, "tag3", tags[2])
}
