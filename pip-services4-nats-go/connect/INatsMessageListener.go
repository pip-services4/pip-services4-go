package connect

import (
	"github.com/nats-io/nats.go"
)

type INatsMessageListener interface {
	OnMessage(message *nats.Msg)
}
