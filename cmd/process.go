package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/ottogiron/ferrariworker/config"
	"github.com/ottogiron/ferrariworker/processor"
	"github.com/ottogiron/ferrariworker/registry"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (

	//CommandKey key for the command to run
	CommandKey = "command"
	//CommandPathKey key for the command context path
	CommandPathKey = "command-run-path"
	//ConcurrencyKey key for the concurrency of the command
	ConcurrencyKey = "max-concurrency"
	//WaitTimeoutKey the time the worker will wait for new jobs in ms
	WaitTimeoutKey = "wait-timeout"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process of jobs",
	Long: `Process a  jobs based on custom configuration:

e.g.
	ferrariworker process \
	--command="node hello.js" \
	--command-run-path="/Users/ogiron" \
	--max-concurrency=8 
	rabbit --host=amqp://guest:guest@localhost:5672/ 
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(processCmd)
	processCmd.PersistentFlags().String(CommandKey, `echo "You should  set your own command :)"`, "Command to be run when a new job arrives ")
	viper.BindPFlag(CommandKey, processCmd.PersistentFlags().Lookup(CommandKey))
	processCmd.PersistentFlags().Int(ConcurrencyKey, 1, "Max number of jobs proccessed concurrently.")
	viper.BindPFlag(ConcurrencyKey, processCmd.PersistentFlags().Lookup(ConcurrencyKey))
	processCmd.PersistentFlags().String(CommandPathKey, ".", "Running path for the command")
	viper.BindPFlag(CommandPathKey, processCmd.PersistentFlags().Lookup(CommandPathKey))
	processCmd.PersistentFlags().Int(WaitTimeoutKey, 200, "Time to wait in milliseconds for existing jobs to be processed. ")
	viper.BindPFlag(WaitTimeoutKey, processCmd.PersistentFlags().Lookup(WaitTimeoutKey))
	initAdaptersSubCommands(processCmd)
}

func initAdaptersSubCommands(command *cobra.Command) {
	schemas := registry.AdapterSchemas()
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
			case config.PropertyTypeString:
				value := cast.ToString(schemaProperty.Default)
				subCmd.Flags().String(schemaProperty.Name, value, schemaProperty.Description)
				viper.BindPFlag(viperPropertyName, subCmd.Flags().Lookup(schemaProperty.Name))
			case config.PropertyTypeInt:
				value := cast.ToInt(schemaProperty.Default)
				subCmd.Flags().Int(schemaProperty.Name, value, schemaProperty.Description)
				viper.BindPFlag(viperPropertyName, subCmd.Flags().Lookup(schemaProperty.Name))
			case config.PropertyTypeBool:
				value := cast.ToBool(schemaProperty.Default)
				subCmd.Flags().Bool(schemaProperty.Name, value, schemaProperty.Description)
				viper.BindPFlag(viperPropertyName, subCmd.Flags().Lookup(schemaProperty.Name))
			}
		}
		command.AddCommand(subCmd)
	}
}

func adapterCommandAction(cmd *cobra.Command, args []string) {

	factory, err := registry.AdapterFactory(cmd.Name())

	if err != nil {
		log.Fatalf("There was an error jcreating starting the processing %s", err)
		return
	}

	command := viper.GetString(CommandKey)
	concurrency := viper.GetInt(ConcurrencyKey)
	commandPath := viper.GetString(CommandPathKey)
	waitTimeout := time.Duration(viper.GetInt(WaitTimeoutKey))
	config, err := parseAdapterConfiguration(cmd.Name())
	if err != nil {
		log.Fatalf("Couldn't parse configuration for %s %s", cmd.Name(), err)
	}
	adapter := factory(config)
	processorConfig := &processor.Config{
		Adapter:     adapter,
		Command:     command,
		CommandPath: commandPath,
		Concurrency: concurrency,
		WaitTimeout: waitTimeout,
	}
	//Default stdout and stder to os.Stdout and os.Stderr
	sp := processor.New(processorConfig, nil, nil)
	err = sp.Start()
	if err != nil {
		log.Fatalf("Could not start  processing %s", err)
	}
}

func parseAdapterConfiguration(name string) (config.AdapterConfig, error) {
	schema, err := registry.AdapterSchema(name)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load schema for adapter %s", name)
	}
	config := config.NewAdapterConfig()
	for _, propertyDefinition := range schema.Properties {
		viperPropertyName := viperConfigKey(name, propertyDefinition.Name)
		config.Set(propertyDefinition.Name, viper.Get(viperPropertyName))
	}
	return config, nil
}

func viperConfigKey(cmdName string, adapterPropertyName string) string {
	return cmdName + "-" + adapterPropertyName
}
