package elastic

import (
	"testing"

	"github.com/ottogiron/ferrariworker/backend"
)

func TestBackend(t *testing.T) {
	b := New()
	backend.TestBackend(t, b)
}
