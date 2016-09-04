package processor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var factoryRegistry = map[string]*configurationRegistry{}

//Processor defines a queue generic processor
type Processor interface {
	Start() error
}

type configurationRegistry struct {
	factory             Factory
	configurationSchema *AdapterConfigurationSchema
}

//RegisterProcessorAdapterFactory registers a new factory for creating adapters
func RegisterProcessorAdapterFactory(factory Factory, adapterConfigurationSchema *AdapterConfigurationSchema) error {
	if factoryRegistry[adapterConfigurationSchema.Name] != nil {
		return fmt.Errorf("The factory already exists %s", adapterConfigurationSchema.Name)
	}
	factoryRegistry[adapterConfigurationSchema.Name] = &configurationRegistry{factory, adapterConfigurationSchema}
	return nil
}

//AdapterFactory returns a factory for a registered adapter
func AdapterFactory(factoryName string) (Factory, error) {
	if factoryRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return factoryRegistry[factoryName].factory, nil
}

//AdapterSchema returns the configuration schema for the adapter factory
func AdapterSchema(factoryName string) (*AdapterConfigurationSchema, error) {
	if factoryRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return factoryRegistry[factoryName].configurationSchema, nil
}

//AdapterSchemas returns the schemas for all the available adapters
func AdapterSchemas() []*AdapterConfigurationSchema {
	configurationSchemas := []*AdapterConfigurationSchema{}
	for _, registry := range factoryRegistry {
		configurationSchemas = append(configurationSchemas, registry.configurationSchema)
	}
	return configurationSchemas
}

type processor struct {
	config *Config
}

type job struct {
	command     string
	commandPath string
	payload     []byte
}

//Message A generic message to be processed by a job
type Message struct {
	Payload         []byte
	OriginalMessage interface{}
}

//NewMessage creates a new instance of a message
func NewMessage(payload []byte, originalMessage interface{}) *Message {
	return &Message{
		payload,
		originalMessage,
	}
}

//JobResult Represents the result of a processed Job
type JobResult struct {
	Success bool
	Output  []byte
}

//JobResultHanlder Handler for the result of a processed Job
type JobResultHanlder func(jobResult *JobResult, message *Message)

//Config configuration for a Processor processor
type Config struct {
	Adapter     Adapter
	Command     string
	CommandPath string
	Concurrency int
	WaitTimeout time.Duration
}

//New creates a new instance of Processor processor
func New(config *Config) Processor {
	return &processor{config}
}

//Start starts processing
func (sp *processor) Start() error {
	//open the connection
	err := sp.config.Adapter.Open()
	if err != nil {
		return errors.New("Couldn't open the Processor connection")
	}
	defer sp.config.Adapter.Close()
	wg := sync.WaitGroup{}
	wg.Add(sp.config.Concurrency)
	msgs := make(chan *Message)
	go sp.config.Adapter.Messages(msgs)
	for i := 0; i < sp.config.Concurrency; i++ {
		go func(threadNumber int) {
			for {
				select {
				case m := <-msgs:
					start := time.Now()
					j := &job{sp.config.Command, sp.config.CommandPath, m.Payload}
					jobResult := processJob(j)
					sp.config.Adapter.ResultHandler(jobResult, m)
					elapsed := time.Since(start)
					log.Printf("Thread %d from %d end processing took %s", threadNumber+1, sp.config.Concurrency, elapsed)
				case <-time.After(sp.config.WaitTimeout * time.Millisecond):
					wg.Done()
					return
				}
			}
		}(i)
	}

	wg.Wait()
	return nil
}

func processJob(job *job) *JobResult {
	encodedPayload := base64.StdEncoding.EncodeToString(job.payload)
	cmdStr := job.command + " " + string(encodedPayload)
	cmd, _ := executeCommand(cmdStr, job.commandPath)
	jobResult := new(JobResult)
	jobResult.Success = cmd.ProcessState.Success()
	jobResult.Output, _ = cmd.Output()
	return jobResult
}

//ExecuteCommand executes a command outputs to stdout
func executeCommand(commandStr string, path string) (*exec.Cmd, error) {
	cmdTokens := strings.Fields(commandStr)
	commandName := cmdTokens[:1]
	args := cmdTokens[1:]
	os.Chdir(path)
	cmd := exec.Command(commandName[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return cmd, err
}
