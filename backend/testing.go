package backend

import (
	"testing"
	"time"

	"github.com/ottogiron/ferraritrunk/worker"
)

var persistCases = []struct {
	jobResult []*worker.JobResult
}{
	{[]*worker.JobResult{
		&worker.JobResult{
			ID:        "job_test",
			WorkerID:  "test_worker_1",
			Status:    worker.JobStatusSuccess,
			Output:    []byte("A test output"),
			StartTime: time.Now(),
			EndTime:   time.Now()},
	}},
}

func TestBackend(t *testing.T, backend Backend, teardown func(tb testing.TB)) {

	for _, tc := range persistCases {

		t.Run("Persist Jobs", func(t *testing.T) {
			err := backend.Persist(tc.jobResult)
			if err != nil {
				t.Errorf("backend.Persist() => err:%s is not expected", err)
			}
		})

		t.Run("Get JobResults", func(t *testing.T) {
			workerID := tc.jobResult[0].WorkerID
			persistedJobs, err := backend.JobResults(workerID)

			if err != nil {
				t.Fatalf("backend.JobResults(%s) error %s was not expected", workerID, err)
			}

			plen := len(persistedJobs)

			if plen != len(tc.jobResult) {
				t.Errorf("len %d expected %d", plen, len(tc.jobResult))
			}
		})

		t.Run("Get Job", func(t *testing.T) {
			singleJob := tc.jobResult[0]
			jobID := singleJob.ID
			jobResult, err := backend.Job(jobID)

			if err != nil {
				t.Errorf("backend.Job(%s) error %s was not expected", jobID, err)
			}

			if jobResult.String() != singleJob.String() {
				t.Errorf("backend.Job(%s) => %v expected %v", jobID, jobResult, singleJob)
			}
		})
		teardown(t)
	}

}
