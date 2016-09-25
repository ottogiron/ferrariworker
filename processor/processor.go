package processor

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ottogiron/ferrariworker/adapter"
	"github.com/ottogiron/ferrariworker/worker"
)

//Processor defines a queue generic processor
type Processor interface {
	Start() error
}

type processor struct {
	config *Config
	stdout io.Writer
	stderr io.Writer
}

type job struct {
	command     string
	commandPath string
	payload     []byte
}

//JobResultHanlder Handler for the result of a processed Job
type JobResultHanlder func(jobResult *worker.JobResult, message *worker.Message)

//Config configuration for a Processor processor
type Config struct {
	Adapter     adapter.Adapter
	Command     string
	CommandPath string
	Concurrency int
	WaitTimeout time.Duration
}

//New creates a new instance of Processor processor
func New(config *Config, stdout io.Writer, stderr io.Writer) Processor {
	if stdout == nil {
		stdout = os.Stdout
	}

	if stderr == nil {
		stdout = os.Stderr
	}
	return &processor{config, stdout, stderr}
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
						sp.config.Adapter.ResultHandler(jobResult, m)
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

func (sp *processor) processJob(job job) *worker.JobResult {
	encodedPayload := base64.StdEncoding.EncodeToString(job.payload)
	cmdStr := job.command + " " + string(encodedPayload)
	var output bytes.Buffer
	cmd := sp.prepareCommand(cmdStr, job.commandPath, &output)
	err := cmd.Run()
	status := worker.JobStatusSuccess
	if err != nil {
		status = worker.JobStatusFailed
		errMsg := fmt.Sprintf("-Failed to run command  %s %s", job.command, err)
		output.WriteString(errMsg)
	} else {
		if success := cmd.ProcessState.Success(); !success {
			status = worker.JobStatusFailed
		}
	}
	jobResult := &worker.JobResult{
		Status: status,
		Output: output.Bytes(),
	}

	return jobResult
}

//prepareCommand executes a command outputs to stdout
func (sp *processor) prepareCommand(commandStr string, path string, output io.Writer) *exec.Cmd {
	cmdTokens := strings.Fields(commandStr)
	commandName := cmdTokens[:1]
	args := cmdTokens[1:]
	os.Chdir(path)
	cmd := exec.Command(commandName[0], args...)
	stdout := io.MultiWriter(output, sp.stdout)
	stderr := io.MultiWriter(output, sp.stderr)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}
