# Ferrari Worker

## Process   programs in any language from different  sources. 

[![Build Status](https://travis-ci.org/ottogiron/ferrariworker.svg?branch=master)](https://travis-ci.org/ottogiron/ferrariworker)
## Example
This example shows jobs processing from rabbitmq using a Node.js script as job processor.

```
ferrariworker process rabbit \
    --uri=amqp://guest:guest@localhost:5672 \
    --queue-name=hello \
    --consumer-auto-ack=true \
    --command="node hello.js" \
    --command-run-path="/Users/ogiron" \
    --max-concurrency=8 
```

### Example job script for Node

***hello.js***

```javascript
console.time("hello task");

const payloadbase64 = process.argv[2];`
var buf = new Buffer(payloadbase64, 'base64');
console.log(buf.toString());
console.timeEnd("hello task")
process.exit(0);
/* Task was successful */
```

### Example output

```
...
hello 96
hello task: 3ms
2016/05/26 10:48:07 Thread 2 from 8 end processing took 169.387777ms
2016/05/26 10:48:07 Thread 7 from 8 end processing took 197.641664ms
2016/05/26 10:48:07 Thread 8 from 8 end processing took 198.734157ms
hello 99
hello task: 2ms
2016/05/26 10:48:07 Thread 3 from 8 end processing took 136.934854ms

```
## Processor Adapters

Processors  Adapters are in charge of consuming jobs from different sources.

```
ferrariworker process <processor_adapter_name> <[flags]>
```

### Available Processor Adapters

* [RabbitMQ](processor/rabbit)

### Processing

## Job Payload

The program will receive the payload of a job as an argument encoded in base64. 

### Example in Node

```
//You can get a job payload using arguments array from the language you are using. 
const payloadbase64 = process.argv[2];
var buf = new Buffer(payloadbase64, 'base64');
console.log(buf.toString());
```

## Global Flags
This flags apply to all the available processors

<table>
    <tr>
        <th>Variable</th>
        <th>Default</th>
        <th>Description</th>
    </tr>
      <tr>
        <td>command</td>
        <td>echo "You should  set your own command :)"</td>
        <td>Command to be run when a new job arrives to the queue.</td>
    </tr>
      <tr>
        <td>command-path</td>
        <td>.</td>
        <td>Context path of the command to be run.</td>
    </tr>
      <tr>
        <td>max-concurrency</td>
        <td>1</td>
        <td>Max number of jobs proccessed concurrently.</td>
    </tr>
      <tr>
        <td>wait-timeout</td>
        <td>200</td>
        <td>Time to wait in milliseconds for existing jobs to be processed. </td>
    </tr>
</table>

## Development

### Processor Message

Represents a generic processor message

```
//Message A generic message to be processed by a job
type Message struct {
	Payload         []byte
	OriginalMessage interface{}
}
```

### Processor Adapters
An adapter is an interface that defines the functionallity for processing jobs from any source.

```
//Adapter defines an processor source
type Adapter interface {
	Open() error
	Close() error
  Messages() (<-chan Message, error)
	ResultHandler(jobResult *JobResult, message Message) error
}
```

#### Processor Adapter Factory
An processor adapter factory is interface which defines a "New" method that the ferrari core will use to create an instance of an specific adapter.
"New" receives   a "Config" object which contains all the configuration values provided by the user.

```
//Factory defines a actory from stream adapters
type Factory interface {
	New(config *AdapterConfig) Adapter
}
```
#### Processor Adapter Factory Registration

Every Adapter should register itself in the processor adapter registry

Example: 

```
func init() {
	processor.RegisterAdapterFactory(&factory{}, schema)
}

Then the adapter has to be imported in ***cmd/modules.go** using the blank identifier for the processor to be registered



package cmd

import (
  _ "github.com/ottogiron/ferrariworker/processor/rabbit"
  )


``` 

#### Processor Adapters Configuration Schema
Every adapters has to define their configuration metadata, that means the adapter name/identifier and all the related configuration fields.
This information is necessary for the adapter to be registered as processor command, and to be able to parse the configuration values that will be provided to the factory.


Please check the [RabbitMQ](processor/rabbit/rabbit.go) adapter for an example of a working processor adapter.


### Prerequisites

Go 1.6 +

### Makefile

 **targets**

* **all**: Runs tests an binaries
* **build-release**: Build a docker container and binaries
* **lint**: Runs lint tools (fmt, vet)
* **test**: Runs unit tests

**Please check each target for details**
