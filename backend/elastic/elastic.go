package elastic

import (
	"reflect"
	"strings"

	"github.com/ottogiron/ferraritrunk/backend"
	"github.com/ottogiron/ferraritrunk/config"

	"github.com/ottogiron/ferraritrunk/worker"
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v2"
)

const (
	setSniffKey     = "set-sniff"
	indexKey        = "index"
	urlsKey         = "urls"
	refreshIndexKey = "refresh-index"
)

const (
	docType       = "job_result"
	workerIDField = "worker_id"
	jobIDField    = "job_id"
)

func factory(config config.Config) (backend.Backend, error) {
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
	client  *elastic.Client
	index   string
	refresh bool
}

func new(client *elastic.Client, config config.Config) (backend.Backend, error) {

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
	index := config.GetString(indexKey)
	refresh := config.GetBoolean(refreshIndexKey)
	return &elasticBackend{client, index, refresh}, nil
}

func (e *elasticBackend) Persist(jobResults []*worker.JobResult) error {

	var workerID = ""

	if len(jobResults) > 0 {
		workerID = jobResults[0].WorkerID
	}

	bulkReq := e.client.
		Bulk().
		Index(e.index).
		Type(docType).
		Refresh(e.refresh)

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
	query := elastic.NewTermQuery(workerIDField, workerID)
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
		} else {
			return nil, errors.Errorf("Failed to deserialize job result %s", item)
		}
	}

	return jobResults, nil
}

func (e *elasticBackend) Job(jobID string) (*worker.JobResult, error) {
	query := elastic.NewTermQuery(jobIDField, jobID)
	results, err := e.client.
		Search(e.index).
		Type(docType).
		Query(query).
		Do()

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retreive job for jobID %s", jobID)
	}

	if results.TotalHits() != 1 {
		return nil, errors.Wrapf(err, "No hits for jobID %s", err)
	}
	var ttyp *worker.JobResult

	for _, item := range results.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(*worker.JobResult); ok {
			return t, nil
		} else {
			return nil, errors.Errorf("Failed to deserialize job result %s", item)
		}
	}
	return nil, nil
}
