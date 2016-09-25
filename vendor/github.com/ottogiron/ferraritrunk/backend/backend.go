package backend

import (
	"github.com/ottogiron/ferraritrunk/config"
	"github.com/ottogiron/ferraritrunk/worker"
)

//Backend defines a data store for jobs to be persisted
type Backend interface {
	Persist(jobResults []*worker.JobResult) error
	JobResults(workerId string) ([]*worker.JobResult, error)
	Job(jobID string) (*worker.JobResult, error)
}

//Factory defines a actory from message adapters
type Factory func(config config.Config) (Backend, error)
