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

type OpenSearchCompanyStruct struct {
	OpenSearchClient *opensearch.Client
}

func OpenSearchCompanyRepository(client *opensearch.Client) OpenSearchCompanySvcRepo {
	return &OpenSearchCompanyStruct{
		OpenSearchClient: client,
	}
}

type OpenSearchCompanySvcRepo interface {
	ListByQueryMap(query map[string]any) ([]*OpenSearchCompanySearchHit, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(companies []*OpenSearchCompany) (int64, error)
}

func (t *OpenSearchCompanyStruct) ListByQueryMap(query map[string]any) ([]*OpenSearchCompanySearchHit, error) {
	if t.OpenSearchClient == nil {
		return nil, constants.OpenSearchNotConnectedError
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.OpenSearchClient.Search(
		t.OpenSearchClient.Search.WithIndex(constants.CompanyIndex),
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

	var searchResponse OpenSearchCompanySearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	return searchResponse.Hits.Hits, nil
}

func (t *OpenSearchCompanyStruct) CountByQueryMap(query map[string]any) (int64, error) {
	if t.OpenSearchClient == nil {
		return 0, constants.OpenSearchNotConnectedError
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.OpenSearchClient.Count(
		t.OpenSearchClient.Count.WithIndex(constants.CompanyIndex),
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

func (t *OpenSearchCompanyStruct) BulkUpsert(companies []*OpenSearchCompany) (int64, error) {
	if t.OpenSearchClient == nil {
		return 0, constants.OpenSearchNotConnectedError
	}
	var buf bytes.Buffer
	for _, company := range companies {
		meta := map[string]any{
			"index": map[string]any{
				"_index": constants.CompanyIndex,
				"_id":    company.UUID,
			},
		}
		if utilities.AddToBuffer(&buf, meta) != nil || utilities.AddToBuffer(&buf, company) != nil {
			log.Error().Msgf("Failed to add company to buffer: %v", company.UUID)
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

	return int64(len(companies)), nil
}
