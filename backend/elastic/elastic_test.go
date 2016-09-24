package elastic

import (
	"os"
	"testing"

	"github.com/ottogiron/ferrariworker/backend"
)

func TestBackend(t *testing.T) {
	if os.Getenv("TEST_ELASTIC_BACKEND") != "true" {
		t.Skip("Skiping elastic test")
	}
	b := New()
	backend.TestBackend(t, b)
}
