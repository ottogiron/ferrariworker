package backend

import "github.com/ottogiron/ferrariworker/processor"

//Backend defines a data store for jobs to be persisted
type Backend interface {
	Persist([]processor.JobResult) error
}
