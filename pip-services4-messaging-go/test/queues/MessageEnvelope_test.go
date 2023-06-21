package test_queues

import (
	"encoding/json"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	"github.com/stretchr/testify/assert"
)

type testType struct {
	Value string `json:"value"`
}

type messageEnvelopeTest struct{}

func NewMessageEnvelopTest() *messageEnvelopeTest {
	c := messageEnvelopeTest{}
	return &c
}

func (c *messageEnvelopeTest) TestSerializeMessage(t *testing.T) {
	message := queues.NewMessageEnvelope("123", "TestMessage", []byte("This is a test message"))
	assert.Equal(t, "123", message.TraceId)
	assert.Equal(t, "TestMessage", message.MessageType)
	assert.Equal(t, []byte("This is a test message"), message.Message)
	assert.NotEqual(t, "", message.MessageId)

	buffer, err := json.Marshal(message)
	assert.Nil(t, err)
	assert.True(t, len(buffer) > 0)

	message2 := queues.NewEmptyMessageEnvelope()
	err = json.Unmarshal(buffer, message2)
	assert.Nil(t, err)
	assert.Equal(t, message.MessageId, message2.MessageId)
	assert.Equal(t, message.TraceId, message2.TraceId)
	assert.Equal(t, message.MessageType, message2.MessageType)
	assert.Equal(t, message.Message, message2.Message)
}

func (c *messageEnvelopeTest) TestMessageEnvelopMethods(t *testing.T) {
	message := queues.NewMessageEnvelope("123", "TestMessage", []byte("This is a test message"))
	assert.Equal(t, "123", message.TraceId)
	assert.Equal(t, "TestMessage", message.MessageType)
	assert.Equal(t, []byte("This is a test message"), message.Message)
	assert.NotEqual(t, "", message.MessageId)

	assert.Equal(t, message.String(), "[123,TestMessage,This is a test message]")

	testObj := testType{Value: "This is a test message"}
	message.SetMessageAsObject(testObj)

	resultObj, err := queues.GetMessageAs[testType](message)
	assert.Nil(t, err)
	assert.Equal(t, testObj, resultObj)
}

func TestMessageEnvelop(t *testing.T) {
	test := NewMessageEnvelopTest()

	t.Run("MessageEnvelop:Serialize Message", test.TestSerializeMessage)
	t.Run("MessageEnvelop:Methods", test.TestMessageEnvelopMethods)
}
