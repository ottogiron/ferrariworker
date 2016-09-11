package processor

import "context"

//Adapter defines an messages source
type Adapter interface {
	Open() error
	Close() error
	Messages(context context.Context) (<-chan Message, error)
	ResultHandler(jobResult *JobResult, message Message) error
}

//Factory defines a actory from message adapters
type Factory interface {
	New(config AdapterConfig) Adapter
}
