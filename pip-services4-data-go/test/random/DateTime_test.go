package test_random

import (
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/random"
	"github.com/stretchr/testify/assert"
)

func TestNextDate(t *testing.T) {
	date := random.DateTime.NextDate(time.Date(2015, 0, 1, 0, 0, 0, 0, time.UTC), time.Date(2016, 0, 1, 0, 0, 0, 0, time.UTC))
	assert.True(t, date.Year() == 2015 || date.Year() == 2016)
}

func TestNextDateTime(t *testing.T) {
	startDate := time.Date(2015, 0, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2016, 0, 1, 0, 0, 0, 0, time.UTC)
	date := random.DateTime.NextDateTime(startDate, endDate)
	assert.True(t, date.Unix() >= startDate.Unix() && date.Unix() <= endDate.Unix())
}

func TestUpdateTime(t *testing.T) {
	oldDate := time.Date(2016, 10, 10, 0, 0, 0, 0, time.UTC)

	date := random.DateTime.UpdateDateTime(oldDate, 0)
	assert.True(t, date.Unix() >= oldDate.Unix()-10*24*3600 ||
		date.Unix() >= oldDate.Unix()+10*24*3600)

	date = random.DateTime.UpdateDateTime(oldDate, 3*24*3600)
	assert.True(t, date.Unix() >= oldDate.Unix()-3*24*3600 ||
		date.Unix() >= oldDate.Unix()+3*24*3600)

	date = random.DateTime.UpdateDateTime(oldDate, -3)
	assert.True(t, date.Unix() == oldDate.Unix())
}
