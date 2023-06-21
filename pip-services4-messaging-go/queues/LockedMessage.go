package queues

import "time"

// LockedMessage data object used to store and lock incoming messages in MemoryMessageQueue.
//	see: MemoryMessageQueue
type LockedMessage struct {
	Message        *MessageEnvelope `json:"message" bson:"message"`                 // The incoming message.
	ExpirationTime time.Time        `json:"expiration_time" bson:"expiration_time"` // The expiration time for the message lock. If it is null then the message is not locked.
	Timeout        time.Duration    `json:"timeout" bson:"timeout"`                 // The lock timeout in milliseconds.
}
