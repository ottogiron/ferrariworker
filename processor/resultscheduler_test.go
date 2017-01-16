package processor

import (
	"reflect"
	"testing"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/ottogiron/ferraritrunk/backend"
	"github.com/ottogiron/ferraritrunk/worker"
)

func TestNewResultScheduler(t *testing.T) {
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name string
		args args
		want ResultPersistScheduler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResultScheduler(tt.args.interval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResultScheduler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_backendScheduler_Schedule(t *testing.T) {
	type fields struct {
		jobResults []*worker.JobResult
		backend    backend.Backend
		interval   time.Duration
		ticker     *time.Ticker
		logger     log15.Logger
	}
	type args struct {
		jobResult *worker.JobResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := b.Schedule(tt.args.jobResult); (err != nil) != tt.wantErr {
				t.Errorf("backendScheduler.Schedule() error = %v, wantErr %v", err, tt.wantErr)
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
	type fields struct {
		jobResults []*worker.JobResult
		backend    backend.Backend
		interval   time.Duration
		ticker     *time.Ticker
		logger     log15.Logger
	}
	tests := []struct {
		name   string
		fields fields
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
			b.Start()
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
	type fields struct {
		jobResults []*worker.JobResult
		backend    backend.Backend
		interval   time.Duration
		ticker     *time.Ticker
		logger     log15.Logger
	}
	tests := []struct {
		name   string
		fields fields
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
			b.Stop()
		})
	}
}
