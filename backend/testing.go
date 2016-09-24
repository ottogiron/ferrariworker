package backend

import (
	"testing"
	"time"

	"github.com/ottogiron/ferrariworker/processor"
)

var persistCases = []struct {
	jobResult []processor.JobResult
}{
	{[]processor.JobResult{
		processor.JobResult{
			Status:    processor.JobStatusSuccess,
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
