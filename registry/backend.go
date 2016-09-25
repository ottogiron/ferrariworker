package registry

import (
	"fmt"

	"github.com/ottogiron/ferraritrunk/backend"
	"github.com/ottogiron/ferraritrunk/config"
)

type backendConfigurationRegistry struct {
	factory             backend.Factory
	configurationSchema *config.ConfigurationSchema
}

var backendRegistry = map[string]*backendConfigurationRegistry{}

//RegisterBackendFactory registers a new factory for creating adapters
func RegisterBackendFactory(factory backend.Factory, adapterConfigurationSchema *config.ConfigurationSchema) error {
	if backendRegistry[adapterConfigurationSchema.Name] != nil {
		return fmt.Errorf("The factory already exists %s", adapterConfigurationSchema.Name)
	}
	backendRegistry[adapterConfigurationSchema.Name] = &backendConfigurationRegistry{factory, adapterConfigurationSchema}
	return nil
}

//BackendFactory returns a factory for a registered adapter
func BackendFactory(factoryName string) (backend.Factory, error) {
	if backendRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return backendRegistry[factoryName].factory, nil
}

//BackendSchema returns the configuration schema for the adapter factory
func BackendSchema(factoryName string) (*config.ConfigurationSchema, error) {
	if backendRegistry[factoryName] == nil {
		return nil, fmt.Errorf("The adapter %s is not registered Processor cannot be created", factoryName)
	}
	return backendRegistry[factoryName].configurationSchema, nil
}

//BackendSchemas returns the schemas for all the available adapters
func BackendSchemas() []*config.ConfigurationSchema {
	configurationSchemas := []*config.ConfigurationSchema{}
	for _, registry := range backendRegistry {
		configurationSchemas = append(configurationSchemas, registry.configurationSchema)
	}
	return configurationSchemas
}
