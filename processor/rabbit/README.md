# Rabbit Processor Adapter
Provides jobs  from RabbitMQ

## Usage

```
 jobs from rabbibmq

Usage:
  ferrariprocessor process rabbit [flags]

Flags:
      --consumer-auto-ack               Consumer Auto ACK
      --consumer-no-local               Consumer no local
      --consumer-no-wait                Consumer no wait
      --consumer-tag string             Consumer tag
      --exchange string                 Exchange name. If exchange name is empty all other exchange flags are ignored
      --exchange-delete-when-complete   Exchange delete when complete
      --exchange-durable                Exchange durable (default true)
      --exchange-internal               Exchange internal
      --exchange-no-wait                Exchange no wait
      --exchange-type string            Exchange type - direct|fanout|topic|x-custom (default "direct")
      --queue-delete-when-used          Queue delete queue when used
      --queue-durable                   Queue durable
      --queue-exclusive                 Queue exclusive
      --queue-name string               Rabbit queue name
      --queue-no-wait                   Queue no wait
      --uri string                      Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/ (default "amqp://guest:guest@localhost:5672/")

Global Flags:
      --command string            Command to be run when a new job arrives  (default "echo \"You should  set your own command :)\"")
      --command-run-path string   Running path for the command (default ".")
      --config string             config file (default is $HOME/.ferrariprocessor.yaml)
      --max-concurrency int       Max number of jobs proccessed concurrently. (default 1)
      --wait-timeout int          Time to wait in milliseconds for existing jobs to be processed.  (default 200)
```