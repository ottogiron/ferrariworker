package rabbit

import (
	"fmt"

	"github.com/ottogiron/ferraristream/stream"

	"github.com/streadway/amqp"
)

const (
	uriKey = "uri"
	//Queue Arguments
	queueNameKey           = "queue-name"
	queueDurableKey        = "queue-durable"
	queueDeleteWhenUsedKey = "queue-delete-when-used"
	queueExclusiveKey      = "queue-exclusive"
	queueNoWaitKey         = "queue-no-wait"
	queuArgumentsKey       = "queue-arguments"
	//Consumer arguments
	consumerTagKey       = "consumer-tag"
	consumerAutoAckKey   = "consumer-auto-ack"
	consumerExclusiveKey = "consumer-exclusive"
	consumerNoLocalKey   = "consumer-no-local"
	consumerNoWaitKey    = "consumer-no-wait"
	consumerArgsKey      = "consumer-args"
	//Exchange configurations
	exchangeNameKey               = "exchange"
	exchangeTypeKey               = "exchange-type"
	exchangeDurableKey            = "exchange-durable"
	exchangeDeleteWhenCompleteKey = "exchange-delete-when-complete"
	exchangeInternalKey           = "exchange-internal"
	exchangeNowaitKey             = "exchange-no-wait"
)

type factory struct{}

func (f *factory) New(config stream.Config) stream.Adapter {

	return newRabbitStreamAdapter(config)
}

//Register registers a stream adapter to  be used by the process command
func Register() {
	stream.RegisterStreamAdapterFactory(&factory{}, schema)
}

type rabbitStreamAdapter struct {
	config     stream.Config
	connection *amqp.Connection
}

//NewRabbitStreamAdapter creates a new rabbitStreamAdapter
func newRabbitStreamAdapter(config stream.Config) stream.Adapter {
	return &rabbitStreamAdapter{
		config: config,
	}
}

func (m *rabbitStreamAdapter) Open() error {
	queueURL := m.config.GetString(uriKey)
	con, err := amqp.Dial(queueURL)
	if err != nil {
		return fmt.Errorf("Error when dialing to %s %s", queueURL, err)
	}
	m.connection = con
	return nil
}

func (m *rabbitStreamAdapter) Close() error {
	return m.connection.Close()
}

func (m *rabbitStreamAdapter) StreamMessages(msgChannel chan<- *stream.Message) error {
	ch, err := m.connection.Channel()

	if err != nil {
		return fmt.Errorf("Could not open a channel  %s", err)
	}
	defer ch.Close()

	if m.config.GetString(exchangeNameKey) != "" {
		err := ch.ExchangeDeclare(
			m.config.GetString(exchangeNameKey),
			m.config.GetString(exchangeTypeKey),
			m.config.GetBoolean(exchangeDurableKey),
			m.config.GetBoolean(exchangeDeleteWhenCompleteKey),
			m.config.GetBoolean(exchangeInternalKey),
			m.config.GetBoolean(exchangeNowaitKey),
			nil,
		)
		if err != nil {
			return fmt.Errorf("Could not create an exchange %s", err)
		}
	}

	q, err := ch.QueueDeclare(
		m.config.GetString(queueNameKey),            // name
		m.config.GetBoolean(queueDurableKey),        // durable
		m.config.GetBoolean(queueDeleteWhenUsedKey), // delete when usused
		m.config.GetBoolean(queueExclusiveKey),      // exclusive
		m.config.GetBoolean(queueNoWaitKey),         // no-wait
		nil, // arguments
	)

	if err != nil {
		return fmt.Errorf("Could not declare the queue %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		m.config.GetString(consumerTagKey),        // consumer
		m.config.GetBoolean(consumerAutoAckKey),   // auto-ack
		m.config.GetBoolean(consumerExclusiveKey), // exclusive
		m.config.GetBoolean(consumerNoLocalKey),   // no-local
		m.config.GetBoolean(consumerNoWaitKey),    // no-wait
		nil, // args
	)

	if err != nil {
		return fmt.Errorf("Could not start consuming queue %s", err)
	}

	for d := range msgs {
		msgChannel <- stream.NewMessage(d.Body, d)
	}

	return nil
}

//RabbitResultHanlder post process when the job is already done
func (m *rabbitStreamAdapter) ResultHandler(jobResult *stream.JobResult, message *stream.Message) error {

	return nil
}
