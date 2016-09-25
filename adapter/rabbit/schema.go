package rabbit

import "github.com/ottogiron/ferraritrunk/config"

var schema = &config.ConfigurationSchema{
	Name:             "rabbit",
	ShortDescription: "processors jobs from rabbitmq",
	LongDescription:  "processors jobs from rabbibmq",
	Properties: []config.ConfigurationProperty{
		config.ConfigurationProperty{
			Name:        uriKey,
			Type:        config.PropertyTypeString,
			Description: "Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/",
			Default:     "amqp://guest:guest@localhost:5672/",
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        bindingKey,
			Type:        config.PropertyTypeString,
			Description: "Binding Key",
			Optional:    false,
		},
		config.ConfigurationProperty{
			Name:        bindingWaitKey,
			Type:        config.PropertyTypeString,
			Description: "Binding wait",
			Default:     false,
			Optional:    false,
		},
		config.ConfigurationProperty{
			Name:        queueNameKey,
			Type:        config.PropertyTypeString,
			Description: "Rabbit queue name",
			Optional:    false,
		},
		config.ConfigurationProperty{
			Name:        queueDurableKey,
			Type:        config.PropertyTypeBool,
			Description: "Queue durable",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        queueDeleteWhenUsedKey,
			Type:        config.PropertyTypeBool,
			Description: "Queue delete queue when used",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        queueExclusiveKey,
			Type:        config.PropertyTypeBool,
			Description: "Queue exclusive",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        queueNoWaitKey,
			Type:        config.PropertyTypeBool,
			Description: "Queue no wait",
			Default:     false,
			Optional:    true,
		},
		//Consumer configurations
		config.ConfigurationProperty{
			Name:        consumerTagKey,
			Type:        config.PropertyTypeString,
			Description: "Consumer tag",
			Default:     "",
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        consumerAutoAckKey,
			Type:        config.PropertyTypeBool,
			Description: "Consumer Auto ACK",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        consumerNoLocalKey,
			Type:        config.PropertyTypeBool,
			Description: "Consumer no local",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        consumerNoWaitKey,
			Type:        config.PropertyTypeBool,
			Description: "Consumer no wait",
			Default:     false,
			Optional:    true,
		},
		//Exchange configurations
		config.ConfigurationProperty{
			Name:        exchangeNameKey,
			Type:        config.PropertyTypeString,
			Description: "Exchange name. If exchange name is empty all other exchange flags are ignored",
			Default:     "",
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        exchangeTypeKey,
			Type:        config.PropertyTypeString,
			Description: "Exchange type - direct|fanout|topic|x-custom",
			Default:     "direct",
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        exchangeDurableKey,
			Type:        config.PropertyTypeBool,
			Description: "Exchange durable",
			Default:     true,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        exchangeDeleteWhenCompleteKey,
			Type:        config.PropertyTypeBool,
			Description: "Exchange delete when complete",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        exchangeInternalKey,
			Type:        config.PropertyTypeBool,
			Description: "Exchange internal ",
			Default:     false,
			Optional:    true,
		},
		config.ConfigurationProperty{
			Name:        exchangeNowaitKey,
			Type:        config.PropertyTypeBool,
			Description: "Exchange no wait",
			Default:     false,
			Optional:    true,
		},
	},
}
