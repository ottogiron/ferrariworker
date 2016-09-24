package backend

import (
	"testing"
	"time"

	"github.com/ottogiron/ferrariworker/worker"
)

var persistCases = []struct {
	jobResult []worker.JobResult
}{
	{[]worker.JobResult{
		worker.JobResult{
			Status:    worker.JobStatusSuccess,
			Output:    []byte{},
			StartTime: time.Now(),
			EndTime:   time.Now()},
	}},
}

func TestBackend(t *testing.T, backend Backend) {
	for _, tc := range persistCases {
		err := backend.Persist("worker_test", "job_test", tc.jobResult)
		if err != nil {
			t.Errorf("backend.Persist() => err:%s is not expected", err)
		}
	}
}
