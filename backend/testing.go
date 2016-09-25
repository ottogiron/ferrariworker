package backend

import (
	"testing"
	"time"

	"github.com/ottogiron/ferrariworker/worker"
)

var persistCases = []struct {
	jobResult []*worker.JobResult
}{
	{[]*worker.JobResult{
		&worker.JobResult{
			ID:        "job_test",
			WorkerID:  "worker_test",
			Status:    worker.JobStatusSuccess,
			Output:    []byte{},
			StartTime: time.Now(),
			EndTime:   time.Now()},
	}},
}

func TestBackend(t *testing.T, backend Backend) {
	for _, tc := range persistCases {
		err := backend.Persist(tc.jobResult)
		if err != nil {
			t.Errorf("backend.Persist() => err:%s is not expected", err)
		}
	}
}
