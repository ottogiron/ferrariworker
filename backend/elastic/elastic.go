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

func (e *elasticBackend) Persist(jobResults []processor.JobResult) error {
	return nil
}
