package connect

import (
	kafka "github.com/Shopify/sarama"
)

type IKafkaMessageListener interface {
	// Setup is run at the beginning of a new session, before ConsumeClaim.
	Setup(kafka.ConsumerGroupSession) error

	// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
	// but before the offsets are committed for the very last time.
	Cleanup(kafka.ConsumerGroupSession) error

	// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
	// Once the Messages() channel is closed, the Handler must finish its processing
	// loop and exit.
	ConsumeClaim(kafka.ConsumerGroupSession, kafka.ConsumerGroupClaim) error

	// channel that recive signal that consummer is already start
	Ready() chan bool

	// set new channel for send ready signal
	SetReady(chFlag chan bool)
}
