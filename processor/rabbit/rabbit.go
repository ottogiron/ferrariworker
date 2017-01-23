package rabbit

import (
	"context"
	"fmt"
	"log"

	"strings"

	"github.com/ferrariframework/ferrariworker/processor"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/streadway/amqp"
)

const (
	uriKey         = "uri"
	routingKey     = "routing_key"
	bindingWaitKey = "binding_wait"
	//Queue Arguments
	queueNameKey           = "queue_name"
	queueDurableKey        = "queue_durable"
	queueDeleteWhenUsedKey = "queue_delete_when_used"
	queueExclusiveKey      = "queue_exclusive"
	queueNoWaitKey         = "queue_no_wait"
	queuArgumentsKey       = "queue_arguments"
	//Consumer arguments
	consumerTagKey       = "consumer_tag"
	consumerAutoAckKey   = "consumer_auto_ack"
	consumerExclusiveKey = "consumer_exclusive"
	consumerNoLocalKey   = "consumer_no_local"
	consumerNoWaitKey    = "consumer_no_wait"
	consumerArgsKey      = "consumer_args"
	//Exchange configurations
	exchangeNameKey               = "exchange"
	exchangeTypeKey               = "exchange_type"
	exchangeDurableKey            = "exchange_durable"
	exchangeDeleteWhenCompleteKey = "exchange_delete_when_complete"
	exchangeInternalKey           = "exchange_internal"
	exchangeNowaitKey             = "exchange_no_wait"
	//Message configurations
	onSuccessKey = "on_sucess"
	onFailedKey  = "on_failed"
)

type factory struct{}

func (f *factory) New(config processor.AdapterConfig) processor.Adapter {
	return newRabbitProcessorAdapter(config)
}

func init() {
	processor.RegisterAdapterFactory(&factory{}, schema)
}

type rabbitProcessorAdapter struct {
	config     processor.AdapterConfig
	connection *amqp.Connection
}

//NewRabbitProcessorAdapter creates a new rabbitStreamAdapter
func newRabbitProcessorAdapter(config processor.AdapterConfig) processor.Adapter {
	return &rabbitProcessorAdapter{
		config: config,
	}
}

func (m *rabbitProcessorAdapter) Open() error {
	queueURL := m.config.GetString(uriKey)
	con, err := amqp.Dial(queueURL)
	if err != nil {
		return fmt.Errorf("Error when dialing to %s %s", queueURL, err)
	}
	m.connection = con
	return nil
}

func (m *rabbitProcessorAdapter) Close() error {
	return m.connection.Close()
}

func (m *rabbitProcessorAdapter) Messages(ctx context.Context) (<-chan processor.Message, error) {

	ch, err := m.connection.Channel()

	if err != nil {
		return nil, fmt.Errorf("Could not open a channel  %s", err)
	}

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
			return nil, fmt.Errorf("Could not create an exchange %s", err)
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
		return nil, fmt.Errorf("Could not declare the queue %s", err)
	}

	if err = ch.QueueBind(
		m.config.GetString(queueNameKey),    // name of the queue
		m.config.GetString(routingKey),      // bindingKey
		m.config.GetString(exchangeNameKey), // sourceExchange
		m.config.GetBoolean(bindingWaitKey), // noWait
		nil, // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
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
		return nil, fmt.Errorf("Could not start consuming queue %s", err)
	}

	msgChannel := make(chan processor.Message)
	go func() {
		for {
			select {
			case d := <-msgs:
				msgChannel <- processor.Message{Payload: d.Body, OriginalMessage: d}
			case <-ctx.Done():
				log.Println("Closing rabbit channel")
				ch.Cancel(m.config.GetString(consumerTagKey), true)
				ch.Close()
				close(msgChannel)
				return
			}
		}

	}()
	return msgChannel, nil
}

//RabbitResultHanlder post process when the job is already done
func (m *rabbitProcessorAdapter) ResultHandler(jobResult *processor.JobResult, message processor.Message) error {
	autoAck := m.config.GetBoolean(consumerAutoAckKey)
	//Only perform this logic if autoAck == false
	if !autoAck {
		originalMessage := message.OriginalMessage.(amqp.Delivery)
		switch jobResult.Status {
		case processor.JobStatusSuccess:
			return m.applyAction(jobResult, &originalMessage, onSuccessKey)

		case processor.JobStatusFailed:
			return m.applyAction(jobResult, &originalMessage, onFailedKey)
		}
	}

	return nil
}

func (m *rabbitProcessorAdapter) applyAction(jobResult *processor.JobResult, delivery *amqp.Delivery, actionKey string) error {
	messageAction, value := parseMessageAction(m.config.GetString(actionKey))
	if messageAction == "ack" {
		err := delivery.Ack(cast.ToBool(value))
		return errors.Wrapf(err, "Failed to apply ack on_success action for message workerID=%s jobID=%s", jobResult.WorkerID, jobResult.JobID)
	} else if messageAction == "reject" {
		err := delivery.Reject(cast.ToBool(value))
		return errors.Wrapf(err, "Failed to apply reject on_success action for message workerID=%s jobID=%s", jobResult.WorkerID, jobResult.JobID)
	} else if messageAction == "nack" {
		tokens := strings.Split(cast.ToString(value), ",")
		var multiple bool
		var requeue bool
		if len(tokens) == 2 {
			multiple = cast.ToBool(tokens[0])
			requeue = cast.ToBool(tokens[1])
		}
		err := delivery.Nack(multiple, requeue)
		if err != nil {
			return errors.Wrapf(err, "Failed to apply nack on_success action for message workerID=%s jobID=%s", jobResult.WorkerID, jobResult.JobID)
		}

	} else {
		//Apply default action
		if actionKey == onSuccessKey {
			delivery.Ack(false)
		} else if actionKey == onFailedKey {
			delivery.Ack(true)
		}
	}
	return nil
}

func parseMessageAction(action string) (string, interface{}) {
	tokens := strings.Split(action, ":")
	if len(tokens) == 2 {
		return tokens[0], tokens[1]
	}
	return "", false
}
