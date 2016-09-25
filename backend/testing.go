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
			WorkerID:  "test_worker_1",
			Status:    worker.JobStatusSuccess,
			Output:    []byte{},
			StartTime: time.Now(),
			EndTime:   time.Now()},
	}},
}

func TestBackend(t *testing.T, backend Backend) {

	for _, tc := range persistCases {

		t.Run("Persist", func(t *testing.T) {
			err := backend.Persist(tc.jobResult)
			if err != nil {
				t.Errorf("backend.Persist() => err:%s is not expected", err)
			}
		})

		t.Run("JobResults", func(t *testing.T) {
			workerID := tc.jobResult[0].WorkerID
			persistedJobs, err := backend.JobResults(workerID)

			if err != nil {
				t.Fatalf("backend.JobResults(%s) err: %s was not expected", workerID, err)
			}

			plen := len(persistedJobs)

			if plen != len(tc.jobResult) {
				t.Errorf("len %d expected %d", plen, len(tc.jobResult))
			}
		})

	}

}
