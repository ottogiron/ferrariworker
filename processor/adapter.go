package processor

import (
	"context"

	"github.com/ottogiron/ferrariworker/config"
)

//Adapter defines an messages source
type Adapter interface {
	Open() error
	Close() error
	Messages(context context.Context) (<-chan Message, error)
	ResultHandler(jobResult *JobResult, message Message) error
}

//Factory defines a actory from message adapters
type Factory interface {
	New(config config.AdapterConfig) Adapter
}
