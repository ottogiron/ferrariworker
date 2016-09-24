package elastic

import (
	"github.com/ottogiron/ferrariworker/backend"
	"github.com/ottogiron/ferrariworker/processor"
)

type elasticBackend struct {
}

func New() backend.Backend {
	return &elasticBackend{}
}

func (e *elasticBackend) Persist(workerID string, jobID string, jobResults []processor.JobResult) error {
	return nil
}

func (e *elasticBackend) JobResults(workerID string, jobID string) ([]processor.JobResult, error) {
	return nil, nil
}
