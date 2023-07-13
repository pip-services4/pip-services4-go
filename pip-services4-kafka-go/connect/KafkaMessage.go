package connect

import (
	kafka "github.com/Shopify/sarama"
)

// Kafka message structure
type KafkaMessage struct {
	// Kafka consummer message
	Message *kafka.ConsumerMessage
	// Counsummer session
	Session kafka.ConsumerGroupSession
}
