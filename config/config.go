package config

import "github.com/spf13/cast"

const (
	//PropertyTypeString type for a string
	PropertyTypeString = iota
	//PropertyTypeInt type for an int
	PropertyTypeInt
	//PropertyTypeBool type for a boolean
	PropertyTypeBool
)

//AdapterConfig defines an adapter configuration
type AdapterConfig interface {
	Set(string, interface{})
	GetString(key string) string
	GetInt(key string) int
	GetBoolean(key string) bool
}

//AdapterConfig a map of AdapterConfigurations for a processor;
type adapterConfig struct {
	adapterConfig map[string]interface{}
}

//NewAdapterConfig returns a new instance of the AdapterConfiguration
func NewAdapterConfig() AdapterConfig {
	return &adapterConfig{
		adapterConfig: map[string]interface{}{},
	}
}

//Set sets a AdapterConfiguration value
func (c *adapterConfig) Set(key string, value interface{}) {
	c.adapterConfig[key] = value
}

//GetString returns a AdapterConfiguration value as string
func (c *adapterConfig) GetString(key string) string {
	return cast.ToString(c.adapterConfig[key])
}

//GetInt return a AdapterConfiguration value as int
func (c *adapterConfig) GetInt(key string) int {
	return cast.ToInt(c.adapterConfig[key])
}

//GetBoolean return a AdapterConfiguration value as boolean
func (c *adapterConfig) GetBoolean(key string) bool {
	return cast.ToBool(c.adapterConfig[key])
}

//AdapterConfigurationSchema definition of AdapterConfigurations for an adapter
type AdapterConfigurationSchema struct {
	Name             string
	ShortDescription string
	LongDescription  string
	Properties       []AdapterConfigurationProperty
}

//AdapterConfigurationProperty returns an adapter configuration property
func (c *AdapterConfigurationSchema) AdapterConfigurationProperty(key string) *AdapterConfigurationProperty {
	for _, property := range c.Properties {
		if property.Name == key {
			return &property
		}
	}
	return nil
}

//PropertyType for a AdapterConfiguration property
type PropertyType int

//AdapterConfigurationProperty definition of an specific AdapterConfiguration property
type AdapterConfigurationProperty struct {
	Name        string
	Type        PropertyType
	Default     interface{}
	Description string
	Optional    bool
}
