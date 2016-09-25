package processor

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ottogiron/ferraritrunk/worker"
	"github.com/ottogiron/ferrariworker/adapter"
)

var successfullJobs = []worker.Message{
	worker.Message{Payload: []byte("message 1")},
	worker.Message{Payload: []byte("message 2")},
	worker.Message{Payload: []byte("message 3")},
	worker.Message{Payload: []byte("message 4")},
	worker.Message{Payload: []byte("message 5")},
	worker.Message{Payload: []byte("message 6")},
}

type dummyWriter struct{}

func (dw *dummyWriter) Write(p []byte) (n int, err error) {
	return 1, nil
}

func createProcessorAdapterMock(t testing.TB, messages []worker.Message) *processorAdapterMock {
	return &processorAdapterMock{tb: t, messages: messages}
}

var failingJobs = []worker.Message{
	worker.Message{Payload: []byte("message for failing job 1")},
}

type processorAdapterFactoryMock struct {
	t *testing.T
}

func (s *processorAdapterMock) New(config *Config) adapter.Adapter {
	return createProcessorAdapterMock(s.tb, successfullJobs)
}

type processorAdapterMock struct {
	tb       testing.TB
	messages []worker.Message
}

func (s *processorAdapterMock) Open() error {
	return nil
}

func (s *processorAdapterMock) Close() error {
	return nil
}

func (s *processorAdapterMock) Messages(context context.Context) (<-chan worker.Message, error) {
	msgChannel := make(chan worker.Message)
	go func() {
		for _, message := range s.messages {
			msgChannel <- message
		}
	}()
	return msgChannel, nil
}

func (s *processorAdapterMock) ResultHandler(jobResult *worker.JobResult, message worker.Message) error {
	if jobResult.Status != worker.JobStatusSuccess {
		s.tb.Errorf("Running should be successful  status %d output %s", jobResult.Status, jobResult.Output)
		return fmt.Errorf("Running should be successful  status %d output %s", jobResult.Status, jobResult.Output)
	}
	originalMessage := base64.StdEncoding.EncodeToString(message.Payload)
	outputMessage := string(jobResult.Output)
	if !strings.Contains(outputMessage, originalMessage) {
		s.tb.Errorf("outputMessage => %s expected to contain %s", outputMessage, originalMessage)
	}
	return nil
}

type failProcessorAdapterMock struct {
	*processorAdapterMock
}

type failAdapterOpenCloseMock struct {
	*processorAdapterMock
}

func (s failAdapterOpenCloseMock) Open() error {
	return errors.New("Could not open the connection")
}

func (s failAdapterOpenCloseMock) Close() error {
	return errors.New("There was an erro closing the connection")
}

func (s *failProcessorAdapterMock) ResultHandler(jobResult *worker.JobResult, message worker.Message) error {
	if jobResult.Status != worker.JobStatusFailed {
		s.tb.Errorf("Running job should be unsuccesful %s", jobResult.Output)
	}
	return nil
}

func TestProcessorSuccessfulJobs(t *testing.T) {

	processorAdapterMock := &processorAdapterMock{tb: t, messages: successfullJobs}
	processorConfig := &Config{
		Adapter:     processorAdapterMock,
		Command:     `echo "Hello successful test"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}
	var w bytes.Buffer
	processor := New(processorConfig, &w, &w)
	processor.Start()
}

func BenchmarkProcessorSuccessfulJobs(b *testing.B) {
	processorAdapterMock := &processorAdapterMock{tb: b, messages: successfullJobs}
	processorConfig := &Config{
		Adapter:     processorAdapterMock,
		Command:     `echo "Hello successful test"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}

	var w bytes.Buffer
	processor := New(processorConfig, &w, &w)
	for n := 0; n < b.N; n++ {
		processor.Start()
	}
}

func TestProcessorFailingJob(t *testing.T) {
	adapterMock := createProcessorAdapterMock(t, failingJobs)
	failStreamAdapterMock := &failProcessorAdapterMock{adapterMock}

	processorConfig := &Config{
		Adapter:     failStreamAdapterMock,
		Command:     `cd nonexistingdir"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}
	w := &dummyWriter{}
	processor := New(processorConfig, w, w)
	processor.Start()
}

func BenchmarkProcessorFailedJobs(b *testing.B) {
	adapterMock := createProcessorAdapterMock(b, failingJobs)
	failStreamAdapterMock := &failProcessorAdapterMock{adapterMock}

	processorConfig := &Config{
		Adapter:     failStreamAdapterMock,
		Command:     `cd nonexistingdir"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}
	w := &dummyWriter{}
	processor := New(processorConfig, w, w)
	for n := 0; n < b.N; n++ {
		processor.Start()
	}
}

func TestNewMessage(t *testing.T) {
	messageStr := "hello world"
	m := worker.Message{
		Payload:         []byte(messageStr),
		OriginalMessage: nil,
	}

	if string(m.Payload) != messageStr {
		t.Fatalf("Message should be %s was %s", messageStr, string(m.Payload))
	}
}

func TestAdapterOpenError(t *testing.T) {
	adapterMock := createProcessorAdapterMock(t, successfullJobs)
	failAdapterOpenCloseMock := &failAdapterOpenCloseMock{adapterMock}
	processorConfig := &Config{
		Adapter: failAdapterOpenCloseMock,
	}
	sp := New(processorConfig, nil, nil)
	err := sp.Start()
	if err == nil {
		t.Error("Expected connection open to fail")
	}
}
