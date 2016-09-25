package elastic

import (
	"os"
	"testing"

	"gopkg.in/olivere/elastic.v2"

	"github.com/ottogiron/ferraritrunk/backend"
	"github.com/ottogiron/ferraritrunk/config"
)

var client *elastic.Client
var testIndex = "workers_test"

func init() {
	var err error
	client, err = elastic.NewClient(
		elastic.SetSniff(false),
	)

	if err != nil {
		panic("Failed to create a test elastic client")
	}
}

func setUp(tb testing.TB) {

}

func tearDown(tb testing.TB) {
	_, err := client.DeleteIndex(testIndex).Do()
	if err != nil {
		tb.Fatalf("Failed to remove index %s", err)
	}
}

func TestBackend(t *testing.T) {
	if os.Getenv("TEST_ELASTIC_BACKEND") != "true" {
		t.Skip("Skiping elastic test")
	}
	config := config.NewConfig()
	config.Set(setSniffKey, false)
	config.Set(urlsKey, "http://127.0.0.1:9200")
	config.Set(indexKey, "workers_test")
	config.Set(refreshIndexKey, true)
	b, err := factory(config)

	if err != nil {
		t.Fatalf("elastic.NewClient() => err:%s while creating elastic client", err)
	}
	backend.TestBackend(t, b, tearDown)
}
