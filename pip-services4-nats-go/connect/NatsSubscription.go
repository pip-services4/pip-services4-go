package connect

import "github.com/nats-io/nats.go"

type NatsSubscription struct {
	Subject    string
	QueueGroup string
	Listener   INatsMessageListener
	Handler    *nats.Subscription
}
