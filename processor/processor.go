package processor

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ferrariframework/ferrariserver/grpc/gen"
)

//JobStatus represents a job final status
type JobStatus int

const (
	//JobStatusSuccess represents a job success
	JobStatusSuccess JobStatus = 0
	//JobStatusFailed represents a failed job
	JobStatusFailed JobStatus = 1
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
	config           *Config
	stdout           io.Writer
	stderr           io.Writer
	jobServiceClient gen.JobServiceClient
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
	WorkerID string
	JobID    string
	Status   JobStatus
	Output   []byte
}

//JobResultHanlder Handler for the result of a processed Job
type JobResultHanlder func(jobResult *JobResult, message *Message)

//Config configuration for a Processor processor
type Config struct {
	WorkerID    string
	Adapter     Adapter
	Command     string
	CommandPath string
	Concurrency int
	WaitTimeout time.Duration
}

//New creates a new instance of Processor processor
func New(config *Config, jobServiceClient gen.JobServiceClient, stdout io.Writer, stderr io.Writer) Processor {
	if stdout == nil {
		stdout = os.Stdout
	}

	if stderr == nil {
		stderr = os.Stderr
	}
	return &processor{config: config, jobServiceClient: jobServiceClient, stdout: stdout, stderr: stderr}
}

//Start starts processing
func (sp *processor) Start() error {
	//open the connection
	err := sp.config.Adapter.Open()
	if err != nil {
		return fmt.Errorf("Failed to open the processor Adapter connection %s", err)
	}
	defer sp.config.Adapter.Close()
	wg := sync.WaitGroup{}
	//Wait for the timeout once then call done to exit the processing
	wg.Add(sp.config.Concurrency)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgs, err := sp.config.Adapter.Messages(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get messages from adapter %s", err)

	}
	for i := 0; i < sp.config.Concurrency; i++ {
		go func() {
			for {
				select {
				case m, ok := <-msgs:
					if ok {
						j := job{sp.config.Command, sp.config.CommandPath, m.Payload}
						jobResult := sp.processJob(j)
						err := sp.config.Adapter.ResultHandler(jobResult, m)
						if err != nil {
							log.Print("Error when running adapter ResultHandler \n", err)
						}
					} else {

						wg.Done()
						return
					}
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

func (sp *processor) processJob(job job) *JobResult {
	encodedPayload := base64.StdEncoding.EncodeToString(job.payload)
	cmdStr := job.command + " " + string(encodedPayload)
	var output bytes.Buffer

	sjob, err := sp.jobServiceClient.RegisterJob(context.Background(), &gen.Job{
		WorkerId: sp.config.WorkerID,
	})

	if err != nil {
		log.Printf("Failed to register job workerID %s", sp.config.WorkerID)
	}
	cmd, close := sp.prepareCommand(cmdStr, job.commandPath, sjob.GetWorkerId(), sjob.GetId(), &output)
	defer close()
	err = cmd.Run()
	status := JobStatusSuccess
	if err != nil {
		status = JobStatusFailed
		errMsg := fmt.Sprintf("Failed to run command  %s %s", job.command, err)
		log.Print(errMsg)
		output.WriteString(errMsg)
	} else {
		if success := cmd.ProcessState.Success(); !success {
			status = JobStatusFailed
		}
	}
	jobResult := &JobResult{
		WorkerID: sjob.GetWorkerId(),
		JobID:    sjob.GetId(),
		Status:   status,
		Output:   output.Bytes(),
	}

	return jobResult
}

//prepareCommand executes a command outputs to stdout
func (sp *processor) prepareCommand(commandStr, path, workerID, jobID string, output io.Writer) (*exec.Cmd, func()) {
	cmdTokens := strings.Fields(commandStr)
	commandName := cmdTokens[:1]
	args := cmdTokens[1:]
	os.Chdir(path)
	cmd := exec.Command(commandName[0], args...)

	stdout := io.MultiWriter(output, sp.stdout)
	stderr := io.MultiWriter(output, sp.stderr)

	logStream, err := sp.jobServiceClient.RecordLog(context.Background())

	if err != nil {
		log.Printf("Failed to get stdout stream for job workerID=%s jobID=%s", workerID, jobID)
	}

	stdoutForwarder := NewJobLogForwarder(workerID, jobID, logStream, stdout)

	stdErrForwarder := NewJobLogForwarder(workerID, jobID, logStream, stderr)

	cmd.Stdout = stdoutForwarder
	cmd.Stderr = stdErrForwarder

	return cmd, func() {
		logStream.CloseAndRecv()
	}
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
