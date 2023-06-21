package queues

// MessagingCapabilities data object that contains supported capabilities of a message queue.
// If certain capability is not supported a queue will throw NotImplemented exception.
type MessagingCapabilities struct {
	canMessageCount bool
	canSend         bool
	canReceive      bool
	canPeek         bool
	canPeekBatch    bool
	canRenewLock    bool
	canAbandon      bool
	canDeadLetter   bool
	canClear        bool
}

// NewMessagingCapabilities method are creates a new instance of the capabilities object.
//	Parameters:
//		- canMessageCount   true if queue supports reading message count.
//		- canSend           true if queue is able to send messages.
//		- canReceive        true if queue is able to receive messages.
//		- canPeek           true if queue is able to peek messages.
//		- canPeekBatch      true if queue is able to peek multiple messages in one batch.
//		- canRenewLock      true if queue is able to renew message lock.
//		- canAbandon        true if queue is able to abandon messages.
//		- canDeadLetter     true if queue is able to send messages to dead letter queue.
//		- canClear          true if queue can be cleared.
//	Returns: *MessagingCapabilities
func NewMessagingCapabilities(canMessageCount bool, canSend bool, canReceive bool,
	canPeek bool, canPeekBatch bool, canRenewLock bool, canAbandon bool,
	canDeadLetter bool, canClear bool) *MessagingCapabilities {

	c := MessagingCapabilities{}
	c.canMessageCount = canMessageCount
	c.canSend = canSend
	c.canReceive = canReceive
	c.canPeek = canPeek
	c.canPeekBatch = canPeekBatch
	c.canRenewLock = canRenewLock
	c.canAbandon = canAbandon
	c.canDeadLetter = canDeadLetter
	c.canClear = canClear
	return &c
}

// CanMessageCount method are informs if the queue is able to read number of messages.
//	Returns: true if queue supports reading message count.
func (c *MessagingCapabilities) CanMessageCount() bool {
	return c.canMessageCount
}

// CanSend method are informs if the queue is able to send messages.
//	Returns: true if queue is able to send messages.
func (c *MessagingCapabilities) CanSend() bool {
	return c.canSend
}

// CanReceive method are informs if the queue is able to receive messages.
//	Returns: true if queue is able to receive messages.
func (c *MessagingCapabilities) CanReceive() bool {
	return c.canReceive
}

// CanPeek method are informs if the queue is able to peek messages.
//	Returns: true if queue is able to peek messages.
func (c *MessagingCapabilities) CanPeek() bool {
	return c.canPeek
}

// CanPeekBatch method are informs if the queue is able to peek multiple messages in one batch.
//	Returns: true if queue is able to peek multiple messages in one batch.
func (c *MessagingCapabilities) CanPeekBatch() bool {
	return c.canPeekBatch
}

// CanRenewLock method are informs if the queue is able to renew message lock.
//	Returns: true if queue is able to renew message lock.
func (c *MessagingCapabilities) CanRenewLock() bool {
	return c.canRenewLock
}

// CanAbandon method are informs if the queue is able to abandon messages.
//	Returns: true if queue is able to abandon.
func (c *MessagingCapabilities) CanAbandon() bool {
	return c.canAbandon
}

// CanDeadLetter method are informs if the queue is able to send messages to dead letter queue.
//	Returns: true if queue is able to send messages to dead letter queue.
func (c *MessagingCapabilities) CanDeadLetter() bool {
	return c.canDeadLetter
}

// CanClear method are informs if the queue can be cleared.
//	Returns: true if queue can be cleared.
func (c *MessagingCapabilities) CanClear() bool {
	return c.canClear
}
