package stream

import "github.com/spf13/cast"

const (
	//PropertyTypeString type for a string
	PropertyTypeString = iota
	//PropertyTypeInt type for an int
	PropertyTypeInt
	//PropertyTypeBool type for a boolean
	PropertyTypeBool
)

type Config interface {
	Set(string, interface{})
	GetString(key string) string
	GetInt(key string) int
	GetBoolean(key string) bool
}

//Config a map of configurations for a stream;
type config struct {
	config map[string]interface{}
}

//NewConfig returns a new instance of the configuration
func NewConfig() Config {
	return &config{
		config: map[string]interface{}{},
	}
}

//Set sets a configuration value
func (c *config) Set(key string, value interface{}) {
	c.config[key] = value
}

//GetString returns a configuration value as string
func (c *config) GetString(key string) string {
	return cast.ToString(c.config[key])
}

//GetInt return a configuration value as int
func (c *config) GetInt(key string) int {
	return cast.ToInt(c.config[key])
}

//GetBoolean return a configuration value as boolean
func (c *config) GetBoolean(key string) bool {
	return cast.ToBool(c.config[key])
}

//ConfigurationSchema definition of configurations for an adapter
type ConfigurationSchema struct {
	Name             string
	ShortDescription string
	LongDescription  string
	Properties       []ConfigurationProperty
}

func (c *ConfigurationSchema) configurationProperty(key string) *ConfigurationProperty {
	for _, property := range c.Properties {
		if property.Name == key {
			return &property
		}
	}
	return nil
}

//PropertyType for a configuration property
type PropertyType int

//ConfigurationProperty definition of an specific configuration property
type ConfigurationProperty struct {
	Name        string
	Type        PropertyType
	Default     interface{}
	Description string
	Optional    bool
}
