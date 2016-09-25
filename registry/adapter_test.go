package registry

import (
	"testing"

	"github.com/ottogiron/ferraritrunk/config"
)

func TestRegisterAdapterFactory(t *testing.T) {

	cs := &config.ConfigurationSchema{
		Name: "test",
	}

	err := RegisterAdapterFactory(nil, cs)

	if err != nil {
		t.Errorf("The first factory should be registered correctly for %s", cs.Name)
	}

	err = RegisterAdapterFactory(nil, cs)

	if err == nil {
		t.Errorf("The registration should fail for %s", cs.Name)
	}
}

func TestGetConfigurationSchemas(t *testing.T) {
	RegisterAdapterFactory(nil, &config.ConfigurationSchema{
		Name: "test",
	})

	schemas := AdapterSchemas()

	slen := len(schemas)
	if slen != 1 {
		t.Errorf("Expected schemas size %d was", slen)
	}
}

func TestGetAdapterSchema(t *testing.T) {
	RegisterAdapterFactory(nil, &config.ConfigurationSchema{Name: "test"})

	schema, _ := AdapterSchema("test")

	if schema.Name != "test" {
		t.Errorf("expected schema name to be 'test' was %s", schema.Name)
	}
}
