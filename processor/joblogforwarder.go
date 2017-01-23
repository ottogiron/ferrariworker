package processor

import (
	"io"

	"github.com/ferrariframework/ferrariserver/grpc/gen"
	"github.com/pkg/errors"
)

var _ io.Writer = (*JobLogForwarder)(nil)

//JobLogForwarder forward job logs to a ferrariworker server and wraps an existing writer
type JobLogForwarder struct {
	workerID           string
	jobID              string
	jobRecordLogStream gen.JobService_RecordLogClient
	writer             io.Writer
}

//NewJobLogForwarder returns a new instance of a JobLogForwarder
func NewJobLogForwarder(workerID, jobID string, jobRecordLogStream gen.JobService_RecordLogClient, writer io.Writer) io.Writer {
	return &JobLogForwarder{
		workerID:           workerID,
		jobID:              jobID,
		jobRecordLogStream: jobRecordLogStream,
		writer:             writer,
	}
}

func (j *JobLogForwarder) Write(b []byte) (int, error) {

	//First try to write to the wrapped writer
	n, err := j.writer.Write(b)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to write to wrapped writer workerID=%s jobID=%s", j.workerID, j.jobID)
	}

	if n != len(b) {
		err = io.ErrShortWrite
		return 0, err
	}

	err = j.jobRecordLogStream.Send(&gen.Log{
		WorkerId: j.workerID,
		JobId:    j.jobID,
		Message:  b,
	})

	if err != nil {
		return 0, errors.Wrapf(err, "Failed to record log for workerID=%s jobID=%s", j.workerID, j.jobID)
	}
	return len(b), nil
}
