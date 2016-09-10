package processor

import (
	"errors"
	"fmt"
	"testing"
)

var successfullJobs = []Message{
	Message{Payload: []byte("message 1")},
	Message{Payload: []byte("message 2")},
	Message{Payload: []byte("message 3")},
	Message{Payload: []byte("message 4")},
	Message{Payload: []byte("message 5")},
	Message{Payload: []byte("message 6")},
}

func createProcessorAdapterMock(t *testing.T, messages []Message) *processorAdapterMock {
	return &processorAdapterMock{t: t, messages: messages}
}

var failingJobs = []Message{
	Message{Payload: []byte("message for failing job 1")},
}

type processorAdapterFactoryMock struct {
	t *testing.T
}

func (s *processorAdapterMock) New(config *Config) Adapter {
	return createProcessorAdapterMock(s.t, successfullJobs)
}

type processorAdapterMock struct {
	t        *testing.T
	messages []Message
}

func (s *processorAdapterMock) Open() error {
	return nil
}

func (s *processorAdapterMock) Close() error {
	return nil
}

func (s *processorAdapterMock) Messages() (<-chan Message, error) {
	msgChannel := make(chan Message)
	go func() {
		for _, message := range s.messages {
			msgChannel <- message
		}
	}()
	return msgChannel, nil
}

func (s *processorAdapterMock) ResultHandler(jobResult *JobResult, message Message) error {
	if jobResult.Status != JobStatusSuccess {
		s.t.Errorf("Running should be successful %s", jobResult.Output)
		return fmt.Errorf("Running should be successful %s", jobResult.Output)
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

func (s *failProcessorAdapterMock) ResultHandler(jobResult *JobResult, message Message) error {
	if jobResult.Status != JobStatusSuccess {
		s.t.Errorf("Running job should be unsuccesful %s", jobResult.Output)
	}
	return nil
}

func TestProcessorSuccessfulJobs(t *testing.T) {

	processorAdapterMock := &processorAdapterMock{t: t, messages: successfullJobs}
	processorConfig := &Config{
		Adapter:     processorAdapterMock,
		Command:     `echo "Hello successful test"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}

	processor := New(processorConfig)
	processor.Start()
}

func TestNewMessage(t *testing.T) {
	messageStr := "hello world"
	m := Message{[]byte(messageStr), nil}

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
	sp := New(processorConfig)
	err := sp.Start()
	if err == nil {
		t.Error("Expected connection open to fail")
	}
}

func TestRegisterAdapterFactory(t *testing.T) {

	cs := &AdapterConfigurationSchema{
		Name: "test",
	}

	err := RegisterAdapterFactory(nil, cs)

	if err != nil {
		t.Errorf("The first factory should be registered correctly for %s", cs.Name)
	}

	err = RegisterAdapterFactory(nil, cs)

	if err == nil {
		t.Errorf("The registration should fail for %s", cs.Name)
	}
}

func TestGetConfigurationSchemas(t *testing.T) {
	RegisterAdapterFactory(nil, &AdapterConfigurationSchema{
		Name: "test",
	})

	schemas := AdapterSchemas()

	slen := len(schemas)
	if slen != 1 {
		t.Errorf("Expected schemas size %d was", slen)
	}
}

func TestGetProcessorAdapterSchema(t *testing.T) {
	RegisterAdapterFactory(nil, &AdapterConfigurationSchema{Name: "test"})

	schema, _ := AdapterSchema("test")

	if schema.Name != "test" {
		t.Errorf("expected schema name to be 'test' was %s", schema.Name)
	}
}
