# Ferrari Stream

## Creates a continuous stream of jobs to be processed by custom programs in any language from different stream sources. 

## Example
This example shows jobs processing from rabbitmq using a Node.js script as job processor.

```
ferraristream process rabbit \
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
## Streams
Streams are the sources from which you can process jobs e.g rabbit.

```
ferraristream process <stream_name> <[flags]>
```

### Available Streams

* [RabbitMQ](/users/ogiron/repos/ferraristream/browse/stream/rabbit)

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
This flags apply to all the available streams

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

### Stream Adapters
An adapter is an interface that defines the functionallity for streaming jobs from any source.

```
//Adapter defines an stream source
type Adapter interface {
	Open() error
	Close() error
	StreamMessages(chan<- *Message) error
	ResultHandler(jobResult *JobResult, message *Message) error
}
```

#### Stream Factory
An stream factory is interface which defines a "New" method that the ferrari core will use to create an instance of an specific adapter.
"New" receives   a "Config" object which contains all the configuration values provided by the user.

```
//Factory defines a actory from stream adapters
type Factory interface {
	New(config *Config) Adapter
}
```
#### Stream Factory Registration
By convention every adapter must implement a package level function called Register, which will be in charge of registering the adapter factory.

Example: 

```
//Register registers a stream adapter to  be used by the process command
func Register() {
	configSchema := &stream.ConfigurationSchema{
		Name:             "rabbit",
		ShortDescription: "Executes jobs coming from rabbit",
		LongDescription:  "Executes jobs coming from rabbit",
		Properties: []stream.ConfigurationProperty{
			stream.ConfigurationProperty{
				Name:        hostPropertyKey,
				Type:        stream.PropertyTypeString,
				Description: "Rabbit host url e.g. amqp://guest:guest@localhost:5672/",
				Default:     "amqp://guest:guest@localhost:5672/",
				Optional:    false,
			},
		},
	}
	stream.RegisterStreamAdapterFactory(&factory{}, configSchema)
}
```

This function needs to be called explictly in cmd/modules.go


```
package cmd

import "github.com/ottogiron/ferraristream/stream/rabbit"

func init() {
	rabbit.Register()
  // Other adapters registration
}
``` 

#### Stream Configuration Schema
Every adapters has to define their configuration metadata, that means the adapter name/identifier and all the related configuration fields.
This information is necessary for the adapter to be registered as stream command, and to be able to parse the configuration values that will be provided to the factory.


Please check the [RabbitMQ](stream/rabbit/rabbit.go) adapter for an example of a working stream adapter.


### Prerequisites

Go 1.6 +

### Makefile

 **targets**

* **all**: Runs tests an binaries
* **build-release**: Build a docker container and binaries
* **lint**: Runs lint tools (fmt, vet)
* **test**: Runs unit tests
* **release-docker**: Pushes docker image to artifactory
* **binaries**: Builds binaries for linux and darwin
* **packages**: Builds packages for rpm and brew
* **release**: releases rpm and brew packages

**Please check each target for details**
