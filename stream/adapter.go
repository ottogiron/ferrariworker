package stream

//Adapter defines an stream source
type Adapter interface {
	Open() error
	Close() error
	StreamMessages(chan<- *Message) error
	ResultHandler(jobResult *JobResult, message *Message) error
}

//Factory defines a actory from stream adapters
type Factory interface {
	New(config Config) Adapter
}
