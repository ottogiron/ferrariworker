package worker

import "time"

//Message A generic message to be processed by a job
type Message struct {
	Payload         []byte
	OriginalMessage interface{}
}

//JobResult Represents the result of a processed Job
type JobResult struct {
	ID        string
	WorkerID  string
	Status    JobStatus
	Output    []byte
	StartTime time.Time
	EndTime   time.Time
}

type JobStatus int

const (
	JobStatusSuccess JobStatus = 0
	JobStatusFailed  JobStatus = 1
)
