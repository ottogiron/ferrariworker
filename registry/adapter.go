package registry

import (
	"fmt"

	"github.com/ottogiron/ferrariworker/adapter"
	"github.com/ottogiron/ferrariworker/config"
)

type configurationRegistry struct {
	factory             adapter.Factory
	configurationSchema *config.AdapterConfigurationSchema
}

var factoryRegistry = map[string]*configurationRegistry{}

//RegisterAdapterFactory registers a new factory for creating adapters
func RegisterAdapterFactory(factory adapter.Factory, adapterConfigurationSchema *config.AdapterConfigurationSchema) error {
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
func AdapterSchema(factoryName string) (*config.AdapterConfigurationSchema, error) {
	if factoryRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return factoryRegistry[factoryName].configurationSchema, nil
}

//AdapterSchemas returns the schemas for all the available adapters
func AdapterSchemas() []*config.AdapterConfigurationSchema {
	configurationSchemas := []*config.AdapterConfigurationSchema{}
	for _, registry := range factoryRegistry {
		configurationSchemas = append(configurationSchemas, registry.configurationSchema)
	}
	return configurationSchemas
}
