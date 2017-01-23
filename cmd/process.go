package cmd

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/ferrariframework/ferrariserver/grpc/gen"
	"github.com/ferrariframework/ferrariworker/processor"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (

	//CommandKey key for the command to run
	commandKey = "command"
	//commandPathKey key for the command context path
	commandPathKey = "command_run_path"
	//concurrencyKey key for the concurrency of the command
	concurrencyKey = "max_concurrency"
	//waitTimeoutKey the time the worker will wait for new jobs in ms
	waitTimeoutKey        = "wait_timeout"
	serverAddressKey      = "server_addr"
	tlsKey                = "tls"
	caFileKey             = "ca_file"
	serverHostOverrideKey = "server_host_override"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process of jobs",
	Long: `Process a  jobs based on custom configuration:

e.g.
	ferrariworker process \
	--command="node hello.js" \
	--command_run_path="/Users/ogiron" \
	--max_concurrency=8 
	rabbit --host=amqp://guest:guest@localhost:5672/ 
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(processCmd)
	processCmd.PersistentFlags().StringP(commandKey, "c", `echo "You should  set your own command :)"`, "Command to be run when a new job arrives ")
	processCmd.PersistentFlags().IntP(concurrencyKey, "m", 1, "Max number of jobs proccessed concurrently.")
	processCmd.PersistentFlags().StringP(commandPathKey, "p", ".", "Running path for the command")
	processCmd.PersistentFlags().IntP(waitTimeoutKey, "w", 200, "Time to wait in milliseconds for existing jobs to be processed. ")
	processCmd.PersistentFlags().BoolP(tlsKey, "t", false, "Connection uses TLS if true, else plain TCP")
	processCmd.PersistentFlags().StringP(serverAddressKey, "s", "127.0.0.1:4051", "The server address in the format of host:port")
	processCmd.PersistentFlags().String(caFileKey, "ca.pem", "he file containning the CA root cert file")
	processCmd.PersistentFlags().String(serverHostOverrideKey, "x.test.ferrariframework.com", "The server name use to verify the hostname returned by TLS handshake")

	viper.BindPFlag(commandPathKey, processCmd.PersistentFlags().Lookup(commandPathKey))
	viper.BindPFlag(concurrencyKey, processCmd.PersistentFlags().Lookup(concurrencyKey))
	viper.BindPFlag(commandKey, processCmd.PersistentFlags().Lookup(commandKey))
	viper.BindPFlag(waitTimeoutKey, processCmd.PersistentFlags().Lookup(waitTimeoutKey))
	viper.BindPFlag(tlsKey, processCmd.PersistentFlags().Lookup(tlsKey))
	viper.BindPFlag(serverAddressKey, processCmd.PersistentFlags().Lookup(serverAddressKey))
	viper.BindPFlag(caFileKey, processCmd.PersistentFlags().Lookup(caFileKey))

	initAdaptersSubCommands(processCmd)
}

func initAdaptersSubCommands(command *cobra.Command) {
	schemas := processor.AdapterSchemas()
	for _, schema := range schemas {
		subCmd := &cobra.Command{
			Use:   schema.Name,
			Short: schema.ShortDescription,
			Long:  schema.LongDescription,
			Run:   adapterCommandAction,
		}

		for _, schemaProperty := range schema.Properties {
			viperPropertyName := viperConfigKey(schema.Name, schemaProperty.Name)
			switch schemaProperty.Type {
			case processor.PropertyTypeString:
				value := cast.ToString(schemaProperty.Default)
				subCmd.Flags().String(schemaProperty.Name, value, schemaProperty.Description)
				viper.BindPFlag(viperPropertyName, subCmd.Flags().Lookup(schemaProperty.Name))
			case processor.PropertyTypeInt:
				value := cast.ToInt(schemaProperty.Default)
				subCmd.Flags().Int(schemaProperty.Name, value, schemaProperty.Description)
				viper.BindPFlag(viperPropertyName, subCmd.Flags().Lookup(schemaProperty.Name))
			case processor.PropertyTypeBool:
				value := cast.ToBool(schemaProperty.Default)
				subCmd.Flags().Bool(schemaProperty.Name, value, schemaProperty.Description)
				viper.BindPFlag(viperPropertyName, subCmd.Flags().Lookup(schemaProperty.Name))
			}
		}
		command.AddCommand(subCmd)
	}
}

func adapterCommandAction(cmd *cobra.Command, args []string) {

	factory, err := processor.AdapterFactory(cmd.Name())

	if err != nil {
		log.Fatalf("There was an error jcreating starting the processing %s", err)
		return
	}

	//Worker flags
	command := viper.GetString(commandKey)
	concurrency := viper.GetInt(concurrencyKey)
	commandPath := viper.GetString(commandPathKey)
	waitTimeout := time.Duration(viper.GetInt(waitTimeoutKey))

	//Server flags
	tls := viper.GetBool(tlsKey)
	serverHostOverride := viper.GetString(serverHostOverrideKey)
	caFile := viper.GetString(caFileKey)
	serverAddr := viper.GetString(serverAddressKey)

	jClient, close := jobServiceClient(serverAddr, caFile, serverHostOverride, tls)
	defer close()
	config, err := parseAdapterConfiguration(cmd.Name())

	if err != nil {
		log.Fatalf("Couldn't parse configuration for %s %s", cmd.Name(), err)
	}

	adapter := factory.New(config)
	processorConfig := &processor.Config{
		WorkerID:    "dummyWorkerID",
		Adapter:     adapter,
		Command:     command,
		CommandPath: commandPath,
		Concurrency: concurrency,
		WaitTimeout: waitTimeout,
	}
	//Default stdout and stder to os.Stdout and os.Stderr
	sp := processor.New(processorConfig, jClient, nil, nil)
	err = sp.Start()
	if err != nil {
		log.Fatalf("Could not start  processing %s", err)
	}
}

//JobServiceClient
func jobServiceClient(serverAddr, caFile, serverHostOverride string, tls bool) (gen.JobServiceClient, func()) {

	var opts []grpc.DialOption
	if tls {
		var sn string
		if serverHostOverride != "" {
			sn = serverHostOverride
		}
		var creds credentials.TransportCredentials
		if caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(caFile, sn)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	client := gen.NewJobServiceClient(conn)
	return client, func() {
		conn.Close()
	}
}

func parseAdapterConfiguration(name string) (processor.AdapterConfig, error) {
	schema, err := processor.AdapterSchema(name)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load schema for adapter %s", name)
	}
	config := processor.NewAdapterConfig()
	for _, propertyDefinition := range schema.Properties {
		viperPropertyName := viperConfigKey(name, propertyDefinition.Name)
		config.Set(propertyDefinition.Name, viper.Get(viperPropertyName))
	}
	return config, nil
}

func viperConfigKey(cmdName string, adapterPropertyName string) string {
	return cmdName + "_" + adapterPropertyName
}
