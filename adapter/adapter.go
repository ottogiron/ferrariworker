package adapter

import (
	"context"

	"github.com/ottogiron/ferrariworker/config"
	"github.com/ottogiron/ferrariworker/worker"
)

//Adapter defines an messages source
type Adapter interface {
	Open() error
	Close() error
	Messages(context context.Context) (<-chan worker.Message, error)
	ResultHandler(jobResult *worker.JobResult, message worker.Message) error
}

//Factory defines a actory from message adapters
type Factory interface {
	New(config config.AdapterConfig) Adapter
}
