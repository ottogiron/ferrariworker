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
	Schedule(jobResult *worker.JobResult)
	Flush() error
	Start()
	Stop()
}

//NewResultScheduler returns a new instance of a result scheduler
func NewResultScheduler(backend backend.Backend, logger log15.Logger, interval time.Duration) ResultPersistScheduler {
	return &backendScheduler{
		backend:  backend,
		logger:   logger,
		interval: interval,
	}
}

type backendScheduler struct {
	mu         sync.Mutex
	jobResults []*worker.JobResult
	backend    backend.Backend
	interval   time.Duration
	ticker     *time.Ticker
	logger     log15.Logger
}

func (b *backendScheduler) Schedule(jobResult *worker.JobResult) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.jobResults = append(b.jobResults, jobResult)

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

func (b *backendScheduler) Flush() error {
	err := b.persistAllJobResults()
	return errors.Wrap(err, "Failed to Flush the remaining jobResults")
}

func (b *backendScheduler) Stop() {
	b.ticker.Stop()
}
