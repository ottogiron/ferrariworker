package elastic

import (
	"os"
	"testing"

	"github.com/ottogiron/ferrariworker/backend"
	"github.com/ottogiron/ferrariworker/config"
)

func setUp(tb testing.TB) {

}

func tearDown(tb testing.TB) {

}

func TestBackend(t *testing.T) {
	if os.Getenv("TEST_ELASTIC_BACKEND") != "true" {
		t.Skip("Skiping elastic test")
	}
	config := config.NewAdapterConfig()
	config.Set(setSniffKey, false)
	config.Set(urlsKey, "http://127.0.0.1:9200")
	config.Set(indexKey, "workers_test")
	b, err := factory(config)

	if err != nil {
		t.Fatalf("elastic.NewClient() => err:%s while creating elastic client", err)
	}
	backend.TestBackend(t, b)
}
