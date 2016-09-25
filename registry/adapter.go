package registry

import (
	"fmt"

	"github.com/ottogiron/ferraritrunk/config"
	"github.com/ottogiron/ferrariworker/adapter"
)

type configurationRegistry struct {
	factory             adapter.Factory
	configurationSchema *config.ConfigurationSchema
}

var factoryRegistry = map[string]*configurationRegistry{}

//RegisterAdapterFactory registers a new factory for creating adapters
func RegisterAdapterFactory(factory adapter.Factory, adapterConfigurationSchema *config.ConfigurationSchema) error {
	if factoryRegistry[adapterConfigurationSchema.Name] != nil {
		return fmt.Errorf("The factory already exists %s", adapterConfigurationSchema.Name)
	}
	factoryRegistry[adapterConfigurationSchema.Name] = &configurationRegistry{factory, adapterConfigurationSchema}
	return nil
}

//AdapterFactory returns a factory for a registered adapter
func AdapterFactory(factoryName string) (adapter.Factory, error) {
	if factoryRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return factoryRegistry[factoryName].factory, nil
}

//AdapterSchema returns the configuration schema for the adapter factory
func AdapterSchema(factoryName string) (*config.ConfigurationSchema, error) {
	if factoryRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return factoryRegistry[factoryName].configurationSchema, nil
}

//AdapterSchemas returns the schemas for all the available adapters
func AdapterSchemas() []*config.ConfigurationSchema {
	configurationSchemas := []*config.ConfigurationSchema{}
	for _, registry := range factoryRegistry {
		configurationSchemas = append(configurationSchemas, registry.configurationSchema)
	}
	return configurationSchemas
}
