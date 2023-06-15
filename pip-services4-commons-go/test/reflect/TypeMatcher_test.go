package test_reflect

import (
	refl "reflect"
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestMatchInteger(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("int", 123))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("Integer", 123))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(int(1)), 123))
}

func TestMatchLong(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("long", int64(123)))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("Long", int64(123)))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(int64(1)), int64(123)))

	assert.True(t, reflect.TypeMatcher.MatchType(convert.Long, refl.TypeOf(int8(1))))
	assert.True(t, reflect.TypeMatcher.MatchType(convert.Long, refl.TypeOf(int16(1))))
	assert.True(t, reflect.TypeMatcher.MatchType(convert.Long, refl.TypeOf(int32(1))))
	assert.True(t, reflect.TypeMatcher.MatchType(convert.Long, refl.TypeOf(int64(1))))
	assert.True(t, reflect.TypeMatcher.MatchType(convert.Long, refl.TypeOf(float32(1))))
	assert.True(t, reflect.TypeMatcher.MatchType(convert.Long, refl.TypeOf(float64(1))))
}

func TestMatchBoolean(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("bool", true))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("Boolean", true))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(false), true))
}

func TestMatchFloat(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("float", float32(123.456)))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("Float", float32(123.456)))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(float32(0.1)), float32(123.456)))
}

func TestMatchDouble(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("double", float64(123.456)))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("Double", float64(123.456)))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(float64(0)), float64(123.456)))
}

func TestMatchString(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("string", "ABC"))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(""), "ABC"))
}

func TestMatchDateTime(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("date", time.Now()))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("DateTime", &time.Time{}))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(&time.Time{}), time.Time{}))
}

func TestMatchDuration(t *testing.T) {
	assert.True(t, reflect.TypeMatcher.MatchValueByName("duration", 123))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("TimeSpan", int64(123)))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("TimeSpan", time.Duration(1)))
}

func TestMatchMap(t *testing.T) {
	dict := map[string]any{}
	assert.True(t, reflect.TypeMatcher.MatchValueByName("map", dict))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("dict", dict))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("Dictionary", dict))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf(map[string]any{}), dict))
}

func TestMatchArray(t *testing.T) {
	list := []any{}
	assert.True(t, reflect.TypeMatcher.MatchValueByName("list", list))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("array", list))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("object[]", list))
	assert.True(t, reflect.TypeMatcher.MatchValue(refl.TypeOf([]any{}), list))

	array := [...]int{}
	assert.True(t, reflect.TypeMatcher.MatchValueByName("list", array))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("array", array))
	assert.True(t, reflect.TypeMatcher.MatchValueByName("int[]", array))
}
