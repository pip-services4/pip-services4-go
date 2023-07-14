package connect

type MqttSubscription struct {
	Topic    string
	Qos      byte
	Listener IMqttMessageListener
	Skip     int32
}
