package test_query

import (
	"encoding/json"
	"testing"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"

	"github.com/stretchr/testify/assert"
)

func TestFilterParamsCreate(t *testing.T) {
	filter := query.NewFilterParamsFromTuples(
		"value1", 123,
		"value2", "ABC",
	)

	assert.Equal(t, 2, filter.Len())
}

func TestFilterParamsJsonSerialization(t *testing.T) {
	json1 := []byte("{\"key1\":\"1\",\"key2\":\"A\"}")

	var value *cdata.StringValueMap
	err := json.Unmarshal(json1, &value)
	assert.Empty(t, err)
	val, ok := value.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "1", val)
	val, ok = value.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "A", val)

	json2, err2 := json.Marshal(value)
	assert.Empty(t, err2)
	assert.Equal(t, json1, json2)
}
