package rabbit

import "github.com/ferrariframework/ferrariworker/processor"

var schema = &processor.AdapterConfigurationSchema{
	Name:             "rabbit",
	ShortDescription: "processors jobs from rabbitmq",
	LongDescription:  "processors jobs from rabbibmq",
	Properties: []processor.AdapterConfigurationProperty{
		processor.AdapterConfigurationProperty{
			Name:        uriKey,
			Type:        processor.PropertyTypeString,
			Description: "Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/",
			Default:     "amqp://guest:guest@localhost:5672/",
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        routingKey,
			Type:        processor.PropertyTypeString,
			Description: "Routing Key",
			Optional:    false,
		},
		processor.AdapterConfigurationProperty{
			Name:        bindingWaitKey,
			Type:        processor.PropertyTypeString,
			Description: "Binding wait",
			Default:     false,
			Optional:    false,
		},
		processor.AdapterConfigurationProperty{
			Name:        queueNameKey,
			Type:        processor.PropertyTypeString,
			Description: "Rabbit queue name",
			Optional:    false,
		},
		processor.AdapterConfigurationProperty{
			Name:        queueDurableKey,
			Type:        processor.PropertyTypeBool,
			Description: "Queue durable",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        queueDeleteWhenUsedKey,
			Type:        processor.PropertyTypeBool,
			Description: "Queue delete queue when used",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        queueExclusiveKey,
			Type:        processor.PropertyTypeBool,
			Description: "Queue exclusive",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        queueNoWaitKey,
			Type:        processor.PropertyTypeBool,
			Description: "Queue no wait",
			Default:     false,
			Optional:    true,
		},
		//Consumer configurations
		processor.AdapterConfigurationProperty{
			Name:        consumerTagKey,
			Type:        processor.PropertyTypeString,
			Description: "Consumer tag",
			Default:     "simple-consumer",
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        consumerAutoAckKey,
			Type:        processor.PropertyTypeBool,
			Description: "Consumer Auto ACK",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        consumerNoLocalKey,
			Type:        processor.PropertyTypeBool,
			Description: "Consumer no local",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        consumerNoWaitKey,
			Type:        processor.PropertyTypeBool,
			Description: "Consumer no wait",
			Default:     false,
			Optional:    true,
		},
		//Exchange configurations
		processor.AdapterConfigurationProperty{
			Name:        exchangeNameKey,
			Type:        processor.PropertyTypeString,
			Description: "Exchange name. If exchange name is empty all other exchange flags are ignored",
			Default:     "",
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        exchangeTypeKey,
			Type:        processor.PropertyTypeString,
			Description: "Exchange type - direct|fanout|topic|x-custom",
			Default:     "direct",
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        exchangeDurableKey,
			Type:        processor.PropertyTypeBool,
			Description: "Exchange durable",
			Default:     true,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        exchangeDeleteWhenCompleteKey,
			Type:        processor.PropertyTypeBool,
			Description: "Exchange delete when complete",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        exchangeInternalKey,
			Type:        processor.PropertyTypeBool,
			Description: "Exchange internal ",
			Default:     false,
			Optional:    true,
		},
		processor.AdapterConfigurationProperty{
			Name:        exchangeNowaitKey,
			Type:        processor.PropertyTypeBool,
			Description: "Exchange no wait",
			Default:     false,
			Optional:    true,
		},
		//Message configuration
		processor.AdapterConfigurationProperty{
			Name:        onSuccessKey,
			Type:        processor.PropertyTypeString,
			Description: `Action to execute on a message when the job succeded. Applied only when consumer_auto_ack=false. Example. --on_sucess="ack:false". Possible values |ack:<bool>|reject:<bool>|neck:multiple:<bool>,requeue:<bool>|`,
			Default:     "ack:false",
		},
		processor.AdapterConfigurationProperty{
			Name:        onFailedKey,
			Type:        processor.PropertyTypeString,
			Description: `Action to execute on a message when the job failed. Applied only when consumer_auto_ack=false. Example. --on_failed="ack:true". Possible values |ack:<bool>|reject:<bool>|neck:multiple:<bool>,requeue:<bool>|`,
			Default:     "ack:true",
		},
	},
}
