package processor

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type JobStatus int

const (
	JobStatusSuccess JobStatus = iota
	JobStatusFailed
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

//JobResult Represents the result of a processed Job
type JobResult struct {
	Status JobStatus
	Output []byte
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

	msgs, err := sp.config.Adapter.Messages()
	if err != nil {
		return fmt.Errorf("Failed to get messages from adapter %s", err)
	}

	for i := 0; i < sp.config.Concurrency; i++ {
		go func() {
			for {
				select {
				case m := <-msgs:
					j := job{sp.config.Command, sp.config.CommandPath, m.Payload}
					jobResult := processJob(j)
					sp.config.Adapter.ResultHandler(jobResult, m)
				case <-time.After(sp.config.WaitTimeout * time.Millisecond):
					wg.Done()
					return
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func processJob(job job) *JobResult {
	encodedPayload := base64.StdEncoding.EncodeToString(job.payload)
	cmdStr := job.command + " " + string(encodedPayload)
	var output bytes.Buffer
	cmd := prepareCommand(cmdStr, job.commandPath, &output)
	err := cmd.Run()
	jobResult := &JobResult{}
	jobResult.Output = output.Bytes()
	if err != nil {
		jobResult.Status = JobStatusFailed
	} else {
		jobResult.Status = JobStatusSuccess
		if success := cmd.ProcessState.Success(); !success {
			jobResult.Status = JobStatusFailed
		}
	}
	return jobResult
}

//prepareCommand executes a command outputs to stdout
func prepareCommand(commandStr string, path string, output io.Writer) *exec.Cmd {
	cmdTokens := strings.Fields(commandStr)
	commandName := cmdTokens[:1]
	args := cmdTokens[1:]
	os.Chdir(path)
	cmd := exec.Command(commandName[0], args...)
	stdout := io.MultiWriter(output, os.Stdout)
	stderr := io.MultiWriter(output, os.Stderr)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

//RegisterAdapterFactory registers a new factory for creating adapters
func RegisterAdapterFactory(factory Factory, adapterConfigurationSchema *AdapterConfigurationSchema) error {
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
