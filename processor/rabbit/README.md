# Rabbit Processor Adapter
Provides jobs  from RabbitMQ

## Usage

```
 jobs from rabbibmq

Usage:
  ferrariworker process rabbit [flags]

Flags:
      --binding_wait string             Binding wait (default "false")
      --consumer_auto_ack               Consumer Auto ACK
      --consumer_no_local               Consumer no local
      --consumer_no_wait                Consumer no wait
      --consumer_tag string             Consumer tag (default "simple-consumer")
      --exchange string                 Exchange name. If exchange name is empty all other exchange flags are ignored
      --exchange_delete_when_complete   Exchange delete when complete
      --exchange_durable                Exchange durable (default true)
      --exchange_internal               Exchange internal
      --exchange_no_wait                Exchange no wait
      --exchange_type string            Exchange type - direct|fanout|topic|x-custom (default "direct")
      --on_failed string                Action to execute on a message when the job failed. Applied only when consumer_auto_ack=false. Example. --on_failed="ack:true". Possible values |ack:<bool>|reject:<bool>|neck:multiple:<bool>,requeue:<bool>| (default "ack:true")
      --on_sucess string                Action to execute on a message when the job succeded. Applied only when consumer_auto_ack=false. Example. --on_sucess="ack:false". Possible values |ack:<bool>|reject:<bool>|neck:multiple:<bool>,requeue:<bool>| (default "ack:false")
      --queue_delete_when_used          Queue delete queue when used
      --queue_durable                   Queue durable
      --queue_exclusive                 Queue exclusive
      --queue_name string               Rabbit queue name
      --queue_no_wait                   Queue no wait
      --routing_key string              Routing Key
      --uri string                      Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/ (default "amqp://guest:guest@localhost:5672/")
```