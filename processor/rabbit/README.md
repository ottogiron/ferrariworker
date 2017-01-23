# Rabbit Processor Adapter
Provides jobs  from RabbitMQ

## Usage

```
 jobs from rabbibmq

Usage:
  ferrariworker process rabbit [flags]

Flags:
      --consumer_auto_ack               Consumer Auto ACK
      --consumer_no_local               Consumer no local
      --consumer_no_wait                Consumer no wait
      --consumer_tag string             Consumer tag
      --exchange string                 Exchange name. If exchange name is empty all other exchange flags are ignored
      --exchange_delete_when_complete   Exchange delete when complete
      --exchange_durable                Exchange durable (default true)
      --exchange_internal               Exchange internal
      --exchange_no_wait                Exchange no wait
      --exchange_type string            Exchange type _ direct|fanout|topic|x_custom (default "direct")
      --queue_delete_when_used          Queue delete queue when used
      --queue_durable                   Queue durable
      --queue_exclusive                 Queue exclusive
      --queue_name string               Rabbit queue name
      --queue_no_wait                   Queue no wait
      --uri string                      Rabbit instance uri e.g. amqp://guest:guest@localhost:5672/ (default "amqp://guest:guest@localhost:5672/")

Global Flags:
      --command string            Command to be run when a new job arrives  (default "echo \"You should  set your own command :)\"")
      --command_run_path string   Running path for the command (default ".")
      --config string             config file (default is $HOME/.ferrariworker.yaml)
      --max_concurrency int       Max number of jobs proccessed concurrently. (default 1)
      --wait_timeout int          Time to wait in milliseconds for existing jobs to be processed.  (default 200)
```