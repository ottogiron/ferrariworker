package adapter

import (
	"context"

	"github.com/ottogiron/ferraritrunk/config"
	"github.com/ottogiron/ferraritrunk/worker"
)

//Adapter defines an messages source
type Adapter interface {
	Open() error
	Close() error
	Messages(context context.Context) (<-chan worker.Message, error)
	ResultHandler(jobResult *worker.JobResult, message worker.Message) error
}

//Factory defines a actory from message adapters
type Factory func(config config.Config) Adapter
