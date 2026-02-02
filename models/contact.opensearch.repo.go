package models

import (
	"bytes"
	"encoding/json"
	"io"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/rs/zerolog/log"
)

type OpenSearchContactStruct struct {
	OpenSearchClient *opensearch.Client
}

func OpenSearchContactRepository(client *opensearch.Client) OpenSearchContactSvcRepo {
	return &OpenSearchContactStruct{
		OpenSearchClient: client,
	}
}

type OpenSearchContactSvcRepo interface {
	ListByQueryMap(query map[string]any) ([]*OpenSearchContactSearchHit, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(contacts []*OpenSearchContact) (int64, error)
}

func (t *OpenSearchContactStruct) ListByQueryMap(query map[string]any) ([]*OpenSearchContactSearchHit, error) {
	if t.OpenSearchClient == nil {
		return nil, constants.OpenSearchNotConnectedError
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.OpenSearchClient.Search(
		t.OpenSearchClient.Search.WithIndex(constants.ContactIndex),
		t.OpenSearchClient.Search.WithBody(queryReader),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, constants.OpenSearchError(response.StatusCode, string(bodyBytes))
	}

	var searchResponse OpenSearchContactSearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return searchResponse.Hits.Hits, nil
}

func (t *OpenSearchContactStruct) CountByQueryMap(query map[string]any) (int64, error) {
	if t.OpenSearchClient == nil {
		return 0, constants.OpenSearchNotConnectedError
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.OpenSearchClient.Count(
		t.OpenSearchClient.Count.WithIndex(constants.ContactIndex),
		t.OpenSearchClient.Count.WithBody(queryReader),
	)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, constants.OpenSearchError(response.StatusCode, string(bodyBytes))
	}

	var countResponse utilities.OpenSearchCount
	if err := json.NewDecoder(response.Body).Decode(&countResponse); err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func (t *OpenSearchContactStruct) BulkUpsert(contacts []*OpenSearchContact) (int64, error) {
	if t.OpenSearchClient == nil {
		return 0, constants.OpenSearchNotConnectedError
	}
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

	response, err := t.OpenSearchClient.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, constants.OpenSearchBulkError(response.StatusCode, string(bodyBytes))
	}

	return int64(len(contacts)), nil
}
