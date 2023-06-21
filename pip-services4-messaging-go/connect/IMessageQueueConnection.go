package connect

// IMessageQueueConnection Interface for queue connections
type IMessageQueueConnection interface {
	ReadQueueNames() ([]string, error)
	CreateQueue(name string) error
	DeleteQueue(name string) error
}
