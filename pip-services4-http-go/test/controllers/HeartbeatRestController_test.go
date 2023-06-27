package test_controllers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatRestController(t *testing.T) {

	url := fmt.Sprintf("http://localhost:%d", HeartbeatRestControllerPort)
	// Test "Heartbeat"
	getRes, getErr := http.Get(url + "/heartbeat")
	assert.Nil(t, getErr)
	assert.NotNil(t, getRes)
}
