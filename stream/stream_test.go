package stream

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

func createStreamAdapterMock(t *testing.T, messages []Message) *streamAdapterMock {
	return &streamAdapterMock{t, messages}
}

var failingJobs = []Message{
	Message{Payload: []byte("message for failing job 1")},
}

type streamAdapterFactoryMock struct {
	t *testing.T
}

func (s *streamAdapterMock) New(config *Config) Adapter {
	return createStreamAdapterMock(s.t, successfullJobs)
}

type streamAdapterMock struct {
	t        *testing.T
	messages []Message
}

func (s *streamAdapterMock) Open() error {
	return nil
}

func (s *streamAdapterMock) Close() error {
	return nil
}

func (s *streamAdapterMock) StreamMessages(msgChannel chan<- *Message) error {
	for _, message := range s.messages {
		msgChannel <- &message
	}
	return nil
}

func (s *streamAdapterMock) ResultHandler(jobResult *JobResult, message *Message) error {
	if !jobResult.Success {
		s.t.Fatalf("Running should be successful %s", jobResult.Output)
		return fmt.Errorf("Running should be successful %s", jobResult.Output)
	}
	return nil
}

type failStreamAdapterMock struct {
	*streamAdapterMock
}

type failAdapterOpenCloseMock struct {
	*streamAdapterMock
}

func (s failAdapterOpenCloseMock) Open() error {
	return errors.New("Could not open the connection")
}

func (s failAdapterOpenCloseMock) Close() error {
	return errors.New("There was an erro closing the connection")
}

func (s *failStreamAdapterMock) ResultHandler(jobResult *JobResult, message *Message) error {
	if jobResult.Success {
		s.t.Fatalf("Running job should be unsuccesful %s", jobResult.Output)
	}
	return nil
}

func TestProcessorSuccessfulJobs(t *testing.T) {

	streamAdapterMock := &streamAdapterMock{t, successfullJobs}
	processorConfig := &ProcessorConfig{
		Adapter:     streamAdapterMock,
		Command:     `echo "Hello successful test"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}

	streamProcessor := NewStreamProcessor(processorConfig)
	streamProcessor.Start()
}

func TestFailingJob(t *testing.T) {
	adapterMock := createStreamAdapterMock(t, failingJobs)
	failStreamAdapterMock := &failStreamAdapterMock{adapterMock}

	processorConfig := &ProcessorConfig{
		Adapter:     failStreamAdapterMock,
		Command:     `cd nonexistingdir"`,
		CommandPath: ".",
		Concurrency: 1,
		WaitTimeout: 200,
	}

	streamProcessor := NewStreamProcessor(processorConfig)
	streamProcessor.Start()
}

func TestNewMessage(t *testing.T) {
	messageStr := "hello world"
	m := NewMessage([]byte(messageStr), nil)

	if string(m.Payload) != messageStr {
		t.Fatalf("Message should be %s was %s", messageStr, string(m.Payload))
	}
}

func TestAdapterOpenError(t *testing.T) {
	adapterMock := createStreamAdapterMock(t, successfullJobs)
	failAdapterOpenCloseMock := &failAdapterOpenCloseMock{adapterMock}
	processorConfig := &ProcessorConfig{
		Adapter: failAdapterOpenCloseMock,
	}
	sp := NewStreamProcessor(processorConfig)
	err := sp.Start()
	if err == nil {
		t.Error("Expected connection open to fail")
	}
}

func TestRegisterStreamAdapterFactory(t *testing.T) {

	cs := &ConfigurationSchema{
		Name: "test",
	}

	err := RegisterStreamAdapterFactory(nil, cs)

	if err != nil {
		t.Errorf("The first factory should be registered correctly for %s", cs.Name)
	}

	err = RegisterStreamAdapterFactory(nil, cs)

	if err == nil {
		t.Errorf("The registration should fail for %s", cs.Name)
	}
}

func TestGetConfigurationSchemas(t *testing.T) {
	RegisterStreamAdapterFactory(nil, &ConfigurationSchema{
		Name: "test",
	})

	schemas := StreamAdapterSchemas()

	slen := len(schemas)
	if slen != 1 {
		t.Errorf("Expected schemas size %d was", slen)
	}
}

func TestGetStreamAdapterSchema(t *testing.T) {
	RegisterStreamAdapterFactory(nil, &ConfigurationSchema{Name: "test"})

	schema, _ := StreamAdapterSchema("test")

	if schema.Name != "test" {
		t.Errorf("expected schema name to be 'test' was %s", schema.Name)
	}
}
