package elastic

import (
	"reflect"
	"strings"

	"github.com/ottogiron/ferrariworker/backend"
	"github.com/ottogiron/ferrariworker/config"
	"github.com/ottogiron/ferrariworker/registry"
	"github.com/ottogiron/ferrariworker/worker"
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v2"
)

const (
	setSniffKey = "set-sniff"
	indexKey    = "index"
	urlsKey     = "urls"
)

const (
	docType        = "job_result"
	workerDocField = "worker_id"
)

func init() {
	registry.RegisterBackendFactory(factory, schema)
}

func factory(config config.AdapterConfig) (backend.Backend, error) {
	urls := strings.Split(config.GetString(urlsKey), ",")
	client, err := elastic.NewClient(
		elastic.SetSniff(config.GetBoolean(setSniffKey)),
		elastic.SetURL(urls...),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create elastic backend connection")
	}
	return new(client, config)
}

type elasticBackend struct {
	client *elastic.Client
	index  string
}

func new(client *elastic.Client, config config.AdapterConfig) (backend.Backend, error) {
	index := config.GetString(indexKey)
	mapping := `{
		"job_result":{
			"properties":{
				"worker_id":{
					"type":"string",
					 "index" : "not_analyzed"
				},
				"job_id":{
					"type":"string",
					 "index" : "not_analyzed"
				}
			}
		}
	}`
	_, err := client.PutMapping().Type(docType).BodyString(mapping).Do()

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to define default mappings for worker_id and job_id %s", mapping)
	}
	return &elasticBackend{client, index}, nil
}

func (e *elasticBackend) Persist(jobResults []*worker.JobResult) error {

	var workerID = ""

	if len(jobResults) > 0 {
		workerID = jobResults[0].WorkerID
	}

	bulkReq := e.client.
		Bulk().
		Index(e.index).
		Type(docType)

	for _, jr := range jobResults {
		doc := elastic.NewBulkIndexRequest().Index(e.index).Type(docType).Doc(jr)
		bulkReq.Add(doc)
	}
	_, err := bulkReq.Do()

	if err != nil {
		return errors.Wrapf(err, "Failed to persist list of processed jobs workerID", workerID)
	}
	return nil
}

func (e *elasticBackend) JobResults(workerID string) ([]*worker.JobResult, error) {
	query := elastic.NewTermQuery(workerDocField, workerID)
	results, err := e.client.
		Search(e.index).
		Type(docType).
		Query(query).
		Do()

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retreive jobs for workerID %s", workerID)
	}

	jobResults := []*worker.JobResult{}

	var ttyp *worker.JobResult

	for _, item := range results.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(*worker.JobResult); ok {
			jobResults = append(jobResults, t)
		}
	}

	return jobResults, nil
}

func (e *elasticBackend) Job(jobID string) (*worker.JobResult, error) {
	return nil, nil
}
