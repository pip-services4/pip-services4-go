package test_validate

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	"github.com/stretchr/testify/assert"
)

func TestObjectCompare(t *testing.T) {
	assert.True(t, validate.ObjectComparator.Compare(1, "eq", 1))
	assert.True(t, validate.ObjectComparator.Compare("ABC", "eq", "ABC"))
	assert.False(t, validate.ObjectComparator.Compare(1, "eq", 2))
	assert.False(t, validate.ObjectComparator.Compare("ABC", "eq", "XYZ"))

	assert.False(t, validate.ObjectComparator.Compare(1, "ne", 1))
	assert.True(t, validate.ObjectComparator.Compare("ABC", "ne", "XYZ"))
	assert.True(t, validate.ObjectComparator.Compare(2, "ne", 1))
	assert.True(t, validate.ObjectComparator.Compare(1, "ne", 2))

	assert.False(t, validate.ObjectComparator.Compare(1, "gt", 1))
	assert.False(t, validate.ObjectComparator.Compare("ABC", "gte", "XYZ"))
	assert.True(t, validate.ObjectComparator.Compare(2, "gt", 1))
	assert.False(t, validate.ObjectComparator.Compare(1, "gt", 2))

	assert.True(t, validate.ObjectComparator.Compare(1, "gte", 1))
	assert.False(t, validate.ObjectComparator.Compare("ABC", "gte", "XYZ"))
	assert.True(t, validate.ObjectComparator.Compare(2, "gte", 1))
	assert.False(t, validate.ObjectComparator.Compare(1, "gte", 2))

	assert.False(t, validate.ObjectComparator.Compare(1, "lt", 1))
	assert.False(t, validate.ObjectComparator.Compare("ABC", "lt", "XYZ"))
	assert.False(t, validate.ObjectComparator.Compare(2, "lt", 1))
	assert.True(t, validate.ObjectComparator.Compare(1, "lt", 2))

	assert.True(t, validate.ObjectComparator.Compare(1, "lte", 1))
	assert.False(t, validate.ObjectComparator.Compare("ABC", "lte", "XYZ"))
	assert.False(t, validate.ObjectComparator.Compare(2, "lte", 1))
	assert.True(t, validate.ObjectComparator.Compare(1, "lte", 2))

	assert.False(t, validate.ObjectComparator.Compare(1, "match", 1))
	assert.False(t, validate.ObjectComparator.Compare("ABC", "match", "XYZ"))
	assert.True(t, validate.ObjectComparator.Compare("ABC", "like", "A.*C"))
}
