package clients

import "time"

type DataDogLogMessage struct {
	Time         time.Time         `json:"time"`
	Tags         map[string]string `json:"tags"`
	Status       string            `json:"status"`
	Source       string            `json:"source"`
	Service      string            `json:"service"`
	Host         string            `json:"host"`
	Message      string            `json:"message"`
	LoggerName   string            `json:"logger_name"`
	ThreadName   string            `json:"thread_name"`
	ErrorMessage string            `json:"error_message"`
	ErrorKind    string            `json:"error_kind"`
	ErrorStack   string            `json:"error_stack"`
}
