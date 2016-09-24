package elastic

import (
	"os"
	"testing"

	"github.com/ottogiron/ferrariworker/backend"
)

func TestBackend(t *testing.T) {
	if os.Getenv("ELASTIC_BACKEND_TEST") != "true" {
		t.Skip("Skiping elastic test")
	}
	b := New()
	backend.TestBackend(t, b)
}
