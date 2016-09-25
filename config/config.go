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
type Config interface {
	Set(string, interface{})
	GetString(key string) string
	GetInt(key string) int
	GetBoolean(key string) bool
}

//AdapterConfig a map of Configurations for a processor;
type config struct {
	adapterConfig map[string]interface{}
}

//NewAdapterConfig returns a new instance of the Configuration
func NewConfig() Config {
	return &config{
		adapterConfig: map[string]interface{}{},
	}
}

//Set sets a Configuration value
func (c *config) Set(key string, value interface{}) {
	c.adapterConfig[key] = value
}

//GetString returns a Configuration value as string
func (c *config) GetString(key string) string {
	return cast.ToString(c.adapterConfig[key])
}

//GetInt return a Configuration value as int
func (c *config) GetInt(key string) int {
	return cast.ToInt(c.adapterConfig[key])
}

//GetBoolean return a Configuration value as boolean
func (c *config) GetBoolean(key string) bool {
	return cast.ToBool(c.adapterConfig[key])
}

//ConfigurationSchema definition of Configurations for an adapter
type ConfigurationSchema struct {
	Name             string
	ShortDescription string
	LongDescription  string
	Properties       []ConfigurationProperty
}

//ConfigurationProperty returns an adapter configuration property
func (c *ConfigurationSchema) ConfigurationProperty(key string) *ConfigurationProperty {
	for _, property := range c.Properties {
		if property.Name == key {
			return &property
		}
	}
	return nil
}

//PropertyType for a Configuration property
type PropertyType int

//ConfigurationProperty definition of an specific Configuration property
type ConfigurationProperty struct {
	Name        string
	Type        PropertyType
	Default     interface{}
	Description string
	Optional    bool
}
