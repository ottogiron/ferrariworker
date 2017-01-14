package processor

import (
	"sync"

	"time"

	"github.com/inconshreveable/log15"
	"github.com/ottogiron/ferraritrunk/backend"
	"github.com/ottogiron/ferraritrunk/worker"
	"github.com/pkg/errors"
)

//ResultPersistScheduler schedules operations on top of a backend
type ResultPersistScheduler interface {
	Schedule(jobResult *worker.JobResult) error
	Start()
	Stop()
}

type backendScheduler struct {
	mu         sync.Mutex
	jobResults []*worker.JobResult
	backend    backend.Backend
	interval   time.Duration
	ticker     *time.Ticker
	logger     log15.Logger
}

func (b *backendScheduler) Persist(jobResult *worker.JobResult) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.jobResults = append(b.jobResults, jobResult)
	return nil
}

func (b *backendScheduler) persistAllJobResults() error {
	err := b.backend.Persist(b.jobResults)
	if err != nil {
		return errors.Wrap(err, "Failed to persist scheduled JobResults")
	}
	b.jobResults = nil
	return nil
}

func (b *backendScheduler) Start() {
	b.ticker = time.NewTicker(b.interval)
	go func() {
		for _ = range b.ticker.C {
			err := b.persistAllJobResults()
			if err != nil {
				b.logger.Error("Failed to persist job results", "location", "Result Persist Scheduler")
			}
		}
	}()
}

func (b *backendScheduler) Stop() {
	b.ticker.Stop()
}
