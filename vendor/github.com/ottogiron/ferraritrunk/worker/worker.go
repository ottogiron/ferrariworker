package worker

import (
	"encoding/json"
	"time"
)

//Message A generic message to be processed by a job
type Message struct {
	Payload         []byte
	OriginalMessage interface{}
}

//JobResult Represents the result of a processed Job
type JobResult struct {
	ID        string    `json:"job_id"`
	WorkerID  string    `json:"worker_id"`
	Status    JobStatus `json:"status"`
	Output    []byte    `json:"output"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (j *JobResult) String() string {
	b, _ := json.Marshal(j)
	return string(b)
}

type JobStatus int

const (
	JobStatusSuccess JobStatus = 0
	JobStatusFailed  JobStatus = 1
)
