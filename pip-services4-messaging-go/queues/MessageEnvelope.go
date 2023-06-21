package queues

import (
	"encoding/base64"
	"strings"
	"time"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// MessageEnvelope allows adding additional information to messages. A trace id, message id, and a message type
// are added to the data being sent/received. Additionally, a MessageEnvelope can reference a lock token.
// Side note: a MessageEnvelope"s message is stored as a buffer, so strings are converted
// using utf8 conversions.
type MessageEnvelope struct {
	reference        any
	TraceId          string    `json:"trace_id" bson:"trace_id"`         // The unique business transaction id that is used to trace calls across components.
	MessageId        string    `json:"message_id" bson:"message_id"`     // The message"s auto-generated ID.
	MessageType      string    `json:"message_type" bson:"message_type"` // String value that defines the stored message"s type.
	SentTime         time.Time `json:"sent_time" bson:"sent_time"`       // The time at which the message was sent.
	Message          []byte    `json:"message" bson:"message"`           // The stored message.
	JsonMapConvertor cconv.IJSONEngine[map[string]any]
}

// NewEmptyMessageEnvelope method are creates an empty MessageEnvelope
//
//	Returns: *MessageEnvelope new instance
func NewEmptyMessageEnvelope() *MessageEnvelope {
	c := MessageEnvelope{
		JsonMapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[map[string]any](),
	}
	return &c
}

// NewMessageEnvelope method are creates a new MessageEnvelope, which adds a trace id, message id, and a type to the
// data being sent/received.
//
//		Parameters:
//	  - traceId     (optional) transaction id to trace execution through call chain.
//	  - messageType       a string value that defines the message"s type.
//	  - message           the data being sent/received.
//		Returns: *MessageEnvelope new instance
func NewMessageEnvelope(traceId string, messageType string, message []byte) *MessageEnvelope {
	c := MessageEnvelope{
		JsonMapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[map[string]any](),
	}
	c.TraceId = traceId
	c.MessageType = messageType
	c.MessageId = cdata.IdGenerator.NextLong()
	c.Message = message
	return &c
}

// NewMessageEnvelopeFromObject method are creates a new MessageEnvelope, which adds a trace id, message id, and a type to the
// data object being sent/received.
//
//		Parameters:
//	  - traceId     (optional) transaction id to trace execution through call chain.
//	  - messageType       a string value that defines the message"s type.
//	  - message           the data object being sent/received.
//		Returns: *MessageEnvelope new instance
func NewMessageEnvelopeFromObject(traceId string, messageType string, message any) *MessageEnvelope {
	c := MessageEnvelope{
		JsonMapConvertor: cconv.NewDefaultCustomTypeJsonConvertor[map[string]any](),
	}
	c.TraceId = traceId
	c.MessageType = messageType
	c.MessageId = cdata.IdGenerator.NextLong()
	c.SetMessageAsObject(message)
	return &c
}

// GetReference method are returns the lock token that this MessageEnvelope references.
func (c *MessageEnvelope) GetReference() any {
	return c.reference
}

// SetReference method are sets a lock token reference for this MessageEnvelope.
//
//	Parameters:
//		- value     the lock token to reference.
func (c *MessageEnvelope) SetReference(value any) {
	c.reference = value
}

// GetMessageAsString method are returns the information stored in this message as a string.
func (c *MessageEnvelope) GetMessageAsString() string {
	return string(c.Message)
}

// SetMessageAsString method are stores the given string.
//
//	Parameters:
//		- value    the string to set. Will be converted to a bufferg.
func (c *MessageEnvelope) SetMessageAsString(value string) {
	c.Message = []byte(value)
}

// GetMessageAs method are returns the value that was stored in this message as object.
//
//	see  SetMessageAsObject
func GetMessageAs[T any](envelope *MessageEnvelope) (T, error) {
	var defaultValue T
	if envelope.Message == nil {
		return defaultValue, nil
	}
	return cconv.NewDefaultCustomTypeJsonConvertor[T]().FromJson(string(envelope.Message))
}

// SetMessageAsObject method are stores the given value as a JSON string.
//
//	Parameters:
//		- value     the value to convert to JSON and store in this message.
//	see  GetMessageAs
func (c *MessageEnvelope) SetMessageAsObject(value any) {
	if value == nil {
		c.Message = []byte{}
	} else {
		if msg, err := cconv.JsonConverter.ToJson(value); err == nil {
			c.Message = []byte(msg)
		}
	}
}

// String method are convert"s this MessageEnvelope to a string, using the following format:
// <trace_id>,<MessageType>,<message.toString>
// If any of the values are nil, they will be replaced with ---.
//
//	Returns: the generated string.
func (c *MessageEnvelope) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")
	if c.TraceId == "" {
		builder.WriteString("---")
	} else {
		builder.WriteString(c.TraceId)
	}
	builder.WriteString(",")
	if c.MessageType == "" {
		builder.WriteString("---")
	} else {
		builder.WriteString(c.MessageType)
	}
	builder.WriteString(",")
	if c.Message == nil {
		builder.WriteString("---")
	} else {
		builder.Write(c.Message)
	}
	builder.WriteString("]")
	return builder.String()
}

func (c MessageEnvelope) MarshalJSON() ([]byte, error) {
	jsonData := map[string]any{
		"message_id":   c.MessageId,
		"trace_id":     c.TraceId,
		"message_type": c.MessageType,
	}

	if !c.SentTime.IsZero() {
		jsonData["sent_time"] = c.SentTime
	} else {
		jsonData["sent_time"] = time.Now()
	}

	if c.Message != nil {
		base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(c.Message)))
		base64.StdEncoding.Encode(base64Text, c.Message)
		jsonData["message"] = string(base64Text)
	}

	msg, err := cconv.JsonConverter.ToJson(jsonData)
	if err != nil {
		return nil, err
	}
	return []byte(msg), nil
}

func (c *MessageEnvelope) UnmarshalJSON(data []byte) error {
	var jsonData map[string]any
	jsonData, err := c.JsonMapConvertor.FromJson(string(data))
	if err != nil {
		return err
	}

	if _val, ok := jsonData["message_id"].(string); ok {
		c.MessageId = _val
	}
	if _val, ok := jsonData["trace_id"].(string); ok {
		c.TraceId = _val
	}
	if _val, ok := jsonData["message_type"].(string); ok {
		c.MessageType = _val
	}
	c.SentTime = cconv.DateTimeConverter.ToDateTime(jsonData["sent_time"])

	if base64Text, ok := jsonData["message"].(string); ok && base64Text != "" {
		data := make([]byte, base64.StdEncoding.DecodedLen(len(base64Text)))
		n, err := base64.StdEncoding.Decode(data, []byte(base64Text))
		if err != nil {
			return err
		}
		c.Message = data[:n]
	}

	return nil
}
