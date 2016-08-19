package rabbit

import "github.com/ottogiron/ferraristream/stream"

var schema = &stream.ConfigurationSchema{
	Name:             "rabbit",
	ShortDescription: "Streams jobs from rabbitmq",
	LongDescription:  "Streams jobs from rabbibmq",
	Properties: []stream.ConfigurationProperty{
		stream.ConfigurationProperty{
			Name:        uriKey,
			Type:        stream.PropertyTypeString,
			Description: "Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/",
			Default:     "amqp://guest:guest@localhost:5672/",
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        queueNameKey,
			Type:        stream.PropertyTypeString,
			Description: "Rabbit queue name",
			Optional:    false,
		},
		stream.ConfigurationProperty{
			Name:        queueDurableKey,
			Type:        stream.PropertyTypeBool,
			Description: "Queue durable",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        queueDeleteWhenUsedKey,
			Type:        stream.PropertyTypeBool,
			Description: "Queue delete queue when used",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        queueExclusiveKey,
			Type:        stream.PropertyTypeBool,
			Description: "Queue exclusive",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        queueNoWaitKey,
			Type:        stream.PropertyTypeBool,
			Description: "Queue no wait",
			Default:     false,
			Optional:    true,
		},
		//Consumer configurations
		stream.ConfigurationProperty{
			Name:        consumerTagKey,
			Type:        stream.PropertyTypeString,
			Description: "Consumer tag",
			Default:     "",
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        consumerAutoAckKey,
			Type:        stream.PropertyTypeBool,
			Description: "Consumer Auto ACK",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        consumerNoLocalKey,
			Type:        stream.PropertyTypeBool,
			Description: "Consumer no local",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        consumerNoWaitKey,
			Type:        stream.PropertyTypeBool,
			Description: "Consumer no wait",
			Default:     false,
			Optional:    true,
		},
		//Exchange configurations
		stream.ConfigurationProperty{
			Name:        exchangeNameKey,
			Type:        stream.PropertyTypeString,
			Description: "Exchange name. If exchange name is empty all other exchange flags are ignored",
			Default:     "",
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        exchangeTypeKey,
			Type:        stream.PropertyTypeString,
			Description: "Exchange type - direct|fanout|topic|x-custom",
			Default:     "direct",
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        exchangeDurableKey,
			Type:        stream.PropertyTypeBool,
			Description: "Exchange durable",
			Default:     true,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        exchangeDeleteWhenCompleteKey,
			Type:        stream.PropertyTypeBool,
			Description: "Exchange delete when complete",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        exchangeInternalKey,
			Type:        stream.PropertyTypeBool,
			Description: "Exchange internal ",
			Default:     false,
			Optional:    true,
		},
		stream.ConfigurationProperty{
			Name:        exchangeNowaitKey,
			Type:        stream.PropertyTypeBool,
			Description: "Exchange no wait",
			Default:     false,
			Optional:    true,
		},
	},
}
