# Rabbit Processor Adapter
Provides jobs  from RabbitMQ

## Usage

```
 jobs from rabbibmq

Usage:
  ferrariworker process rabbit [flags]

Flags:
      --binding-wait string             Binding wait (default "false")
      --consumer-auto-ack               Consumer Auto ACK
      --consumer-no-local               Consumer no local
      --consumer-no-wait                Consumer no wait
      --consumer-tag string             Consumer tag (default "simple-consumer")
      --exchange string                 Exchange name. If exchange name is empty all other exchange flags are ignored
      --exchange-delete-when-complete   Exchange delete when complete
      --exchange-durable                Exchange durable (default true)
      --exchange-internal               Exchange internal
      --exchange-no-wait                Exchange no wait
      --exchange-type string            Exchange type - direct|fanout|topic|x-custom (default "direct")
      --on-failed string                Action to execute on a message when the job failed. Applied only when consumer-auto-ack=false. Example. --on-failed="ack:true". Possible values |ack:<bool>|reject:<bool>|neck:multiple:<bool>,requeue:<bool>| (default "ack:true")
      --on-sucess string                Action to execute on a message when the job succeded. Applied only when consumer-auto-ack=false. Example. --on-sucess="ack:false". Possible values |ack:<bool>|reject:<bool>|neck:multiple:<bool>,requeue:<bool>| (default "ack:false")
      --queue-delete-when-used          Queue delete queue when used
      --queue-durable                   Queue durable
      --queue-exclusive                 Queue exclusive
      --queue-name string               Rabbit queue name
      --queue-no-wait                   Queue no wait
      --routing-key string              Routing Key
      --uri string                      Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/ (default "amqp://guest:guest@localhost:5672/")
```