package test_controllers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusRestController(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", StatusRestControllerPort)
	// Test "Status"
	getRes, getErr := http.Get(url + "/status")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
}
