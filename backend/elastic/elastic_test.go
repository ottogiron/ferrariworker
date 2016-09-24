package elastic

import (
	"os"
	"testing"

	"gopkg.in/olivere/elastic.v2"

	"github.com/ottogiron/ferrariworker/backend"
)

func TestBackend(t *testing.T) {
	if os.Getenv("TEST_ELASTIC_BACKEND") != "true" {
		t.Skip("Skiping elastic test")
	}

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
	)

	if err != nil {
		t.Fatalf("elastic.NewClient() => err:%s while creating elastic client", err)
	}
	b := New(client)
	backend.TestBackend(t, b)
}
