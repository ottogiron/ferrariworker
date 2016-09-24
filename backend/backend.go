package backend

import (
	"github.com/ottogiron/ferrariworker/config"
	"github.com/ottogiron/ferrariworker/worker"
)

//Backend defines a data store for jobs to be persisted
type Backend interface {
	Persist(workerId string, jobID string, jobResults []worker.JobResult) error
	JobResults(workerId string, jobID string) ([]worker.JobResult, error)
}

//Factory defines a actory from message adapters
type Factory interface {
	New(config config.AdapterConfig) Backend
}
