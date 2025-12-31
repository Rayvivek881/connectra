package models

import (
	"bytes"
	"encoding/json"
	"io"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type ElasticContactStruct struct {
	ElasticClient *elasticsearch.Client
}

func ElasticContactRepository(client *elasticsearch.Client) ElasticContactSvcRepo {
	return &ElasticContactStruct{
		ElasticClient: client,
	}
}

type ElasticContactSvcRepo interface {
	ListByQueryMap(query map[string]any) ([]*ElasticContactSearchHit, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(contacts []*ElasticContact) (int64, error)
}

func (t *ElasticContactStruct) ListByQueryMap(query map[string]any) ([]*ElasticContactSearchHit, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.ElasticClient.Search(
		t.ElasticClient.Search.WithIndex(constants.ContactIndex),
		t.ElasticClient.Search.WithBody(queryReader),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, constants.ElasticsearchError(response.StatusCode, string(bodyBytes))
	}

	var searchResponse ElasticContactSearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return searchResponse.Hits.Hits, nil
}

func (t *ElasticContactStruct) CountByQueryMap(query map[string]any) (int64, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.ElasticClient.Count(
		t.ElasticClient.Count.WithIndex(constants.ContactIndex),
		t.ElasticClient.Count.WithBody(queryReader),
	)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, constants.ElasticsearchError(response.StatusCode, string(bodyBytes))
	}

	var countResponse utilities.ElasticCount
	if err := json.NewDecoder(response.Body).Decode(&countResponse); err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func (t *ElasticContactStruct) BulkUpsert(contacts []*ElasticContact) (int64, error) {
	var buf bytes.Buffer
	for _, contact := range contacts {
		meta := map[string]any{
			"index": map[string]any{
				"_index": constants.ContactIndex,
				"_id":    contact.UUID,
			},
		}
		if utilities.AddToBuffer(&buf, meta) != nil || utilities.AddToBuffer(&buf, contact) != nil {
			log.Error().Msgf("Failed to add contact to buffer: %v", contact.UUID)
			continue
		}
	}

	response, err := t.ElasticClient.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, constants.ElasticsearchBulkError(response.StatusCode, string(bodyBytes))
	}

	return int64(len(contacts)), nil
}
