package test_query

import (
	"encoding/json"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type user struct {
	Name string
	Age  int
}

func TestNewEmptyDataPage(t *testing.T) {
	dataPage := query.NewEmptyDataPage[user]()

	assert.False(t, dataPage.HasData())
	assert.Nil(t, dataPage.Data)

	assert.False(t, dataPage.HasTotal())
	assert.Equal(t, query.EmptyTotalValue, dataPage.Total)
}

func TestNewDataPage(t *testing.T) {
	arr := []user{{
		Name: "User1",
		Age:  26,
	}, {
		Name: "User2",
		Age:  45,
	}}
	dataPage := query.NewDataPage[user](arr, query.EmptyTotalValue)

	assert.True(t, dataPage.HasData())
	assert.Equal(t, 2, len(dataPage.Data))

	assert.False(t, dataPage.HasTotal())
	assert.Equal(t, query.EmptyTotalValue, dataPage.Total)

	// Test with total marshaling
	dataPage.Total = 2
	buf, err := json.Marshal(dataPage)
	assert.Nil(t, err)
	assert.True(t, len(buf) > 0)

	var resultedDataPage query.DataPage[user]
	err = json.Unmarshal(buf, &resultedDataPage)
	assert.Nil(t, err)

	assert.True(t, resultedDataPage.HasData())
	assert.True(t, resultedDataPage.HasTotal())
	assert.Equal(t, 2, len(resultedDataPage.Data))
	assert.Equal(t, 2, resultedDataPage.Total)

	// Test with total marshaling
	dataPage.Total = query.EmptyTotalValue
	buf, err = json.Marshal(dataPage)
	assert.Nil(t, err)
	assert.True(t, len(buf) > 0)

	var resultedDataPageWithoutTotal query.DataPage[user]
	err = json.Unmarshal(buf, &resultedDataPageWithoutTotal)
	assert.Nil(t, err)

	assert.True(t, resultedDataPageWithoutTotal.HasData())
	assert.False(t, resultedDataPageWithoutTotal.HasTotal())
	assert.Equal(t, 2, len(resultedDataPageWithoutTotal.Data))
	assert.True(t, resultedDataPageWithoutTotal.Total >= 0)
}
