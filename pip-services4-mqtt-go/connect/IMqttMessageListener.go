package connect

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IMqttMessageListener interface {
	OnMessage(message mqtt.Message)
}
