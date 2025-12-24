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
	ListByQueryMap(query map[string]any) ([]*ElasticCompany, error)
	CountByQueryMap(query map[string]any) (int64, error)
	BulkUpsert(companies []*ElasticCompany) (int64, error)
	IndexCompany(company *ElasticCompany) error
	DeleteCompany(uuid string) error
}

func (t *ElasticCompanyStruct) ListByQueryMap(query map[string]any) ([]*ElasticCompany, error) {
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

	result := make([]*ElasticCompany, 0, len(searchResponse.Hits.Hits))
	for _, hit := range searchResponse.Hits.Hits {
		result = append(result, &hit.Source)
	}

	return result, nil
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

// IndexCompany indexes a single company in Elasticsearch
func (t *ElasticCompanyStruct) IndexCompany(company *ElasticCompany) error {
	companyJson, err := json.Marshal(company)
	if err != nil {
		return fmt.Errorf("failed to marshal company: %w", err)
	}

	response, err := t.ElasticClient.Index(
		constants.CompanyIndex,
		bytes.NewReader(companyJson),
		t.ElasticClient.Index.WithDocumentID(company.Id),
		t.ElasticClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to index company: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("elasticsearch index error: status %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteCompany removes a company from Elasticsearch
func (t *ElasticCompanyStruct) DeleteCompany(uuid string) error {
	response, err := t.ElasticClient.Delete(
		constants.CompanyIndex,
		uuid,
		t.ElasticClient.Delete.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() && response.StatusCode != 404 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("elasticsearch delete error: status %d, body: %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}
