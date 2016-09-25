package registry

import (
	"testing"

	"github.com/ottogiron/ferrariworker/config"
)

func TestRegisterBackendFactory(t *testing.T) {

	cs := &config.AdapterConfigurationSchema{
		Name: "test",
	}

	err := RegisterBackendFactory(nil, cs)

	if err != nil {
		t.Errorf("The first factory should be registered correctly for %s", cs.Name)
	}

	err = RegisterAdapterFactory(nil, cs)

	if err == nil {
		t.Errorf("The registration should fail for %s", cs.Name)
	}
}

func TestGetBackendSchemas(t *testing.T) {
	RegisterBackendFactory(nil, &config.AdapterConfigurationSchema{
		Name: "test",
	})

	schemas := AdapterSchemas()

	slen := len(schemas)
	if slen != 1 {
		t.Errorf("Expected schemas size %d was", slen)
	}
}

func TestGetBackendrSchema(t *testing.T) {
	RegisterBackendFactory(nil, &config.AdapterConfigurationSchema{Name: "test"})

	schema, _ := AdapterSchema("test")

	if schema.Name != "test" {
		t.Errorf("expected schema name to be 'test' was %s", schema.Name)
	}
}
