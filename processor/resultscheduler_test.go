package processor

import (
	"errors"
	"testing"
	"time"

	"reflect"

	"github.com/inconshreveable/log15"
	"github.com/ottogiron/ferraritrunk/backend"
	"github.com/ottogiron/ferraritrunk/worker"
)

func testLogger() log15.Logger {
	l := log15.New()
	l.SetHandler(log15.DiscardHandler())
	return log15.New()
}

func testBackend(wantErr bool) backend.Backend {
	return &mockBackend{wantErr: wantErr}
}

func newTestBackendScheduler(interval time.Duration, wantErr bool) ResultPersistScheduler {
	r := NewResultScheduler(
		testBackend(wantErr),
		testLogger(),
		interval,
	)

	return r
}

type mockBackend struct {
	wantErr bool
	backend.Backend
}

func (mb *mockBackend) Persist(jobResults []*worker.JobResult) error {
	if mb.wantErr {
		return errors.New("There was an error persisting")
	}
	return nil
}

func Test_backendScheduler_Schedule(t *testing.T) {

	type args struct {
		jobResult *worker.JobResult
	}
	tests := []struct {
		name string
		args args
		size int
	}{
		{"Schedule", args{&worker.JobResult{}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newTestBackendScheduler(time.Millisecond, false)
			b.Schedule(tt.args.jobResult)
			rs := b.(*backendScheduler)
			size := len(rs.jobResults)
			if len(rs.jobResults) != tt.size {
				t.Errorf("backendScheduler.Schedule() => size %d  wantErr %d", size, tt.size)
			}
		})
	}
}

func Test_backendScheduler_persistAllJobResults(t *testing.T) {
	type fields struct {
		jobResults []*worker.JobResult
		backend    backend.Backend
		interval   time.Duration
		ticker     *time.Ticker
		logger     log15.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &backendScheduler{

				jobResults: tt.fields.jobResults,
				backend:    tt.fields.backend,
				interval:   tt.fields.interval,
				ticker:     tt.fields.ticker,
				logger:     tt.fields.logger,
			}
			if err := b.persistAllJobResults(); (err != nil) != tt.wantErr {
				t.Errorf("backendScheduler.persistAllJobResults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backendScheduler_Start(t *testing.T) {

	tests := []struct {
		name         string
		jobResult    *worker.JobResult
		jobResultIn  []*worker.JobResult
		jobResultOut []*worker.JobResult
		wantErr      bool
	}{
		{"Start", &worker.JobResult{}, []*worker.JobResult{&worker.JobResult{}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			b := newTestBackendScheduler(time.Millisecond*1, tt.wantErr)
			b.Schedule(tt.jobResult)
			rs := b.(*backendScheduler)

			if !reflect.DeepEqual(rs.jobResults, tt.jobResultIn) {
				t.Errorf("b.Schedule() => rs.jobResults => %v want %v", rs.jobResults, tt.jobResultIn)
				return
			}
			b.Start()
			time.Sleep(2 * time.Millisecond)
			if !reflect.DeepEqual(rs.jobResults, tt.jobResultOut) {
				t.Errorf("b.Schedule() => rs.jobResults => %v want %v", rs.jobResults, tt.jobResultIn)
			}
		})
	}
}

func Test_backendScheduler_Flush(t *testing.T) {
	type fields struct {
		jobResults []*worker.JobResult
		backend    backend.Backend
		interval   time.Duration
		ticker     *time.Ticker
		logger     log15.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &backendScheduler{
				jobResults: tt.fields.jobResults,
				backend:    tt.fields.backend,
				interval:   tt.fields.interval,
				ticker:     tt.fields.ticker,
				logger:     tt.fields.logger,
			}
			if err := b.Flush(); (err != nil) != tt.wantErr {
				t.Errorf("backendScheduler.Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backendScheduler_Stop(t *testing.T) {

	tests := []struct {
		name string
	}{
		{"Flush"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newTestBackendScheduler(time.Millisecond*1, false)
			b.Stop()
		})
	}
}
