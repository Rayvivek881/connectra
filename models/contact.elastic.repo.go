package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/elastic/go-elasticsearch/v8"
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
	ListByQueryMap(query map[string]any) ([]*ElasticContact, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(contacts []*ElasticContact) (int64, error)
	IndexContact(contact *ElasticContact) error
	DeleteContact(uuid string) error
}

func (t *ElasticContactStruct) ListByQueryMap(query map[string]any) ([]*ElasticContact, error) {
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
		return nil, fmt.Errorf("elasticsearch error: status %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	var searchResponse ElasticContactSearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	result := make([]*ElasticContact, 0, len(searchResponse.Hits.Hits))
	for _, hit := range searchResponse.Hits.Hits {
		result = append(result, &hit.Source)
	}
	return result, nil
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
		return 0, fmt.Errorf("elasticsearch error: status %d, body: %s", response.StatusCode, string(bodyBytes))
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
				"_id":    contact.Id,
			},
		}
		if utilities.AddToBuffer(&buf, meta) != nil || utilities.AddToBuffer(&buf, contact) != nil {
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
		return 0, fmt.Errorf("elasticsearch bulk error: status %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	return int64(len(contacts)), nil
}

// IndexContact indexes a single contact in Elasticsearch
func (t *ElasticContactStruct) IndexContact(contact *ElasticContact) error {
	contactJson, err := json.Marshal(contact)
	if err != nil {
		return fmt.Errorf("failed to marshal contact: %w", err)
	}

	response, err := t.ElasticClient.Index(
		constants.ContactIndex,
		bytes.NewReader(contactJson),
		t.ElasticClient.Index.WithDocumentID(contact.Id),
		t.ElasticClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to index contact: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("elasticsearch index error: status %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteContact removes a contact from Elasticsearch
func (t *ElasticContactStruct) DeleteContact(uuid string) error {
	response, err := t.ElasticClient.Delete(
		constants.ContactIndex,
		uuid,
		t.ElasticClient.Delete.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() && response.StatusCode != 404 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("elasticsearch delete error: status %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}
