package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/stretchr/testify/assert"
)

func TestAnyValueGetAndSet(t *testing.T) {
	value := data.NewEmptyAnyValue()
	assert.Nil(t, value.GetAsObject())

	value.SetAsObject(1)
	assert.Equal(t, 1, value.GetAsInteger())
	assert.True(t, 1.0-value.GetAsFloat() < 0.001)
	assert.Equal(t, "1", value.GetAsString())
}

func TestAnyValueEquals(t *testing.T) {
	value := data.NewAnyValue(1)

	assert.True(t, value.Equals(1))
	assert.True(t, value.Equals(1.0))
	assert.True(t, value.Equals("1"))
}

func TestAnyValueGetAsNullable(t *testing.T) {
	value := data.NewAnyValue(1)

	_, ok := value.GetAsNullableInteger()
	assert.True(t, ok)
	_, ok = value.GetAsNullableLong()
	assert.True(t, ok)
	_, ok = value.GetAsNullableFloat()
	assert.True(t, ok)
	_, ok = value.GetAsNullableDouble()
	assert.True(t, ok)
	_, ok = value.GetAsNullableString()
	assert.True(t, ok)
	_, ok = value.GetAsNullableBoolean()
	assert.True(t, ok)
	_, ok = value.GetAsNullableDuration()
	assert.True(t, ok)
	_, ok = value.GetAsNullableDateTime()
	assert.True(t, ok)
}
