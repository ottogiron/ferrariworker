package elastic

import (
	"github.com/ottogiron/ferrariworker/backend"
	"github.com/ottogiron/ferrariworker/worker"
	"gopkg.in/olivere/elastic.v2"
)

type elasticBackend struct {
	client *elastic.Client
}

func New(client *elastic.Client) backend.Backend {
	return &elasticBackend{client}
}

func (e *elasticBackend) Persist(workerID string, jobID string, jobResults []worker.JobResult) error {
	return nil
}

func (e *elasticBackend) JobResults(workerID string, jobID string) ([]worker.JobResult, error) {
	return nil, nil
}
