package queues

// CallbackMessageReceiver allows to wrap message callback into IMessageReceiver
type CallbackMessageReceiver struct {
	Callback func(message *MessageEnvelope, queue IMessageQueue) error
}

func NewCallbackMessageReceiver(callback func(message *MessageEnvelope, queue IMessageQueue) error) *CallbackMessageReceiver {
	c := CallbackMessageReceiver{
		Callback: callback,
	}
	return &c
}

func (c *CallbackMessageReceiver) ReceiveMessage(message *MessageEnvelope, queue IMessageQueue) (err error) {
	return c.Callback(message, queue)
}
