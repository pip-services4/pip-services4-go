package connect

import (
	kafka "github.com/Shopify/sarama"
)

// Subscription structure
type KafkaSubscription struct {
	Topic    string
	GroupId  string
	Listener IKafkaMessageListener
	Handler  *kafka.ConsumerGroup
}
