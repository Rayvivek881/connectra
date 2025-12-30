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

type ElasticCompanyStruct struct {
	ElasticClient *elasticsearch.Client
}

func ElasticCompanyRepository(client *elasticsearch.Client) ElasticCompanySvcRepo {
	return &ElasticCompanyStruct{
		ElasticClient: client,
	}
}

type ElasticCompanySvcRepo interface {
	ListByQueryMap(query map[string]any) ([]*ElasticCompanySearchHit, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(companies []*ElasticCompany) (int64, error)
}

func (t *ElasticCompanyStruct) ListByQueryMap(query map[string]any) ([]*ElasticCompanySearchHit, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.ElasticClient.Search(
		t.ElasticClient.Search.WithIndex(constants.CompanyIndex),
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

	var searchResponse ElasticCompanySearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	return searchResponse.Hits.Hits, nil
}

func (t *ElasticCompanyStruct) CountByQueryMap(query map[string]any) (int64, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	queryReader := bytes.NewReader(queryJson)
	response, err := t.ElasticClient.Count(
		t.ElasticClient.Count.WithIndex(constants.CompanyIndex),
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

func (t *ElasticCompanyStruct) BulkUpsert(companies []*ElasticCompany) (int64, error) {
	var buf bytes.Buffer
	for _, company := range companies {
		meta := map[string]any{
			"index": map[string]any{
				"_index": constants.CompanyIndex,
				"_id":    company.Id,
			},
		}
		if utilities.AddToBuffer(&buf, meta) != nil || utilities.AddToBuffer(&buf, company) != nil {
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

	return int64(len(companies)), nil
}
